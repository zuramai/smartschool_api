package utils

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

func Message(status bool, message string) map[string]interface{} {
	return map[string]interface{}{"status": status, "message": message}
}

func Respond(w http.ResponseWriter, data map[string]interface{}) {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func Paginate(r *http.Request) (int, int, int) {
	currentPage, errC := strconv.Atoi(r.URL.Query().Get("page"))
	perPage, errP := strconv.Atoi(r.URL.Query().Get("per-page"))
	var offset int

	if errC != nil {
		offset = 0
		currentPage = 1
	}
	if errP != nil {
		perPage = 10
	}

	if currentPage > 1 {
		offset = currentPage * perPage
	}

	return currentPage, perPage, offset
}

func StartOfTheDay(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
}
