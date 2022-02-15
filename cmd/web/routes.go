package main

import (
	"github.com/bmizerany/pat"
	"github.com/dragun-igor/fool/internal/handlers"
	"net/http"
)

func routes() http.Handler {
	mux := pat.New()

	mux.Get("/ws", http.HandlerFunc(handlers.WsEndpoint))

	return mux
}
