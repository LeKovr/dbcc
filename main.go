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
)

// Version of this program
const Version = "1.5"

// https://elithrar.github.io/article/custom-handlers-avoiding-globals/

type appContext struct {
	key string
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

// IndexHandler - process args, check key, call DbCheckCreate
func IndexHandler(a *appContext, w http.ResponseWriter, r *http.Request) (int, error) {
	name := r.FormValue("name")

	log.Printf("**** Request: key %s name %s", r.FormValue("key"), name)
	if a.key != r.FormValue("key") {
		return http.StatusForbidden, errors.New("req: API key is absent or wrong")
	} else if r.FormValue("name") == "" {
		return http.StatusNotAcceptable, errors.New("req: name arg is required")
	}

	var status int
	db, err := DbInit()
	if err == nil {
		status, err = DbCheckCreate(db, name, r.FormValue("pass"))
		db.Close()
	}

	if err != nil {
		return http.StatusInternalServerError, err
	}

	fmt.Fprintf(w, "OK: %02b\n", status)
	return 200, nil
}

func main() {

	cli.NewApp()
	app := cli.NewApp()
	app.Name = "dbcc"
	app.Version = Version
	app.Author = "Alexey A. Kovrizhkin"
	app.Email = "lekovr+dbcc@gmail.com"
	app.Usage = "Check if database & user exists and create them if don't"

	log.Printf("dbcc v %s (%s on %s/%s)", Version, runtime.Version(), runtime.GOOS, runtime.GOARCH)

	app.Action = func(c *cli.Context) {

		key := c.String("key")
		if key == "" {
			log.Println("Error: API key does not set")
			os.Exit(1)
		}

		// Check if connect correct
		db, err := DbInit()
		checkErr(nil, err)
		db.Close()

		context := &appContext{key: key}

		r := web.New()
		r.Get(c.String("prefix"), appHandler{context, IndexHandler})

		addr := fmt.Sprintf("%s:%d", c.String("host"), c.Int("port"))
		log.Printf("Start listening at %s%s with key %s", addr, c.String("prefix"), c.String("key"))
		err = graceful.ListenAndServe(addr, r)
		checkErr(nil, err)

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
