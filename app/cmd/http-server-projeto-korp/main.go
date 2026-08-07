package main

import (
	"log"
	"net/http"

	"github.com/projeto-korp/app/internal/httpserver"
	"github.com/projeto-korp/app/internal/metrics"
)

const listenAddress = ":8080"

func main() {
	metricsComponent := metrics.New()
	server := &http.Server{
		Addr:    listenAddress,
		Handler: metricsComponent.Handler(httpserver.NewHandler()),
	}

	log.Printf("%s listening on %s", httpserver.ServiceName, listenAddress)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
