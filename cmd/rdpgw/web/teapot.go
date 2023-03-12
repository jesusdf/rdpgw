package web

import (
	"net/http"
)

// Test function just to check if the application is running
func Teapot(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusTeapot)
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("I'm a teapot"))
}
