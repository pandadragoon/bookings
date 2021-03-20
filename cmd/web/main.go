package main

import (
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/pandadragoon/bookings/internal/config"
	"github.com/pandadragoon/bookings/internal/handlers"
	"github.com/pandadragoon/bookings/internal/render"
	"log"
	"net/http"
	"time"
)

const portNumber = ":4040"

var app config.AppConfig
var session *scs.SessionManager

func main() {

	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatalf("cannot create template cache %v", err)
	}

	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	render.NewTemplates(&app)

	fmt.Printf("Serving on port: %s \n", portNumber)

	srv := &http.Server{
		Addr:    portNumber,
		Handler: Routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}
