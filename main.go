package main

import (
	"fmt"
	"github.com/raptor72/rateLimiter/config"
	"github.com/raptor72/rateLimiter/handlers"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func main() {
	c, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	handler := &handlers.DefaultHandler{Config: c}
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler.FallbackHandler)
	mux.HandleFunc("/ping", handler.LiveHandler)
	mux.HandleFunc("/login", handler.LoginHandler)
	log.Warning("Starting the server on http://127.0.0.1:%d", c.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", c.Port), mux))
}
