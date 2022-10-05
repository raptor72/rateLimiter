package handlers

import (
	"fmt"
	"net/http"
)

func ReturnAccept(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message":"yes"}`))
	w.WriteHeader(http.StatusOK)
}

func ReturnDecline(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message":"no"}`))
	w.WriteHeader(http.StatusOK)
}

func ReturnBadRequest(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	response := fmt.Sprintf("{\"error\":%v}", err.Error())
	w.Write([]byte(response))
	w.WriteHeader(http.StatusBadRequest)
}
