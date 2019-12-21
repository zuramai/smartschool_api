package controllers

import (
	"fmt"
	"net/http"
)


// Home Controller
func HomeController(w http.ResponseWriter, r *http.Request) {
	fmt.Println("home")
}
