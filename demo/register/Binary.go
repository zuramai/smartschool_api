package main

import (
	"fmt"
	"gopkg.in/pipe.v2"
)

func main(){
	p := pipe.Line(
		pipe.Exec("python", "register.py"),
	)
	output, err := pipe.CombinedOutput(p)
	fmt.Println(string(output))

	if err != nil {
		fmt.Printf("%v\n", err)
	}
}