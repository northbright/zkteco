package zkteco

import (
	"log"

	"github.com/northbright/xls2csv-go/xls2csv"
)

var (
	debugMode = false
)

// SetDebugMode sets debug mode for package unifi.
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

func GetAttendances(xlsFile string) ([][]string, error) {
	var err error
	var records [][]string
	sheetId := 2

	defer func() {
		logFnResult("GetAttendances", err)
	}()

	records, err = xls2csv.XLS2CSV(xlsFile, sheetId)
	if err != nil {
		return records, err
	}

	return records, nil
}
