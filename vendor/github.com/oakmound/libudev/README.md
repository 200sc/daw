# libudev
Golang native implementation Udev library

[![Build Status](https://travis-ci.org/oakmound/libudev.svg?branch=master)](https://travis-ci.org/oakmound/libudev)
[![Coverage Status](https://coveralls.io/repos/github/oakmound/libudev/badge.svg?branch=master)](https://coveralls.io/github/oakmound/libudev?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/oakmound/libudev)](https://goreportcard.com/report/github.com/oakmound/libudev)
[![GoDoc](https://godoc.org/github.com/oakmound/libudev?status.svg)](https://godoc.org/github.com/oakmound/libudev)
[![GitHub release](https://img.shields.io/github/release/oakmound/libudev.svg)](https://github.com/oakmound/libudev/releases)


Installation
------------
    go get github.com/oakmound/libudev

Usage
-----

### Scanning devices
```go
sc := libudev.NewScanner()
err, devices := s.ScanDevices()
```

### Filtering devices
```go
m := matcher.NewMatcher()
m.SetStrategy(matcher.StrategyOr)
m.AddRule(matcher.NewRuleAttr("dev", "189:133"))
m.AddRule(matcher.NewRuleEnv("DEVNAME", "usb/lp0"))

filteredDevices := m.Match(devices)
```

### Getting parent device
```go
if device.Parent != nil {
    fmt.Printf("%s\n", device.Parent.Devpath)
}
```

### Getting children devices
```go
fmt.Printf("Count children devices %d\n", len(device.Children))
```

Features
--------
* 100% Native code
* Without external dependencies
* Code is covered by tests

Requirements
------------

* Need at least `go1.10` or newer.

Documentation
-------------

You can read package documentation [here](http:godoc.org/github.com/oakmound/libudev) or read tests.

Testing
-------
Unit-tests:
```bash
go test -race -v ./...
```

Contributing
------------
* Fork
* Write code
* Run unit test: `go test -v ./...`
* Run go vet: `go vet -v ./...`
* Run go fmt: `go fmt ./...`
* Commit changes
* Create pull-request
