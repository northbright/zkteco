package zkteco_test

import (
	"log"
	"path"

	"github.com/northbright/pathhelper"
	"github.com/northbright/zkteco"
)

func Example() {
	// Set debug mode to true.
	zkteco.SetDebugMode(true)

	// Open a DB.
	db := zkteco.Open(":6379", "")

	// Get absolute path of example attendance xls file.
	f, _ := pathhelper.GetCurrentExecDir()
	f = path.Join(f, "files/campus-a.xls")

	// Update attendance.
	err := db.UpdateAttendance(f)
	if err != nil {
		log.Printf("db.UpdateAttendance() failed: %v", err)
		return
	}

	log.Printf("db.UpdateAttendance ok")
	// Output:
}
