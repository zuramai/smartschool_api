package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/zuramai/smartschool_api/models"
)

func respondJSON(w http.ResponseWriter, status int, message string, data interface{}) {
	var payload models.Payload
	payload.Status = true
	payload.Message = message
	payload.Data = data

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
