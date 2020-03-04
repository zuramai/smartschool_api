package controllers

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"image/png"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	gosocketio "github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"
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

func newAttendance(w http.ResponseWriter, attendance models.AttendanceBody) (models.Attendance, error) {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	todayFirst, _ := time.Parse("2006-01-02", time.Now().In(loc).Format("2006-01-02"))
	todayLast, _ := time.Parse("2006-01-02 15:04:05", time.Now().In(loc).Format("2006-01-02")+" 23:59:59")

	query := bson.M{
		"$and": []bson.M{
			bson.M{"attend_at": bson.M{"$gte": todayFirst, "$lte": todayLast}},
			bson.M{"user_id": attendance.UserID},
		},
	}

	var isAttended models.Attendance

	err := models.GetDB("main").Collection("attendances").FindOne(context.TODO(), query).Decode(&isAttended)
	if err != nil {
		fmt.Println("not attend")
	} else {
		fmt.Println("Sudah absen attend")
		return isAttended, errors.New("Sudah absen")

	}

	// cameraID, _ := strconv.ParseInt(r.FormValue("camera_id"), 10, 64)

	timeNowMs := strconv.FormatInt(makeTimestampMilli(), 10)
	timeNow, _ := time.Parse("2006-01-02 15:04:05", time.Now().In(loc).Format("2006-01-02 15:04:05"))
	photoName := timeNowMs + ".png"
	// dir, _ := os.Getwd()
	photoDir := "D:/Programmer/Code/Laravel/smart_office/public/images/attendance/" + photoName

	f, err := os.Create(photoDir)
	if err != nil {
		fmt.Println("err", err)
		return isAttended, err
	}
	unbased, erru := base64.StdEncoding.DecodeString(attendance.PhotoEncoding)
	if erru != nil {
		fmt.Println("erru", erru)
		return isAttended, erru
	}
	pngI, errp := png.Decode(bytes.NewReader(unbased))
	if errp != nil {
		fmt.Println("errp", errp)
		return isAttended, errp
	}
	errr := png.Encode(f, pngI)
	if errr != nil {
		fmt.Println("errr", errr)
		return isAttended, errr
	}

	defer f.Close()
	waktuSekarang, _ := time.Parse("15:04:05", time.Now().In(loc).Format("15:04:05"))
	waktuTerlambat, _ := time.Parse("15:04:05", "06:30:00")
	var keterangan string
	if waktuSekarang.After(waktuTerlambat) {
		keterangan = "Terlambat"
	} else {
		keterangan = "Tepat Waktu"
	}

	newAttendance := models.Attendance{
		UserID:       attendance.UserID,
		Name:         attendance.Name,
		PictureTaken: photoName,
		AttendAt:     timeNow,
		Keterangan:   keterangan,
	}
	c, err := gosocketio.Dial(
		gosocketio.GetUrl("localhost", 8765, false),
		transport.GetDefaultWebsocketTransport(),
	)
	if err != nil {
		fmt.Println("Failed to connect to socket")
		fmt.Println(err)
	}
	type AttendanceSocket struct {
		Name      string    `json:"name"`
		ImageName string    `json:"image_name"`
		Time      time.Time `json:"time"`
	}
	err2 := c.On(gosocketio.OnConnection, func(h *gosocketio.Channel) {
		log.Println("Connected")
		c.Emit("newAbsen", AttendanceSocket{newAttendance.Name, newAttendance.PictureTaken, newAttendance.AttendAt}) //close connection
	})
	err2 = c.On(gosocketio.OnDisconnection, func(h *gosocketio.Channel) {
		log.Println("Disconnected")
	})
	if err2 != nil {
		fmt.Println("Failed to connect to socket")
	}
	// jsonn, _ := json.Marshal(newAttendance)

	models.GetDB("main").Collection("attendances").InsertOne(context.TODO(), &newAttendance)

	// var addr = flag.String("addr", "localhost:8765", "http service address")

	// u := url.URL{Scheme: "ws", Host: *addr, Path: "/ws"}
	// c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)

	// errw := c.WriteMessage(websocket.TextMessage, []byte(jsonn))
	// if errw != nil {
	// 	fmt.Println(errw)
	// }
	return newAttendance, nil
}

func AttendanceV2New(w http.ResponseWriter, r *http.Request) {
	attendance := models.AttendanceBody{}
	errEncode := json.NewDecoder(r.Body).Decode(&attendance)

	if errEncode != nil {
		panic(errEncode)
	}

	result, err := newAttendance(w, attendance)
	if err != nil {
		fmt.Println("the error: ", err)
		// if sudah absen, result kosong
		respondJSON(w, 200, "Anda sudah absen hari ini", map[string]interface{}{})
	} else {
		// if belum absen, no error
		fmt.Println("the error: ", err)
		respondJSON(w, 200, "Success Attend", result)
	}
}
