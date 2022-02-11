package main

import (
	"github.com/dragun-igor/fool/internal/handlers"
	"log"
	"net/http"
)

func main() {
	mux := routes()

	log.Println("starting channel listener")
	go handlers.ListenToWsChannel()

	log.Println("starting web server on port 8081")

	_ = http.ListenAndServe(":8081", mux)
}
