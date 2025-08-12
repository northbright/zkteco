package zkteco_test

import (
	"fmt"
	"log"
	"sort"

	"github.com/northbright/zkteco"
)

// Record presents the clock in / out time.
type Record struct {
	ClockIn  string
	ClockOut string
}

var (
	// Use a map to store attendance records.
	// Key is "date:name".
	records = make(map[string]Record)
)

// Update the clock in / out time of one employee's working day.
// If the key exists, it gets the old clock in / out time.
// It puts clock in / out, old clock in / out time in a string slice and sorts the slice.
// The first element is the new clock in time and last element is the new clock out time.
func UpdateAttendance(name, date, clockIn, clockOut string) {
	s := []string{}

	if clockIn != "" {
		s = append(s, clockIn)
	}

	if clockOut != "" {
		s = append(s, clockOut)
	}

	k := fmt.Sprintf("%s:%s", name, date)

	v, ok := records[k]
	if ok {
		if v.ClockIn != "" {
			s = append(s, v.ClockIn)
		}
		if v.ClockOut != "" {
			s = append(s, v.ClockOut)
		}
	}

	if len(s) == 0 {
		return
	}

	// Sort times(strings).
	sort.Strings(s)
	r := Record{}

	r.ClockIn = s[0]
	if len(s) > 1 {
		// Use last time as clock out time.
		r.ClockOut = s[len(s)-1]
	}

	// Update attendance
	records[k] = r
}

func Example() {

	xlsFiles := []string{
		"files/campus-a.xls",
		"files/campus-b.xls",
	}

	for _, f := range xlsFiles {
		// Parse the XLS and get attendance data.
		s, err := zkteco.Parse(f)
		if err != nil {
			log.Printf("zkteco.Parse() error: %v", err)
			return
		}
		log.Printf("zkteco.Parse(): %s OK", f)

		fmt.Printf("\nrecords from %s\n", f)
		for _, a := range s {
			fmt.Printf("id: %s, name: %s, date: %s, clock in/out: %s/%s\n", a.ID, a.Name, a.Date, a.ClockIn, a.ClockOut)
			UpdateAttendance(a.Name, a.Date, a.ClockIn, a.ClockOut)
		}
	}

	// Show final attendances.
	fmt.Printf("\nfinal records(merged):\n")
	var s []string
	for k, v := range records {
		s = append(s, fmt.Sprintf("%s clock in / out: %s/%s", k, v.ClockIn, v.ClockOut))
	}
	sort.Strings(s)
	for _, str := range s {
		fmt.Printf("%s\n", str)
	}

	// Output:
	// records from files/campus-a.xls
	// id: 7, name: Mary, date: 2017-07-01, clock in/out: 20:49/
	// id: 7, name: Mary, date: 2017-07-02, clock in/out: 09:34/21:24
	// id: 7, name: Mary, date: 2017-07-05, clock in/out: 17:34/
	// id: 7, name: Mary, date: 2017-07-06, clock in/out: 09:21/17:36
	// id: 7, name: Mary, date: 2017-07-07, clock in/out: 11:52/20:44
	// id: 8, name: Bob, date: 2017-07-01, clock in/out: 08:45/20:07
	// id: 8, name: Bob, date: 2017-07-02, clock in/out: 08:46/18:54
	// id: 8, name: Bob, date: 2017-07-03, clock in/out: 08:20/18:18
	// id: 8, name: Bob, date: 2017-07-04, clock in/out: 08:51/
	// id: 8, name: Bob, date: 2017-07-05, clock in/out: 08:50/
	// id: 8, name: Bob, date: 2017-07-07, clock in/out: 11:54/20:44
	// id: 9, name: Jack, date: 2017-07-01, clock in/out: 18:02/
	// id: 9, name: Jack, date: 2017-07-02, clock in/out: 16:59/
	// id: 9, name: Jack, date: 2017-07-03, clock in/out: 18:17/
	// id: 9, name: Jack, date: 2017-07-04, clock in/out: 08:58/
	// id: 9, name: Jack, date: 2017-07-07, clock in/out: 11:56/20:36
	//
	// records from files/campus-b.xls
	// id: 7, name: Mary, date: 2017-07-05, clock in/out: 09:20/
	// id: 8, name: Bob, date: 2017-07-03, clock in/out: 07:00/19:00
	// id: 9, name: Jack, date: 2017-07-04, clock in/out: 17:32/
	//
	// final records(merged):
	// Bob:2017-07-01 clock in / out: 08:45/20:07
	// Bob:2017-07-02 clock in / out: 08:46/18:54
	// Bob:2017-07-03 clock in / out: 07:00/19:00
	// Bob:2017-07-04 clock in / out: 08:51/
	// Bob:2017-07-05 clock in / out: 08:50/
	// Bob:2017-07-07 clock in / out: 11:54/20:44
	// Jack:2017-07-01 clock in / out: 18:02/
	// Jack:2017-07-02 clock in / out: 16:59/
	// Jack:2017-07-03 clock in / out: 18:17/
	// Jack:2017-07-04 clock in / out: 08:58/17:32
	// Jack:2017-07-07 clock in / out: 11:56/20:36
	// Mary:2017-07-01 clock in / out: 20:49/
	// Mary:2017-07-02 clock in / out: 09:34/21:24
	// Mary:2017-07-05 clock in / out: 09:20/17:34
	// Mary:2017-07-06 clock in / out: 09:21/17:36
	// Mary:2017-07-07 clock in / out: 11:52/20:44
}
