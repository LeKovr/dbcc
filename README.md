dbcc
====

[![Build Status][1]][2]
[![Build Status][3]][4]
[![GoRelease][5]][6]
[![GoCard][7]][8]

[1]: https://travis-ci.org/LeKovr/dbcc.svg?branch=master
[2]: https://travis-ci.org/LeKovr/dbcc
[3]: https://drone.io/github.com/LeKovr/dbcc/status.png
[4]: https://drone.io/github.com/LeKovr/dbcc/latest
[5]: https://dn-gorelease.qbox.me/gorelease-download-blue.svg
[6]: http://gorelease.herokuapp.com/LeKovr/dbcc/master
[7]: https://goreportcard.com/badge/LeKovr/dbcc
[8]: https://goreportcard.com/report/github.com/LeKovr/dbcc

[dbcc](https://github.com/LeKovr/dbcc) - Database check & create tool.

This is a database superuser agent which 

* listens http port
* gets authorized requests with `name`
* check if requested database `name` & user `name` exists and create them otherwise.

Only postgresql database supported now.

Make
----

`$ go build`

If you need cross platform build with gox, run
`$ make buildall`

Tests
-----

With mock database:
`$ go test`

With real database server:
```
# set connection vars in ENV and run
$ DBCC_TEST_DB=1 PGUSER=op go test
```

Run
---

`$ gosu postgres ./dbcc --key=YOUR_SECRET_KEY`

or

`$ APP_KEY=YOUR_SECRET_KEY gosu postgres ./dbcc`

Usage
-----

`curl "http://$DB_HOST:8080/?key=YOUR_SECRET_KEY&name=operator&pass=operator_pass[&tmpl=template]"`

Will do the following:

* if user `operator` does not exists then create it with password `operator_pass`
* if database `operator` does not exists then create it with owner `operator` and template `template` (default template1)

and return

* `OK: 00` if nothing was done
* `OK: 10` if db created (user exists)
* `OK: 11` if user & db created

License
-------

The MIT License (MIT)

Copyright (c) 2015 Alexey Kovrizhkin lekovr@gmail.com

