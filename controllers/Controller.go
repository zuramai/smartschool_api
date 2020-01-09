package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/zuramai/smartschool_api/models"
)

func respondJSON(w http.ResponseWriter, status int, message string, data interface{}) {
	var payload models.Payload
	if status == 200 {
		payload.Status = true
	} else {
		payload.Status = false
	}
	payload.Message = message
	payload.Data = data

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func respondErrorValidationJSON(w http.ResponseWriter, status int, message string, data map[string]interface{}) {
	var payload models.ErrorValidation
	payload.Message = "Error"
	payload.Errors = data

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func unixMilli(t time.Time) int64 {
	return t.Round(time.Millisecond).UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond))
}

func makeTimestampMilli() int64 {
	return unixMilli(time.Now())
}
