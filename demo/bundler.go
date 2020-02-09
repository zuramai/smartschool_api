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
		pipe.Exec("python", "CyberNet.py"),
	)
	output, err := pipe.CombinedOutput(p)
	fmt.Println(string(output))

	if err != nil {
		fmt.Printf("%v\n", err)
	}
}
