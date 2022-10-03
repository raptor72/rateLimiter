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
		fmt.Println("Bad config")
		log.Fatal(err)
	}

	handler := &handlers.DefaultHandler{Config: c}
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler.FallbackHandler)
	mux.HandleFunc("/ping", handler.LiveHandler)
	mux.HandleFunc("/login", handler.LoginHandler)
	mux.HandleFunc("/password", handler.PasswordHandler)
	mux.HandleFunc("/ip", handler.IPHandler)
	mux.HandleFunc("/white_list", handler.WhiteListHandler)
	log.Warningf("Starting the server on http://127.0.0.1:%d", c.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", c.Port), mux))
}
