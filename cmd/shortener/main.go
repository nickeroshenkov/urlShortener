package main

import (
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/nickeroshenkov/urlShortener/internal/app/server"
)

func main() {
	cnf := server.NewConfig()
	cnf.Parse()
	
	srv, err := server.New(cnf)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	srv.Shutdown()
	if err != nil {
		log.Fatal(err)
	}
}
