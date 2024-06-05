package main

import (
	"log"
	"net/http"
)

func main() {
	log.Println("Serving on 8081")
	if err := http.ListenAndServe(`:8081`, http.FileServer(http.Dir(`./cmd/server/assets/`))); err != nil {
		panic(err)
	}
}
