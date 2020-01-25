package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	// "github.com/jinzhu/gorm"

	"github.com/zuramai/smartschool_api/app"
	"github.com/zuramai/smartschool_api/controllers"
	// "github.com/zuramai/smartschool_api/models"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "socket-test.html")
		fmt.Println("home")
	})
	// http.Handle("/", http.FileServer(http.Dir("./assets")))
	APP_PORT := os.Getenv("APP_PORT")
	if APP_PORT == "" {
		APP_PORT = "8088"
	}

	apiV1 := router.PathPrefix("/api/v1").Subrouter()
	apiV1.Use(app.JwtAuthentication)
	apiV2 := router.PathPrefix("/api/v2").Subrouter()

	v1Auth := apiV1.PathPrefix("/auth").Subrouter()
	v1Auth.HandleFunc("/login", controllers.AuthLogin).Methods("POST", "OPTION")
	v1Auth.HandleFunc("/register", controllers.AuthRegister).Methods("POST")

	v1Attendance := apiV1.PathPrefix("/attendance").Subrouter()
	v1Attendance.HandleFunc("/teacher", controllers.AttendanceTeacherIndex).Methods("GET", "OPTION")
	v1Attendance.HandleFunc("/teacher/late_today", controllers.AttendanceTeacherLateToday).Methods("GET", "OPTION")
	v1Attendance.HandleFunc("/teacher/report", controllers.AttendanceTeacherReport).Methods("GET", "OPTION")
	v1Attendance.HandleFunc("/student", controllers.AttendanceStudentIndex).Methods("GET", "OPTION")
	v1Attendance.HandleFunc("/student/late_today", controllers.AttendanceStudentLateToday).Methods("GET", "OPTION")
	v1Attendance.HandleFunc("/student/report", controllers.AttendanceStudentReport).Methods("GET", "OPTION")

	v1Camera := apiV1.PathPrefix("/camera").Subrouter()
	v1Camera.HandleFunc("/", controllers.CameraIndex).Methods("GET")         // View All
	v1Camera.HandleFunc("/", controllers.CameraStore).Methods("POST")        // Store
	v1Camera.HandleFunc("/{id}", controllers.CameraDetail).Methods("GET")    // Get Detail
	v1Camera.HandleFunc("/{id}", controllers.CameraUpdate).Methods("PUT")    // Update
	v1Camera.HandleFunc("/{id}", controllers.CameraDelete).Methods("DELETE") // Delete

	v1ClassRoom := apiV1.PathPrefix("/classroom").Subrouter()
	v1ClassRoom.HandleFunc("/", controllers.ClassroomIndex).Methods("GET")         // View All
	v1ClassRoom.HandleFunc("/", controllers.ClassroomStore).Methods("POST")        // Store
	v1ClassRoom.HandleFunc("/{id}", controllers.ClassroomDetail).Methods("GET")    // Get Detail
	v1ClassRoom.HandleFunc("/{id}", controllers.ClassroomUpdate).Methods("PUT")    // Update
	v1ClassRoom.HandleFunc("/{id}", controllers.ClassroomDelete).Methods("DELETE") // Delete

	v1Schedule := apiV1.PathPrefix("/schedule").Subrouter()
	v1Schedule.HandleFunc("/", controllers.ScheduleIndex).Methods("GET")         // View All
	v1Schedule.HandleFunc("/", controllers.ScheduleStore).Methods("POST")        // Store
	v1Schedule.HandleFunc("/{id}", controllers.ScheduleDetail).Methods("GET")    // Get Detail
	v1Schedule.HandleFunc("/{id}", controllers.ScheduleUpdate).Methods("PUT")    // Update
	v1Schedule.HandleFunc("/{id}", controllers.ScheduleDelete).Methods("DELETE") // Delete

	v1Logs := apiV1.PathPrefix("/logs").Subrouter()
	v1Logs.HandleFunc("/", controllers.LogIndex).Methods("GET")      // View All
	v1Logs.HandleFunc("/", controllers.LogStore).Methods("POST")     // Store
	v1Logs.HandleFunc("/{id}", controllers.LogDetail).Methods("GET") // Get Detail

	v1Subject := apiV1.PathPrefix("/subject").Subrouter()
	v1Subject.HandleFunc("/", controllers.SubjectIndex).Methods("GET")         // View All
	v1Subject.HandleFunc("/", controllers.SubjectStore).Methods("POST")        // Store
	v1Subject.HandleFunc("/{id}", controllers.SubjectDetail).Methods("GET")    // Get Detail
	v1Subject.HandleFunc("/{id}", controllers.SubjectUpdate).Methods("PUT")    // Update
	v1Subject.HandleFunc("/{id}", controllers.SubjectDelete).Methods("DELETE") // Delete

	v1Room := apiV1.PathPrefix("/room").Subrouter()
	v1Room.HandleFunc("/", controllers.RoomIndex).Methods("GET")         // View All
	v1Room.HandleFunc("/", controllers.RoomStore).Methods("POST")        // Store
	v1Room.HandleFunc("/{id}", controllers.RoomDetail).Methods("GET")    // Get Detail
	v1Room.HandleFunc("/{id}", controllers.RoomUpdate).Methods("PUT")    // Update
	v1Room.HandleFunc("/{id}", controllers.RoomDelete).Methods("DELETE") // Delete

	v1RoomAccess := apiV1.PathPrefix("/room_access").Subrouter()
	v1RoomAccess.HandleFunc("/", controllers.RoomAccessIndex).Methods("GET")         // View All
	v1RoomAccess.HandleFunc("/", controllers.RoomAccessStore).Methods("POST")        // Store
	v1RoomAccess.HandleFunc("/{id}", controllers.RoomAccessDetail).Methods("GET")    // Get Detail
	v1RoomAccess.HandleFunc("/{id}", controllers.RoomAccessUpdate).Methods("PUT")    // Update
	v1RoomAccess.HandleFunc("/{id}", controllers.RoomAccessDelete).Methods("DELETE") // Delete

	v1User := apiV1.PathPrefix("/user").Subrouter()
	v1User.HandleFunc("/", controllers.UserIndex).Methods("GET")         // View All
	v1User.HandleFunc("/", controllers.UserStore).Methods("POST")        // Store
	v1User.HandleFunc("/{id}", controllers.UserDetail).Methods("GET")    // Get Detail
	v1User.HandleFunc("/{id}", controllers.UserUpdate).Methods("PUT")    // Update
	v1User.HandleFunc("/{id}", controllers.UserDelete).Methods("DELETE") // Delete

	v2Attendance := apiV2.PathPrefix("/attendance").Subrouter()
	v2Attendance.HandleFunc("/", controllers.AttendanceV2).Methods("GET")        // View
	v2Attendance.HandleFunc("/new", controllers.AttendanceV2New).Methods("POST") // Store

	v2User := apiV2.PathPrefix("/user").Subrouter()
	v2User.HandleFunc("/", controllers.UserV2Index).Methods("GET")                                           // View All
	v2User.HandleFunc("/embeddings", controllers.UserV2Embeddings).Methods("GET")                            // Verify
	v2User.HandleFunc("/embeddings/clear", controllers.UserV2EmbeddingsClear).Methods("GET")                 // Verify
	v2User.HandleFunc("/embeddings/clear/{user_id}", controllers.UserV2EmbeddingsClearOnUser).Methods("GET") // Verify
	v2User.HandleFunc("/{id}", controllers.UserV2Detail).Methods("GET")                                      // Detail
	v2User.HandleFunc("/verify", controllers.UserV2Verify).Methods("POST")                                   // Verify
	v2User.HandleFunc("/register", controllers.UserRegister).Methods("POST")                                 // Store
	v2User.HandleFunc("/recognize", controllers.UserRecognize).Methods("OPTION", "POST")                     // Recognize

	apiV2.HandleFunc("/room_accesss/check", controllers.RoomAccessCheck).Methods("POST")
	apiV2.HandleFunc("/classroom", controllers.ClassroomV2Index).Methods("GET")
	apiV2.HandleFunc("/logs", controllers.LogIndex).Methods("GET")
	apiV2.HandleFunc("/camera", controllers.CameraIndex).Methods("GET")
	apiV2.HandleFunc("/test", controllers.ImportCsv).Methods("GET")

	fmt.Println("App running on port " + APP_PORT)
	log.Fatal(http.ListenAndServe(":"+APP_PORT, router))
}
