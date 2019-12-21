package controllers

import (
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/zuramai/smartschool_api/models"
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
	userID, _ := strconv.ParseInt(r.FormValue("user_id"), 10, 64)
	// cameraID, _ := strconv.ParseInt(r.FormValue("camera_id"), 10, 64)
	photo, _, _ := r.FormFile("photo")

	dir, _ := os.Getwd()
	timeNowMs := strconv.FormatInt(makeTimestampMilli(), 10)
	timeNow, _ := time.Parse("2006-01-02 15:04:05", time.Now().Format("2006-01-02 15:04:05"))
	photoName := timeNowMs + ".png"
	photoDir := dir + "/assets/images/attendances/" + photoName
	f, err := os.Create(photoDir)
	if err != nil {
		log.Println(err)
		return
	}
	io.Copy(f, photo)
	defer f.Close()

	attendance := models.Attendance{
		UserID:    userID,
		PhotoName: photoName,
		AttendAt:  timeNow,
	}

	models.GetDB("main").NewRecord(&attendance)
	models.GetDB("main").Create(&attendance)

	respondJSON(w, 200, "Success Attend", attendance)
}
