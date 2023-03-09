package main

import (
	"log"

	"github.com/nickeroshenkov/urlShortener/internal/app/config"
	"github.com/nickeroshenkov/urlShortener/internal/app/server"
)

func main() {
	c := config.Config{}
	c.SetDefaults()
	c.LoadFlagsConditional() // Prioritize flags over the default values
	c.LoadEVarsConditional() // Prioritize environment variables over the flags

	if err := server.Run(&c); err != nil {
		log.Fatal(err)
	}
}
