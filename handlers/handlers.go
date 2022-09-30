package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/raptor72/rateLimiter/config"
	"github.com/raptor72/rateLimiter/pkg/limiter"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

type DefaultHandler struct {
	Config *config.Config
}

func (h *DefaultHandler) FallbackHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("called url %v\n", r.URL)
	http.NotFound(w, r)
}

func (h *DefaultHandler) LiveHandler(w http.ResponseWriter, r *http.Request) {
	// livenes проба сразу отдает 200
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message":"pong"}`))
	w.WriteHeader(http.StatusOK)
	return
}

func (h *DefaultHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	limiterClient := limiter.NewClient(h.Config)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.WithError(err).Error("failed to read request body")
		return
	}

	req := LoginRequest{}
	if err := json.Unmarshal(body, &req); err != nil {
		log.WithError(err).Error("failed to unmarshal request body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	count, err := limiterClient.GetCountPattern(req.Login)
	if err != nil {
		log.WithError(err).Error("have err")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fmt.Printf("Count %d\n", *count)

	if *count > h.Config.LoginLimit {
		w.Write([]byte(`{"message":"no"}`))
		return
	}
	res := limiterClient.IncrementOrBlock(req.Login, h.Config.IpLimit, 25)
	if !res {
		w.Write([]byte(`{"message":"no"}`))
		return
	}
	w.Write([]byte(`{"message":"yes"}`))
	w.WriteHeader(http.StatusOK)
	return
}
