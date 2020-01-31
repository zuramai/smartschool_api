package main

import (
	"gocv.io/x/gocv"
	"image"
	"log"
)

func main() {
	net := gocv.ReadNet("..\\demo\\models\\openface_nn4.small2.v1.t7", "")
	defer net.Close()
	img := gocv.IMRead("C:\\Users\\Athanatius.C\\Pictures\\Dataset\\Lexi\\IMG_20200130_083534.jpg", gocv.IMReadColor)
	blob := gocv.BlobFromImage(img, 1.0/255, image.Pt(96, 96), gocv.NewScalar(0, 0, 0, 0), true, false)
	net.SetInput(blob, "input")
	log.Println(net.Forward("softmax2"))
}
