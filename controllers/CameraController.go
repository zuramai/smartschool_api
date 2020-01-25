package controllers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/zuramai/smartschool_api/models"
	"go.mongodb.org/mongo-driver/bson"
)

var cameras []models.Camera

func CameraIndex(w http.ResponseWriter, r *http.Request) {
	var camera models.Camera
	var cameras []models.Camera
	cursor, err := models.GetDB("main").Collection("cameras").Find(context.TODO(), bson.M{})

	if err != nil {
		fmt.Println(err)
	}

	for cursor.Next(context.TODO()) {
		cursor.Decode(&camera)
		cameras = append(cameras, camera)
		camera = models.Camera{}
	}

	respondJSON(w, 200, "Get All Data Camera", cameras)
	return
}

func CameraStore(w http.ResponseWriter, r *http.Request) {
	// rules := govalidator.MapData{
	// 	"name":     []string{"required"},
	// 	"location": []string{"required"},
	// 	"note":     []string{"required"},
	// 	"status":   []string{"bool"},
	// }
	// opts := govalidator.Options{
	// 	Request:         r,     // request object
	// 	Rules:           rules, // rules map
	// 	RequiredDefault: true,  // all the field to be pass the rules
	// }
	// v := govalidator.New(opts)
	// e := v.Validate()
	// if len(e) > 0 {
	// 	err := map[string]interface{}{"validation_error": e}
	// 	w.Header().Set("Content-type", "application/json")
	// 	json.NewEncoder(w).Encode(err)
	// 	return
	// }

	// status, _ := strconv.ParseBool(r.FormValue("status"))

	// camera := &models.Camera{
	// 	Name:     r.FormValue("name"),
	// 	Location: r.FormValue("location"),
	// 	Note:     r.FormValue("note"),
	// 	Status:   status,
	// }
	// models.GetDB("main").NewRecord(camera)
	// models.GetDB("main").Create(&camera)

	// respondJSON(w, 200, "Successfully create camera", camera)
	// return
}

func CameraDetail(w http.ResponseWriter, r *http.Request) {

}

func CameraUpdate(w http.ResponseWriter, r *http.Request) {

}

func CameraDelete(w http.ResponseWriter, r *http.Request) {

}
