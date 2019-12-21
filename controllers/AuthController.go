package controllers

import (
	"fmt"
	"net/http"
	"os"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"github.com/zuramai/smartschool_api/models"
	"golang.org/x/crypto/bcrypt"
)

func AuthIndex() {
	fmt.Println("auth")
}

type Response struct {
	Status bool          `json:"status"`
	Data   []interface{} `json:"data"`
}

func AuthLogin(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	user := models.User{}

	err := models.GetDB("main").Where("username = ?", username).First(&user).Error

	var returnMsg interface{}

	if err != nil {
		if err := models.GetDB("main").Where("username = ?", username).First(&user).Error; err != nil && err != gorm.ErrRecordNotFound {
			respondJSON(w, http.StatusInternalServerError, "Username not found", returnMsg)
			return
		}

		respondJSON(w, http.StatusInternalServerError, "Internal Server Error", returnMsg)
		return
	}

	hashedPassword := user.Password

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))

	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		respondJSON(w, http.StatusUnauthorized, "Invalid Credentials", user)
		return
	}

	tk := &models.Token{UserID: user.ID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))

	user.Token = tokenString

	respondJSON(w, http.StatusOK, "Successfully Logged In!", user)
	return
	// w.Header().Set("Content-Type", "application/json")

}

func AuthRegister(w http.ResponseWriter, r *http.Request) {

}
