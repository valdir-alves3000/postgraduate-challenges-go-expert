package internal_error

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func writeErrorResp(w http.ResponseWriter, code int, message string) {
	errorResponse := ErrorResponse{
		Code:    code,
		Message: message,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

func InvalidZipcodeError(w http.ResponseWriter) {
	writeErrorResp(w, http.StatusUnprocessableEntity, "invalid zipcode")
}

func CityNotFoundError(w http.ResponseWriter) {
	writeErrorResp(w, http.StatusNotFound, "can not find zipcode")
}

func TemperatureNotFoundError(w http.ResponseWriter) {
	writeErrorResp(w, http.StatusInternalServerError, "can not find temperature")
}
