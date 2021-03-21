package main

import (
	"encoding/gob"
	"fmt"
	"github.com/TwinProduction/go-color"
	"github.com/alexedwards/scs/v2"
	"github.com/pandadragoon/bookings/internal/config"
	"github.com/pandadragoon/bookings/internal/driver"
	"github.com/pandadragoon/bookings/internal/handlers"
	"github.com/pandadragoon/bookings/internal/helpers"
	"github.com/pandadragoon/bookings/internal/models"
	"github.com/pandadragoon/bookings/internal/render"
	"log"
	"net/http"
	"os"
	"time"
)

const portNumber = ":4040"

var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

// main is the main function
func main() {
	db, err := run()
	if err != nil {
		log.Fatal(err)
	}
	defer db.SQL.Close()

	fmt.Println(fmt.Sprintf("Staring application on port %s", portNumber))

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func run() (*driver.DB, error) {
	// what am I going to put in the session
	gob.Register(models.Reservation{})
	gob.Register(models.Restriction{})
	gob.Register(models.Room{})
	gob.Register(models.User{})

	// change this to true when in production
	app.InProduction = false

	infoLog = log.New(os.Stdout, color.Cyan+"INFO\t"+color.Reset, log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog = log.New(os.Stdout, color.Red+"ERROR\t"+color.Reset, log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	// set up the session
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	fmt.Println("connecting to db...")
	db, err := driver.ConnectSQL("host=localhost port=5432 dbname=bookings user=postgres password=")
	if err != nil {
		log.Fatalf("Cannot connect to datase, dying... %v", err)
	}

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
		return nil, err
	}

	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)

	render.NewRenderer(&app)
	helpers.NewHelpers(&app)
	return db, nil
}
