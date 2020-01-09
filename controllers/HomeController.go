package controllers

import (
	"fmt"
	"net/http"

	socketio_client "github.com/zhouhui8915/go-socket.io-client"
)

// Home Controller
func HomeController(w http.ResponseWriter, r *http.Request) {
	fmt.Println("home")
}

var socketClient *socketio_client.Client

func ConnectSocket() {

}
