package zkteco

import (
	"fmt"
	"log"
	"regexp"
	"sort"
	"strconv"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/northbright/xls2csv-go/xls2csv"
)

// DB represents zkteco attendance database.
type DB struct {
	// RedisAddr is Redis address. e.g. ":6379".
	RedisAddr string
	// RedisPassword is Redis password.
	RedisPassword string
}

const (
	defSheetID = 2
)

var (
	debugMode = false
)

// SetDebugMode sets debug mode.
func SetDebugMode(f bool) {
	debugMode = f
}

// IsDebugMode returns if it's in debug mode or not.
func IsDebugMode() bool {
	return debugMode
}

// logFnResult outputs the result of the function.
//
// params:
//     funcName: function name.
//     err: result of function.
func logFnResult(funcName string, err error) {
	if !debugMode {
		return
	}

	if err != nil {
		log.Printf("%v() error: %v", funcName, err)
		return
	}

	log.Printf("%v() ok", funcName)
}

func dbgLog(fmt string, v ...interface{}) {
	if debugMode {
		log.Printf(fmt, v...)
	}
}

// Open opens a DB which contains zkteco attendance data.
func Open(redisAddr, redisPassword string) *DB {
	return &DB{redisAddr, redisPassword}
}

// getAttendance gets the attendance data by given employee name and date.
func getAttendance(c redis.Conn, name, date string) (string, string, error) {
	k := fmt.Sprintf("kaoqin:%v", name)

	field := fmt.Sprintf("%v:in", date)
	clockIn, err := redis.String(c.Do("HGET", k, field))
	if err != nil && err != redis.ErrNil {
		return "", "", err
	}

	field = fmt.Sprintf("%v:out", date)
	clockOut, err := redis.String(c.Do("HGET", k, field))
	if err != nil && err != redis.ErrNil {
		return "", "", err
	}

	return clockIn, clockOut, nil
}

// updateAttendance updates employee attendance data in Redis.
func updateAttendance(c redis.Conn, name, date, clockIn, clockOut string) error {
	var arr []string

	oldClockIn, oldClockOut, err := getAttendance(c, name, date)
	if err != nil {
		return err
	}

	// Use map to emulate set
	m := map[string]bool{}

	if oldClockIn != "" {
		m[oldClockIn] = true
	}

	if oldClockOut != "" {
		m[oldClockOut] = true
	}

	if clockIn != "" {
		m[clockIn] = true
	}

	if clockOut != "" {
		m[clockOut] = true
	}

	for k := range m {
		arr = append(arr, k)
	}

	// Sort
	sort.Strings(arr)
	dbgLog("sorted arr: %v", arr)

	l := len(arr)
	switch l {
	case 0:
		return nil
	case 1:
		clockIn = arr[0]
		k := fmt.Sprintf("kaoqin:%v", name)
		field := fmt.Sprintf("%v:in", date)
		value := clockIn
		if _, err = c.Do("HSET", k, field, value); err != nil {
			return err
		}
	case 2, 3, 4:
		clockIn = arr[0]
		clockOut = arr[l-1]
		k := fmt.Sprintf("kaoqin:%v", name)
		field := fmt.Sprintf("%v:in", date)
		value := clockIn
		if _, err = c.Do("HSET", k, field, value); err != nil {
			return err
		}

		field = fmt.Sprintf("%v:out", date)
		value = clockOut
		if _, err = c.Do("HSET", k, field, value); err != nil {
			return err
		}
	}
	return nil
}

// UpdateAttendance updates employees' attendance data by XLS file which is outputted by zkteco device.
func (db *DB) UpdateAttendance(xlsFile string) error {
	var err error
	var records [][]string

	defer func() {
		logFnResult("UpdateAttendance", err)
	}()

	records, err = xls2csv.XLS2CSV(xlsFile, defSheetID)
	if err != nil {
		return err
	}

	// Check row numbers.
	n := len(records)
	if n <= 4 {
		err = fmt.Errorf("rows of XLS file <= 4")
		return err
	}

	if n%2 != 0 {
		err = fmt.Errorf("rows of XLS file mod 2 != 0")
		return err
	}

	// Get begin date of attendace.
	p := `^(\d{4})-(\d{2})-(\d{2}) ~ \d{4}-\d{2}-\d{2}$`
	re := regexp.MustCompile(p)
	matched := re.FindStringSubmatch(records[2][2])
	if len(matched) != 4 {
		err = fmt.Errorf("can not find begin / end date of attendace")
		return err
	}
	year, err := strconv.Atoi(matched[1])
	if err != nil {
		return err
	}

	month, err := strconv.Atoi(matched[2])
	if err != nil {
		return err
	}

	day, err := strconv.Atoi(matched[3])
	if err != nil {
		return err
	}

	loc, err := time.LoadLocation("Local")
	if err != nil {
		return err
	}

	startTime := time.Date(int(year), time.Month(month), int(day), 0, 0, 0, 0, loc)
	dbgLog("startTime: %v", startTime)

	// Get day numbers.
	colNums := len(records[3])
	days := 0
	for i := 0; i <= colNums-1; i++ {
		if records[3][i] != "" {
			days++
		}
	}
	dbgLog("days: %v", days)

	// Clock in / off pattern
	// p1: have clock in and out.
	p1 := `^(\d{2}:\d{2})(\d{2}:\d{2})$`
	// p2: have clock in only.
	p2 := `^(\d{2}:\d{2})$`
	re1 := regexp.MustCompile(p1)
	re2 := regexp.MustCompile(p2)

	// Get Redis Connection
	c, err := GetRedisConn(db.RedisAddr, db.RedisPassword)
	if err != nil {
		return err
	}
	defer c.Close()

	// Get attendance data.
	for i := 4; i+1 <= n-1; i += 2 {
		// Get name.
		name := records[i][10]
		dbgLog("name: %v", name)
		for j := 0; j <= days-1; j++ {
			// Get date by adding days to start date.
			t := startTime.AddDate(0, 0, j)
			y, m, d := t.Date()
			date := fmt.Sprintf("%04d-%02d-%02d", y, int(m), d)
			dbgLog("date: %v", date)

			// Get clock in / off time for each work day.
			clockIn := ""
			clockOut := ""

			matched := re1.FindStringSubmatch(records[i+1][j])
			if len(matched) == 3 {
				clockIn = matched[1]
				clockOut = matched[2]
			} else {
				matched := re2.FindStringSubmatch(records[i+1][j])
				if len(matched) == 2 {
					clockIn = matched[1]
				}
			}
			dbgLog("clock in/out: %v, %v", clockIn, clockOut)

			if err = updateAttendance(c, name, date, clockIn, clockOut); err != nil {
				return err
			}
		}
	}

	return nil
}
