package main

import (
	"path/filepath"
	"strconv"
	// "time"
	// "bytes"
	// "encoding/base64"
	"fmt"
	"gocv.io/x/gocv"
	"gopkg.in/pipe.v2"
	// "image/jpeg"
	"log"
	// "os"
	// "os/exec"
)

func main() {
	image := gocv.IMRead(filepath.Join("C:\\Users\\Athanatius.C\\Pictures\\IMG_20191030_183641.png"), gocv.IMReadColor)
	buff, err := gocv.IMEncode(gocv.PNGFileExt, image)
	print(buff)
	if err != nil {
		log.Println(err)
	}
	argument := strconv.Quote("import canvas;canvas.main('Test')")

	p := pipe.Line(
		pipe.Exec("python", "-c", argument),
	)
	output, err := pipe.CombinedOutput(p)
	fmt.Println(string(output))

	if err != nil {
		fmt.Printf("%v\n", err)
	}
}
