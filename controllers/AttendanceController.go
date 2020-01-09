package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"flag"

	"github.com/gorilla/websocket"
	"github.com/zuramai/smartschool_api/models"
	"go.mongodb.org/mongo-driver/bson"
)

var attendance models.Attendance

func AttendanceStudentIndex(w http.ResponseWriter, r *http.Request) {

}

func AttendanceStudentLateToday(w http.ResponseWriter, r *http.Request) {

}

func AttendanceStudentReport(w http.ResponseWriter, r *http.Request) {

}

func AttendanceTeacherIndex(w http.ResponseWriter, r *http.Request) {

}

func AttendanceTeacherLateToday(w http.ResponseWriter, r *http.Request) {

}

func AttendanceTeacherReport(w http.ResponseWriter, r *http.Request) {

}

func AttendanceV2(w http.ResponseWriter, r *http.Request) {

}

func AttendanceV2New(w http.ResponseWriter, r *http.Request) {

	userID := r.FormValue("user_id")
	photo, _, _ := r.FormFile("photo")

	todayFirst, _ := time.Parse("2006-01-02", time.Now().Format("2006-01-02"))
	todayLast, _ := time.Parse("2006-01-02 15:04:05", time.Now().Format("2006-01-02")+" 23:59:59")

	query := bson.M{
		"$and": []bson.M{
			bson.M{"attend_at": bson.M{"$gte": todayFirst, "$lte": todayLast}},
			bson.M{"user_id": userID},
		},
	}

	var isAttended models.Attendance

	err := models.GetDB("main").Collection("attendances").FindOne(context.TODO(), query).Decode(&isAttended)

	if err == nil {
		fmt.Println(err)
		respondJSON(w, 403, "Anda sudah absen hari ini", map[string]interface{}{})
		return
	}

	// cameraID, _ := strconv.ParseInt(r.FormValue("camera_id"), 10, 64)

	timeNowMs := strconv.FormatInt(makeTimestampMilli(), 10)
	timeNow, _ := time.Parse("2006-01-02 15:04:05", time.Now().Format("2006-01-02 15:04:05"))
	photoName := timeNowMs + ".png"
	dir, _ := os.Getwd()
	photoDir := dir + "/assets/images/attendances/" + photoName
	f, err := os.Create(photoDir)
	if err != nil {
		log.Println(err)
		return
	}
	io.Copy(f, photo)
	defer f.Close()

	waktuSekarang, _ := time.Parse("15:04:05", time.Now().Format("15:04:05"))
	waktuTerlambat, _ := time.Parse("15:04:05", "06:30:00")
	var keterangan string
	if waktuSekarang.After(waktuTerlambat) {
		keterangan = "Terlambat"
	} else {
		keterangan = "Tepat Waktu"
	}

	attendance := models.Attendance{
		UserID:       userID,
		PictureTaken: photoName,
		AttendAt:     timeNow,
		Keterangan:   keterangan,
	}
	jsonn, _ := json.Marshal(attendance)

	models.GetDB("main").Collection("attendances").InsertOne(context.TODO(), &attendance)

	var addr = flag.String("addr", "localhost:8765", "http service address")

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/ws"}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)

	errw := c.WriteMessage(websocket.TextMessage, []byte(jsonn))
	if errw != nil {
		fmt.Println(errw)
	}
	respondJSON(w, 200, "Success Attend", attendance)
}
