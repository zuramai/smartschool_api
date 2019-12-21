package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/zuramai/smartschool_api/models"
	u "github.com/zuramai/smartschool_api/utils"

	"github.com/thedevsaddam/govalidator"
)

var cameras []models.Camera

func CameraIndex(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	_, perPage, offset := u.Paginate(r)
	models.GetDB("main").Where("name LIKE ?", "%"+search+"%").Offset(offset).Limit(perPage).Order("id asc").Find(&cameras)
	respondJSON(w, 200, "Success get all data camera, search: "+search, cameras)
}

func CameraStore(w http.ResponseWriter, r *http.Request) {
	rules := govalidator.MapData{
		"name":     []string{"required"},
		"location": []string{"required"},
		"note":     []string{"required"},
		"status":   []string{"bool"},
	}
	opts := govalidator.Options{
		Request:         r,     // request object
		Rules:           rules, // rules map
		RequiredDefault: true,  // all the field to be pass the rules
	}
	v := govalidator.New(opts)
	e := v.Validate()
	if len(e) > 0 {
		err := map[string]interface{}{"validation_error": e}
		w.Header().Set("Content-type", "application/json")
		json.NewEncoder(w).Encode(err)
		return
	}

	status, _ := strconv.ParseBool(r.FormValue("status"))

	camera := &models.Camera{
		Name:     r.FormValue("name"),
		Location: r.FormValue("location"),
		Note:     r.FormValue("note"),
		Status:   status,
	}
	models.GetDB("main").NewRecord(camera)
	models.GetDB("main").Create(&camera)

	respondJSON(w, 200, "Successfully create camera", camera)
	return
}

func CameraDetail(w http.ResponseWriter, r *http.Request) {

}

func CameraUpdate(w http.ResponseWriter, r *http.Request) {

}

func CameraDelete(w http.ResponseWriter, r *http.Request) {

}
