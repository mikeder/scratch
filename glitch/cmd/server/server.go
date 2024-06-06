package main

import (
	"log"
	"net/http"
	"time"
)

func main() {
	log.Println("Serving on 8081")

	h := http.FileServer(http.Dir(`./cmd/server/assets/`))

	logMw := func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			h.ServeHTTP(w, r)
			log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
		})
	}

	h = logMw(h)

	if err := http.ListenAndServe(`:8081`, h); err != nil {
		panic(err)
	}
}
