dbcc
====

Database check & create.

Database superuser agent: Listen http port, check if given database & user exists and create them otherwise

Only postgresql database supported now.

Run
---

`$ gosu postgres ./dbcc --key=YOUR_SECRET_KEY`

Usage
-----

`curl "http://localhost:8080/?key=YOUR_SECRET_KEY&user=operator&pass=operator_pass"`


