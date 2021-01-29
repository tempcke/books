package rest

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func decodeRequestData(w http.ResponseWriter, body io.Reader, data interface{}) error {
	err := json.NewDecoder(body).Decode(&data)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, "Request body was not valid json")
		return err
	}
	return nil
}

func jsonResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Println("json.Encode response failed: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func errorResponse(w http.ResponseWriter, code int, msg string) {
	w.WriteHeader(code)
	jsonResponse(w, ErrorResponse{Error: msg})
}
