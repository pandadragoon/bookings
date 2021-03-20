package helpers

import (
	"fmt"
	"github.com/TwinProduction/go-color"
	"github.com/pandadragoon/bookings/internal/config"
	"net/http"
	"runtime/debug"
)

var app *config.AppConfig

// NewHelpers sets up app config for helpers
func NewHelpers(a *config.AppConfig) {
	app = a
}

// ClientError
func ClientError(w http.ResponseWriter, status int) {
	app.InfoLog.Println(color.Cyan+"Client error with status of"+color.Reset, status)
	http.Error(w, http.StatusText(status), status)
}

// ServerError
func ServerError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.ErrorLog.Println(color.Red + trace + color.Reset)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}
