package health

import (
	"net/http"
)

// Handle health check request
func NewRouter() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		_, err := w.Write([]byte("OK"))
		if err != nil {
			panic(err)
		}
	}
}
