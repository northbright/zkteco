package zkteco

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/northbright/xls2csv-go/xls2csv"
)

const (
	defSheetID = 2
)

type Attendance struct {
	ID       string
	Name     string
	Date     string
	ClockIn  string
	ClockOut string
}

// Parse parses the XLS file output from ZKTeco devices and return the attendance data.
func Parse(xlsFile string) ([]Attendance, error) {
	var err error
	var records [][]string

	records, err = xls2csv.XLS2CSV(xlsFile, defSheetID)
	if err != nil {
		return nil, err
	}

	// Check row numbers.
	n := len(records)
	if n <= 3 {
		return nil, fmt.Errorf("rows of XLS file <= 3")
	}

	// Get begin date of attendace.
	p := `^(\d{4})-(\d{2})-(\d{2}) ~ \d{4}-\d{2}-\d{2}$`
	re := regexp.MustCompile(p)
	matched := re.FindStringSubmatch(records[1][2])
	if len(matched) != 4 {
		return nil, fmt.Errorf("can not find begin / end date of attendace")
	}
	year, err := strconv.Atoi(matched[1])
	if err != nil {
		return nil, err
	}

	month, err := strconv.Atoi(matched[2])
	if err != nil {
		return nil, err
	}

	day, err := strconv.Atoi(matched[3])
	if err != nil {
		return nil, err
	}

	loc, err := time.LoadLocation("Local")
	if err != nil {
		return nil, err
	}

	startTime := time.Date(int(year), time.Month(month), int(day), 0, 0, 0, 0, loc)

	// Get day numbers.
	colNums := len(records[2])
	days := 0
	for i := 0; i <= colNums-1; i++ {
		if records[2][i] != "" {
			days++
		}
	}

	// Clock in / off pattern
	// p1: have clock in and out.
	p1 := `^(\d{2}:\d{2})(\d{2}:\d{2})$`
	// p2: have clock in only.
	p2 := `^(\d{2}:\d{2})$`
	re1 := regexp.MustCompile(p1)
	re2 := regexp.MustCompile(p2)

	// Attendance slice to return.
	var s []Attendance

	// Get attendance data.
	for i := 3; i+1 <= n-1; i += 2 {
		// Get ID.
		id := records[i][2]

		// Get name.
		name := records[i][10]

		for j := 0; j <= days-1; j++ {
			// Get date by adding days to start date.
			t := startTime.AddDate(0, 0, j)
			year, month, day := t.Date()
			date := fmt.Sprintf("%04d-%02d-%02d", year, month, day)

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

			if clockIn != "" || clockOut != "" {
				// Make an attendance data.
				a := Attendance{ID: id, Name: name, Date: date, ClockIn: clockIn, ClockOut: clockOut}

				// Push the attendance data to the slice.
				s = append(s, a)
			}
		}
	}

	return s, nil
}
