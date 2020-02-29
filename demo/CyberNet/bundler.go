package main

import (
	// "time"
	// "bytes"
	// "encoding/base64"

	"fmt"
	"log"
	"os/exec"
	"strconv"

	// "gocv.io/x/gocv"
	"encoding/json"
	"os"
	// "image/jpeg"
	// "log"
	// "os"
)

func main() {
	conf := ReadConfig()
	var ProcessList []exec.Cmd
	for i,cam := range(conf.camera){
		fmt.Println("\nCamera "+strconv.Itoa(i)+" Configuration :")
		fmt.Println("Camera ID : "+cam.CameraID)
		fmt.Println("RTSP Link : "+cam.CameraRTSPLink)
		fmt.Println("Stream Port : "+cam.CameraStreamPort)
		fmt.Println("End Point : "+cam.ServerEndPoint+"\n")
		cmd := Parse(cam)
		ProcessList = append(ProcessList,cmd)
	}
	for _,cmd := range(ProcessList){
		fmt.Println("Starting")
		err := cmd.Start()
		if err != nil {
			log.Println(err)
		}
	}
	for _,cmd := range(ProcessList){
		cmd.Wait()
	}
}
type Camera struct {
    CameraID string
	CameraRTSPLink string
	CameraStreamPort string
	ServerEndPoint string
}

type Configuration struct{
	camera []Camera
}

func Parse(cam Camera)(command exec.Cmd){
	cmd := exec.Command("python","CyberNet.py",
		"--camera_id",cam.CameraID,
		"--rtsp_link",cam.CameraRTSPLink,
		"--stream_port",cam.CameraStreamPort,
		"--end_point",cam.ServerEndPoint)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return *cmd
	// fmt.Println(output)
}

func ErrorHandle(err error){
	if err != nil{
		fmt.Println(err)
	}
} 

func ReadConfig()(config Configuration){
	var cameras []Camera
	var configuration Configuration
	file, err := os.Open("configuration.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&cameras)
	ErrorHandle(err)
	configuration.camera = cameras
	return configuration
}