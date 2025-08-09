package zkteco_test

import (
	"fmt"
	"log"

	"github.com/northbright/zkteco"
)

func Example() {
	f := "files/campus-a.xls"

	// Parse the XLS and get attendance data.
	s, err := zkteco.Parse(f)
	if err != nil {
		log.Printf("zkteco.Parse() error: %v", err)
		return
	}
	log.Printf("zkteco.Parse() OK")

	for _, a := range s {
		fmt.Printf("ID: %s, Name: %s, Date: %s, Clock In/Out: %s/%s\n", a.ID, a.Name, a.Date, a.ClockIn, a.ClockOut)
	}

	// Output:
	// ID: 7, Name: Mary, Date: 2017-07-01, Clock In/Out: 20:49/
	// ID: 7, Name: Mary, Date: 2017-07-02, Clock In/Out: 09:34/21:24
	// ID: 7, Name: Mary, Date: 2017-07-05, Clock In/Out: 17:34/
	// ID: 7, Name: Mary, Date: 2017-07-06, Clock In/Out: 09:21/17:36
	// ID: 7, Name: Mary, Date: 2017-07-07, Clock In/Out: 11:52/20:44
	// ID: 8, Name: Bob, Date: 2017-07-01, Clock In/Out: 08:45/20:07
	// ID: 8, Name: Bob, Date: 2017-07-02, Clock In/Out: 08:46/18:54
	// ID: 8, Name: Bob, Date: 2017-07-03, Clock In/Out: 08:20/18:18
	// ID: 8, Name: Bob, Date: 2017-07-04, Clock In/Out: 08:51/
	// ID: 8, Name: Bob, Date: 2017-07-05, Clock In/Out: 08:50/
	// ID: 8, Name: Bob, Date: 2017-07-07, Clock In/Out: 11:54/20:44
	// ID: 9, Name: Jack, Date: 2017-07-01, Clock In/Out: 18:02/
	// ID: 9, Name: Jack, Date: 2017-07-02, Clock In/Out: 16:59/
	// ID: 9, Name: Jack, Date: 2017-07-03, Clock In/Out: 18:17/
	// ID: 9, Name: Jack, Date: 2017-07-04, Clock In/Out: 08:58/
	// ID: 9, Name: Jack, Date: 2017-07-07, Clock In/Out: 11:56/20:36
}
