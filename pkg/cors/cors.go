// Package cors is used for enabling CORS in application API.
package cors

import (
	"net/http"

	"github.com/gorilla/handlers"
)

func EnableCORS(api http.Handler) http.Handler {
	headersOK := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type"})
	originsOK := handlers.AllowedOrigins([]string{"*"})
	methodsOK := handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS", "DELETE", "PUT"})
	exposedHeaders := handlers.ExposedHeaders([]string{"X-Response-Time", "X-Server-Name"})

	return handlers.CORS(headersOK, originsOK, methodsOK, exposedHeaders)(api)
}
