package main

/*

Run as
  gosu postgres ./dbcc --key=YOUR_SECRET_KEY
Call as
  curl "http://localhost:8080/?key=YOUR_SECRET_KEY&name=operator&pass=operator_pass"

*/

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"

	"github.com/codegangsta/cli"

	"github.com/zenazn/goji/graceful"
	"github.com/zenazn/goji/web"

	"database/sql"
	"github.com/lib/pq"
)

const Version = "1.1"

// https://elithrar.github.io/article/custom-handlers-avoiding-globals/

type appContext struct {
	config *cli.Context
	db     *sql.DB
}

type appHandler struct {
	*appContext
	h func(*appContext, http.ResponseWriter, *http.Request) (int, error)
}

// Our ServeHTTP method is mostly the same, and also has the ability to
// access our *appContext's fields (templates, loggers, etc.) as well.
func (ah appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Updated to pass ah.appContext as a parameter to our handler type.
	status, err := ah.h(ah.appContext, w, r)
	if err != nil {
		log.Printf("HTTP %d: %q", status, err)
		switch status {
		case http.StatusNotFound:
			http.NotFound(w, r)
			// And if we wanted a friendlier error page, we can
			// now leverage our context instance - e.g.
			// err := ah.renderTemplate(w, "http_404.tmpl", nil)
		case http.StatusInternalServerError:
			http.Error(w, http.StatusText(status), status)
		default:
			http.Error(w, http.StatusText(status), status)
		}
	}
}

func IndexHandler(a *appContext, w http.ResponseWriter, r *http.Request) (int, error) {
	name := r.FormValue("name")
	nameQuoted := pq.QuoteIdentifier(name)

	log.Printf("**** Request: key %s name %s", r.FormValue("key"), name)
	if a.config.String("key") != r.FormValue("key") {
		return http.StatusForbidden, errors.New("req: API key is absent or wrong")
	} else if r.FormValue("name") == "" {
		return http.StatusNotAcceptable, errors.New("req: name arg is required")
	}
	db := a.db
	err := db.Ping()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	var ret = 0
	rows, err := db.Query("SELECT 1 FROM pg_roles WHERE rolname = $1", name)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	if rows.Next() {
		log.Printf("User %s already exists", name)
	} else {
		//    rows, err := db.Query("create user " + nameQuoted + " with password $1", r.FormValue("pass"))
		_, err := db.Exec(fmt.Sprintf("CREATE USER %s PASSWORD '", nameQuoted) + r.FormValue("pass") + "'")
		if err != nil {
			return http.StatusInternalServerError, err
		}
		log.Printf("User %s created", name)
		ret += 1
	}

	rows, err = db.Query("SELECT 1 FROM pg_database WHERE datname = $1", name)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	if rows.Next() {
		log.Printf("Database %s already exists", name)
	} else {
		_, err := db.Exec(fmt.Sprintf("CREATE DATABASE %s OWNER %s", nameQuoted, nameQuoted))
		if err != nil {
			return http.StatusInternalServerError, err
		}
		log.Printf("Database %s created", name)
		ret += 2
	}
	fmt.Fprintf(w, "OK: %02b\n", ret)
	return 200, nil
}

func main() {

	app := cli.NewApp()
	app.Name = "DBcc"
	app.Version = Version
	app.Author = "Alexey A. Kovrizhkin"
	app.Email = "lekovr+dbcc@gmail.com"
	app.Usage = "Check if database & user exists and create them if don't"

	log.Printf("%s version: %s (%s on %s/%s; %s) compiled at %s", os.Args[0], Version, runtime.Version(), runtime.GOOS, runtime.GOARCH, runtime.Compiler, app.Compiled)
	app.Action = func(c *cli.Context) {

		dbinfo := fmt.Sprintf("sslmode=disable")
		db, err := sql.Open("postgres", dbinfo)
		checkErr(nil, err)
		defer db.Close()

		addr := fmt.Sprintf("%s:%d", c.String("host"), c.Int("port"))
		context := &appContext{config: c, db: db}

		log.Printf("Start listening at %s%s with key %s", addr, c.String("prefix"), c.String("key"))

		r := web.New()
		// We pass an instance to our context pointer, and our handler.
		r.Get(c.String("prefix"), appHandler{context, IndexHandler})
		graceful.ListenAndServe(addr, r)

	}
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "key, k",
			Usage:  "API key required in requests",
			EnvVar: "APP_KEY",
		},
		cli.StringFlag{
			Name:   "port, p",
			Value:  "8080",
			Usage:  "Port listen to",
			EnvVar: "APP_PORT,PORT",
		},
		cli.StringFlag{
			Name:   "host",
			Value:  "",
			Usage:  "IP listen to",
			EnvVar: "APP_HOST,HOST",
		},
		cli.StringFlag{
			Name:   "prefix",
			Value:  "/",
			Usage:  "URL prefix",
			EnvVar: "APP_PREFIX,PREFIX",
		},
	}
	app.Run(os.Args)

}

func checkErr(w http.ResponseWriter, err error) {
	if err != nil {
		if w != nil {
			status := http.StatusInternalServerError
			http.Error(w, http.StatusText(status), status)
		}
		panic(err)
	}
}

/*
func (h *handler) ServeHTTP(
    w http.ResponseWriter,
    r *http.Request,
) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    enc := json.NewEncoder(w)
    if err := enc.Encode(&MyResponse{}); nil != err {
        fmt.Fprintf(w, `{"error":"%s"}`, err)
    }
}
*/
