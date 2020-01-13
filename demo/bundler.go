package main

import (
	// "time"
	// "bytes"
	// "encoding/base64"
	"fmt"
	// "gocv.io/x/gocv"
	"gopkg.in/pipe.v2"
	// "image/jpeg"
	// "log"
	// "os"
	// "os/exec"
)

func main() {
	p := pipe.Line(
		pipe.Exec("conda", "activate"),
		pipe.Exec("python", "demo.py"),
	)
	output, err := pipe.CombinedOutput(p)
	fmt.Println(string(output))
	// decoded, err := base64.StdEncoding.DecodeString(string(output))
	// image, err := jpeg.Decode(bytes.NewReader(decoded))
	// mat, err := gocv.IMDecode(decoded, gocv.IMReadUnchanged)
	// window := gocv.NewWindow("window")
	// window.IMShow(mat)
	// window.WaitKey(0)
	// time.Sleep(10000)

	if err != nil {
		fmt.Printf("%v\n", err)
	}
	// fmt.Println(image.Bounds())

	// cmd := exec.Command("D:\\Users\\Athanatius.C\\Anaconda3\\python.exe", "demo.py")
	// // cmd.Stdin = os.Stdin
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr
	// err := cmd.Run()
	// if err != nil {
	// 	log.Println(err)
	// }
}
