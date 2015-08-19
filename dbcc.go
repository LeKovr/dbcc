package main

import (
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	"log"
)

// Database check and create objects
func DbCheckCreate(db *sql.DB, name, pass string) (ret int, err error) {
	ret = 0
	if err = db.Ping(); err != nil {
		return
	}

	var rows *sql.Rows
	rows, err = db.Query("SELECT 1 FROM pg_roles WHERE rolname = $1", name)
	if err != nil {
		return
	}

	nameQuoted := pq.QuoteIdentifier(name)
	if rows.Next() {
		log.Printf("User %s already exists", name)
	} else {
		//    rows, err := db.Query("create user " + nameQuoted + " with password $1", r.FormValue("pass"))
		_, err = db.Exec(fmt.Sprintf("CREATE USER %s PASSWORD '", nameQuoted) + pass + "'")
		if err != nil {
			return
		}
		log.Printf("User %s created", name)
		ret += 1
	}

	rows, err = db.Query("SELECT 1 FROM pg_database WHERE datname = $1", name)
	if err != nil {
		return
	}
	if rows.Next() {
		log.Printf("Database %s already exists", name)
	} else {
		_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s OWNER %s", nameQuoted, nameQuoted))
		if err != nil {
			return
		}
		log.Printf("Database %s created", name)
		ret += 2
	}
	return //  ret, nil
}
