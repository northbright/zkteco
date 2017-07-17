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

	// Create an Util instance.
	util := zkteco.NewUtil(":6379", "")

	// Get absolute path of example attendance xls file.
	f, _ := pathhelper.GetCurrentExecDir()
	f = path.Join(f, "files/campus-a.xls")

	// Update attendance.
	err := util.UpdateAttendance(f)
	if err != nil {
		log.Printf("util.UpdateAttendance() failed: %v", err)
		return
	}

	log.Printf("util.UpdateAttendance ok")
	// Output:
}
