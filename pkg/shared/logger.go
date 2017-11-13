package shared

import (
	"log"
	"net/http"
	"time"
)

// Logger is a generic logging middleware
func Logger(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	t1 := time.Now()
	next(w, r)
	t2 := time.Now()
	log.Printf("[%s] %q %v", r.Method, r.URL.String(), t2.Sub(t1))
}
