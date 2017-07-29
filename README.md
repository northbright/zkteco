# zkteco

Package zkteco parses the XLS files outputted from zkteco device(KQ803) and update the employee attendance data in [Redis](https://redis.io).

#### Data in Redis
* The attendance data of each emplyee are stored in a Redis Hash.
* Key is `kaoqin:EmplyeeName`. e.g. `kaoqin:Jack`.

        redis-cli --raw
        127.0.0.1:6379> keys *
        kaoqin:Jack
        kaoqin:Bob
        kaoqin:Mary 

* Field of the Redis Hash has two type.
  * `date:in` represents the clock in time of the date. Date format is YYYY-MM-DD.
  * `date:out` represents the clock out time of the date. 

          127.0.0.1:6379> HGETALL kaoqin:Jack
          2017-07-01:in
          18:02
          2017-07-02:in
          16:59
          2017-07-03:in
          18:17
          2017-07-04:in
          08:58
          2017-07-07:in
          11:56
          2017-07-07:out
          20:36

  * If there's only one clock time for a date, it will be recognized as clock in time.

* Value of each field is the clock in / out time. Time format: HH:MM.

#### Examples

        // Open a DB by given Redis address and password.
        db := zkteco.Open(":6379", "")

        // Get absolute path of example attendance xls file.
        f := "/home/xx/campus-a.xls"

        // Update attendance.
        db.UpdateAttendance(f)

#### Documentation
* [API Reference](http://godoc.org/github.com/northbright/zkteco)

#### Requirements
* zkteco requires [libxls](http://libxls.sourceforge.net/) and [xls2csv-go](https://github.com/northbright/xls2csv-go)
* You may find installation guide in [xls2csv-go](https://github.com/northbright/xls2csv-go).
* You may also find instructions in [.travis.yml](.travis.yml) to build or test your apps.

#### License
* [MIT License](LICENSE)
