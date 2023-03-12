package main

import (
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/nickeroshenkov/urlShortener/internal/app/config"
	"github.com/nickeroshenkov/urlShortener/internal/app/server"
)

func main() {
	var c config.ServerConfig
	c.SetDefaults()
	c.LoadFlagsConditional() // Prioritize flags over the default values
	c.LoadEVarsConditional() // Prioritize environment variables over the flags

	srv, err := server.New(&c)
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
