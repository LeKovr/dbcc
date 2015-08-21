dbcc
====

[![Build Status](https://travis-ci.org/LeKovr/dbcc.svg?branch=master)](https://travis-ci.org/LeKovr/dbcc)  [![Build Status](https://drone.io/github.com/LeKovr/dbcc/status.png)](https://drone.io/github.com/LeKovr/dbcc/latest)

Database check & create.

Database superuser agent: Listen http port, check if given database & user exists and create them otherwise

Only postgresql database supported now.

Run
---

`$ gosu postgres ./dbcc --key=YOUR_SECRET_KEY`

Usage
-----

`curl "http://localhost:8080/?key=YOUR_SECRET_KEY&name=operator&pass=operator_pass"`


