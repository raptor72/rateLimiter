package whitelists

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func NewHandler(storage Storage) *Handler {
	return &Handler{
		storage: storage,
	}
}

type Handler struct {
	storage Storage
}

func (h *Handler) GetWhiteLists(w http.ResponseWriter, r *http.Request) {
	whitelists, err := h.storage.Select()
	if err != nil {
		log.WithError(err).Error("failed to select binding_queues")
		return
	}

	for _, cidr := range whitelists {
		ipv4Addr, ipv4Net, err := net.ParseCIDR(cidr.Address)
		if err != nil {
			log.Error(err)
		}
		fmt.Println(ipv4Addr, ipv4Net)
	}

	res := whiteListResult{
		Items: whitelists,
	}

	bts, err := json.Marshal(&res)
	if err != nil {
		log.WithError(err).Error("could not marshal data response")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(bts)
}
