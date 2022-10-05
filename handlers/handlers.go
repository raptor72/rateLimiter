package handlers

import (
	"encoding/json"
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
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message":"pong"}`))
	w.WriteHeader(http.StatusOK)
}

func (u UnionRequest) GetField(tag string) (string, error) {
	err := fmt.Errorf("do not get valid field for tag %v", tag)
	switch {
	case u.Login != "" && tag == "Login":
		return u.Login, nil
	case u.Password != "" && tag == "Password":
		return u.Password, nil
	case u.IP != "" && tag == "Ip":
		return u.IP, nil
	default:
		return "", err
	}
}

func BaseHandler(cfg *config.Config, w http.ResponseWriter, r *http.Request,
	lim config.CountLimit, timeout config.CoolDownTime, tag string,
) {
	limiterClient := limiter.NewClient(cfg)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.WithError(err).Error("failed to read request body")
		ReturnBadRequest(w, err)
		return
	}

	var req UnionRequest

	if err := json.Unmarshal(body, &req); err != nil {
		log.WithError(err).Error("failed to unmarshal request body")
		ReturnBadRequest(w, err)
		return
	}

	field, err := req.GetField(tag)
	if err != nil {
		log.WithError(err).Error("Wrong input body json")
		ReturnBadRequest(w, err)
		return
	}

	if tag == "Ip" {
		accept, decline, err := CheckWhiteBlack(cfg, req)
		if err != nil {
			ReturnBadRequest(w, fmt.Errorf("can not parse ip address from %v", req.IP))
			return
		}
		if accept {
			ReturnAccept(w)
			return
		}
		if decline {
			ReturnDecline(w)
			return
		}
	}

	count, err := limiterClient.GetCountPattern(field)
	if err != nil {
		log.WithError(err).Error("have GetCountPattern err")
		ReturnBadRequest(w, err)
		return
	}

	log.Infof("Count %d\n", *count)

	if *count > lim.Count {
		ReturnDecline(w)
		return
	}
	res := limiterClient.IncrementOrBlock(field, lim.Count, time.Duration(timeout.SecLimit))
	if !res {
		ReturnDecline(w)
		return
	}
	ReturnAccept(w)
}

func (h *DefaultHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	BaseHandler(h.Config, w, r, h.Config.LoginLimit, h.Config.IPCoolDown, "Login")
}

func (h *DefaultHandler) PasswordHandler(w http.ResponseWriter, r *http.Request) {
	BaseHandler(h.Config, w, r, h.Config.PasswordLimit, h.Config.PasswordCoolDown, "Password")
}

func (h *DefaultHandler) IPHandler(w http.ResponseWriter, r *http.Request) {
	BaseHandler(h.Config, w, r, h.Config.IPLimit, h.Config.IPCoolDown, "Ip")
}

func (h *DefaultHandler) WhiteListHandler(w http.ResponseWriter, r *http.Request) {
	db, err := h.Config.NewDB()
	if err != nil {
		log.Errorln("new DB error:", err)
		return
	}

	WhitelistHandler := injectWhiteLists(db)
	if err != nil {
		log.WithError(err).Fatal("failed to inject server disk handler")
	}
	WhitelistHandler.GetWhiteLists(w, r)
}
