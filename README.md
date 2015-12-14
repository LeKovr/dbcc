dbcc
====

[![Build Status](https://travis-ci.org/LeKovr/dbcc.svg?branch=master)](https://travis-ci.org/LeKovr/dbcc) 
 [![Build Status](https://drone.io/github.com/LeKovr/dbcc/status.png)](https://drone.io/github.com/LeKovr/dbcc/latest)
 [![gorelease](https://dn-gorelease.qbox.me/gorelease-download-blue.svg)](http://gorelease.herokuapp.com/LeKovr/dbcc/master)

[dbcc](https://github.com/LeKovr/dbcc) - Database check & create tool.

This is a database superuser agent which listens http port, check if requested database & user exists and create them otherwise.

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

`curl "http://$DB_HOST:8080/?key=YOUR_SECRET_KEY&name=operator&pass=operator_pass"`

Will do the following:

* if user `operator` does not exists then create it with password `operator_pass`
* if database `operator` does not exists then create it with owner `operator`

and return

* `OK: 00` if nothing was done
* `OK: 10` if user created
* `OK: 01` if db created
* `OK: 11` if user & db created

License
-------

The MIT License (MIT)

Copyright (c) 2015 Alexey Kovrizhkin lekovr@gmail.com

