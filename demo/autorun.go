package main

import(
	"github.com/ProtonMail/go-autostart"
	"log"
)

func main(){
	app := &autostart.App{
		Name: "CyberNet",
		DisplayName: "CyberNet",
		Exec: []string{"bash", "-c", "echo autostart >> ~/autostart.txt"},
	}

	if app.IsEnabled() {
		log.Println("App is already enabled")
	} else {
		log.Println("Enabling app...")

		if err := app.Enable(); err != nil {
			log.Fatal(err)
		}
	}
}