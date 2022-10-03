package white_lists

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net/http"
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
	whiteLists, err := h.storage.Select()
	if err != nil {
		log.WithError(err).Error("failed to select binding_queues")
		return
	}

	res := whiteListResult{
		Items: whiteLists,
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
