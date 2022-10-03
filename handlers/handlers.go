package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/raptor72/rateLimiter/config"
	"github.com/raptor72/rateLimiter/pkg/limiter"
	log "github.com/sirupsen/logrus"
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
}

func (u UnionRequest) GetField(tag string) (string, error) {
	err := errors.New("do not have valid field")
	switch {
	case u.Login != "" && tag == "Login":
		return u.Login, nil
	case u.Password != "" && tag == "Password":
		return u.Password, nil
	case u.Ip != "" && tag == "Ip":
		return u.Ip, nil
	default:
		return "", err
	}
}

func BaseHandler(config *config.Config, w http.ResponseWriter, r *http.Request, limit config.CountLimit, couldown config.CoulDownTime, tag string) {
	limiterClient := limiter.NewClient(config)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.WithError(err).Error("failed to read request body")
		return
	}

	var req UnionRequest

	if err := json.Unmarshal(body, &req); err != nil {
		log.WithError(err).Error("failed to unmarshal request body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	field, err := req.GetField(tag)
	if err != nil {
		log.WithError(err).Error("Wrong json")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// fmt.Printf("%+v, %T\n", req, req)

	count, err := limiterClient.GetCountPattern(field)
	if err != nil {
		log.WithError(err).Error("have err")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fmt.Printf("Count %d\n", *count)

	if *count > limit.Count {
		w.Write([]byte(`{"message":"no"}`))
		return
	}
	res := limiterClient.IncrementOrBlock(field, limit.Count, time.Duration(couldown.SecLimit))
	if !res {
		w.Write([]byte(`{"message":"no"}`))
		return
	}
	w.Write([]byte(`{"message":"yes"}`))
	w.WriteHeader(http.StatusOK)
}

func (h *DefaultHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	BaseHandler(h.Config, w, r, h.Config.LoginLimit, h.Config.IpCouldown, "Login")
}

func (h *DefaultHandler) PasswordHandler(w http.ResponseWriter, r *http.Request) {
	BaseHandler(h.Config, w, r, h.Config.PasswordLimit, h.Config.PasswordCouldown, "Password")
}

func (h *DefaultHandler) IPHandler(w http.ResponseWriter, r *http.Request) {
	BaseHandler(h.Config, w, r, h.Config.IpLimit, h.Config.IpCouldown, "Ip")
}

func (h *DefaultHandler) WhiteListHandler(w http.ResponseWriter, r *http.Request) {
	db, err := config.NewDB(h.Config)
	if err != nil {
		log.Errorln("new DB error:", err)
		return
	}

	WhitelistHandler, err := injectWhiteLists(db)
	if err != nil {
		log.WithError(err).Fatal("failed to inject server disk handler")
	}
	WhitelistHandler.GetWhiteLists(w, r)
}
