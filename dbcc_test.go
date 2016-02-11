package main

import (
	"database/sql"
	"fmt"
	"gopkg.in/DATA-DOG/go-sqlmock.v0"
	"testing"
)

const (
	dbccTestName = "username"
	dbccTestPass = "userpass"
)

// columns are prefixed with "o" since we used sqlstruct to generate them
var columns = []string{"o_id"}

func prepare(t *testing.T) *sql.DB {
	// open database stub
	db, err := sqlmock.New()
	if err != nil {
		t.Errorf("An error '%s' was not expected when opening a stub database connection", err)
	}
	return db
}

func process(t *testing.T, db *sql.DB, ret int) {
	status, err := DbCheckCreate(db, dbccTestName, dbccTestPass)
	if err != nil {
		t.Errorf("Expected no error, but got %s instead", err)
	}
	// db.Close() ensures that all expectations have been met
	if err = db.Close(); err != nil {
		t.Errorf("Error '%s' was not expected while closing the database", err)
	}
	if status != ret {
		t.Errorf("Expected status %d, got %d", ret, status)
	}
}

func expectUser(values string) {
	sqlmock.ExpectQuery("SELECT 1 FROM pg_roles WHERE rolname = (.+)").
		WithArgs(dbccTestName).
		WillReturnRows(sqlmock.NewRows(columns).FromCSVString(values))
}
func createUser() {
	sqlmock.ExpectExec(fmt.Sprintf("CREATE USER \"%s\" PASSWORD '%s'", dbccTestName, dbccTestPass)).
		WillReturnResult(sqlmock.NewResult(0, 1)) // no insert id, 1 affected row
}
func expectDb(values string) {
	sqlmock.ExpectQuery("SELECT 1 FROM pg_database WHERE datname = (.+)").
		WithArgs(dbccTestName).
		WillReturnRows(sqlmock.NewRows(columns).FromCSVString(values))
}
func createDb() {
	sqlmock.ExpectExec(fmt.Sprintf("CREATE DATABASE \"%s\" OWNER \"%s\"", dbccTestName, dbccTestName)).
		WillReturnResult(sqlmock.NewResult(0, 1)) // no insert id, 1 affected row
}

// will test that user and database will be created
func TestFull_DbCheckCreate(t *testing.T) {
	db := prepare(t)
	expectUser("")
	createUser()
	expectDb("")
	createDb()
	process(t, db, 3)
}

// will test that only database will be created
func TestDb_DbCheckCreate(t *testing.T) {
	db := prepare(t)
	expectUser("1")
	expectDb("")
	createDb()
	process(t, db, 2)
}

// will test that nothing will be created
func TestNothing_DbCheckCreate(t *testing.T) {
	db := prepare(t)
	expectUser("1")
	expectDb("1")
	process(t, db, 0)
}
