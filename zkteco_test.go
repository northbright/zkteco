package zkteco_test

import (
	"log"
	"path"

	"github.com/northbright/pathhelper"
	"github.com/northbright/zkteco"
)

func ExampleGetAttendances() {

	zkteco.SetDebugMode(true)

	// Get absolute path of example xls file.
	f, _ := pathhelper.GetCurrentExecDir()
	f = path.Join(f, "files/campus-b.xls")

	records, err := zkteco.GetAttendances(f)
	if err != nil {
		log.Printf("zkteco.GetAttendances() failed: %v", err)
		return
	}

	log.Printf("zkteco.GetAttendances() ok")
	for i, v := range records {
		log.Printf("len: %v", len(v))
		log.Printf("%v: %v", i, v)

	}

	// Output:
}
