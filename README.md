# zkteco

Package zkteco parses the XLS files output from zkteco device(KQ803).

## Install
* Install [xls2csv-go](https://github.com/northbright/xls2csv-go)
  * [Install xls2csv-go package](https://github.com/northbright/xls2csv-go?tab=readme-ov-file#install-xls2csv-package)

## Build and Test Your App

```sh
CGO_CFLAGS=-I/usr/local/include CGO_LDFLAGS="-L/usr/local/lib -l xlsreader" go build
```

```sh
CGO_CFLAGS=-I/usr/local/include CGO_LDFLAGS="-L/usr/local/lib -l xlsreader" go test
```

## Documentation
* [API Reference](http://pkg.go.dev/github.com/northbright/zkteco)

#### License
* [MIT License](LICENSE)
