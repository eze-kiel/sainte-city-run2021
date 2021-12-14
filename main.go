package main

import (
	"net/http"
	"time"

	"github.com/eze-kiel/sort-run-time/handlers"
	log "github.com/sirupsen/logrus"
)

func main() {
	srv := &http.Server{
		Addr:         ":8080",
		Handler:      handlers.HandleFunc(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	log.Info("[DEV] Server is starting, wish me luck boys")
	log.Info("http://0.0.0.0" + srv.Addr)
	log.Println(srv.ListenAndServe())
}
