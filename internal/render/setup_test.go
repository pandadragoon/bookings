package render

import (
	"encoding/gob"
	"github.com/TwinProduction/go-color"
	"github.com/alexedwards/scs/v2"
	"github.com/pandadragoon/bookings/internal/config"
	"github.com/pandadragoon/bookings/internal/models"
	"log"
	"net/http"
	"os"
	"testing"
	"time"
)

var session *scs.SessionManager
var testApp config.AppConfig

func TestMain(m *testing.M) {

	gob.Register(models.Reservation{})

	// change this to true when in production
	testApp.InProduction = false

	infoLog := log.New(os.Stdout, color.Cyan+"INFO\t"+color.Reset, log.Ldate|log.Ltime)
	testApp.InfoLog = infoLog

	errorLog := log.New(os.Stdout, color.Red+"ERROR\t"+color.Reset, log.Ldate|log.Ltime|log.Lshortfile)
	testApp.ErrorLog = errorLog

	// set up the session
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = false

	testApp.Session = session

	app = &testApp

	os.Exit(m.Run())
}

type myWriter struct{}

func (tw *myWriter) Header() http.Header {
	var h http.Header
	return h
}

func (tw *myWriter) WriteHeader(i int) {}

func (tw *myWriter) Write(b []byte) (int, error) {
	length := len(b)
	return length, nil
}
