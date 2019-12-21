package controllers

import (
	"net/http"

	"github.com/zuramai/smartschool_api/models"
	"golang.org/x/crypto/bcrypt"
)

func UserIndex(w http.ResponseWriter, r *http.Request) {

}

func UserRegister(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	name := r.FormValue("name")
	phone := r.FormValue("phone")
	role := r.FormValue("role")
	user := models.User{
		Username: username,
		Password: string(hashedPassword),
		Name:     name,
		Phone:    phone,
		Role:     role,
		Status:   true,
	}

	models.GetDB("main").NewRecord(user)
	models.GetDB("main").Create(&user)
	if role == "Superadmin" {

	} else if role == "Moderator" {

	} else if role == "Teacher" {

	} else if role == "Student" {

	}
	respondJSON(w, 200, "Successfully register user!", user)

}

func UserStore(w http.ResponseWriter, r *http.Request) {

}

func UserDetail(w http.ResponseWriter, r *http.Request) {

}

func UserUpdate(w http.ResponseWriter, r *http.Request) {

}

func UserDelete(w http.ResponseWriter, r *http.Request) {

}

func UserV2Index(w http.ResponseWriter, r *http.Request) {

}

func UserV2Detail(w http.ResponseWriter, r *http.Request) {

}

func UserV2Verify(w http.ResponseWriter, r *http.Request) {

}
