package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/zenazn/goji/web"

	"database/sql"
	_ "github.com/lib/pq"
)

const (
	testKey    = "simpletest"
	testUriFmt = "/?key=%s&name=dbcctestuser%d"
)

var httpTests = []struct {
	n        int // input
	dbUsed   bool
	url      string
	expected string
}{
	{1, false, "/", "403 - Forbidden"},                             // API key is absent or wrong (absent)
	{2, false, "/?key=XX", "403 - Forbidden"},                      // API key is absent or wrong (wrong)
	{3, false, "/?key=" + testKey, "406 - Not Acceptable"},         // name arg is required
	{4, true, fmt.Sprintf(testUriFmt, testKey, 0), "200 - OK: 11"}, // no user, no db
	{5, true, fmt.Sprintf(testUriFmt, testKey, 1), "200 - OK: 10"}, // user exists, no db
	{6, true, fmt.Sprintf(testUriFmt, testKey, 2), "200 - OK: 00"}, // user & db exists
}

var sqlBefore = []string{
	"CREATE USER dbcctestuser1",
	"CREATE USER dbcctestuser2",
	"CREATE DATABASE dbcctestuser2 OWNER dbcctestuser2",
}

var sqlAfter = []string{
	"DROP DATABASE dbcctestuser0", "DROP USER dbcctestuser0",
	"DROP DATABASE dbcctestuser1", "DROP USER dbcctestuser1",
	"DROP DATABASE dbcctestuser2", "DROP USER dbcctestuser2",
}

func checkTestErr(t *testing.T, err error) {
	if err != nil {
		t.Errorf("An error '%s' was not expected ", err)
	}
}
func Test_main(t *testing.T) {

	dbinfo := fmt.Sprintf("sslmode=disable")
	db, err := sql.Open("postgres", dbinfo)
	checkTestErr(t, err)
	defer db.Close()

	context := &appContext{key: testKey}

	s := web.New()
	s.Get("/", appHandler{context, IndexHandler})

	dbOn := os.Getenv("DBCC_TEST_DB") != ""
	if dbOn {
		for _, sql := range sqlBefore {
			_, err := db.Exec(sql)
			checkTestErr(t, err)
		}
	}

	for _, tt := range httpTests {

		if tt.dbUsed && !dbOn {
			t.Skip("Skipping db tests: DBCC_TEST_DB not set", tt.n)
		}

		w := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "http://example.com"+tt.url, nil)
		checkTestErr(t, err)
		s.ServeHTTP(w, req)
		resp := fmt.Sprintf("%d - %s", w.Code, w.Body.String())

		if resp != tt.expected+"\n" {
			t.Errorf("Http(%d): expected '%s', got '%s'", tt.n, tt.expected, resp)
		}
	}

	if dbOn {
		for _, sql := range sqlAfter {
			_, err := db.Exec(sql)
			checkTestErr(t, err)
		}
	}
}
