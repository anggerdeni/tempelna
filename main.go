package main

import (
	_ "image/jpeg"
	"log"
	"os"

	"github.com/anggerdeni/tempelna/image_manipulator"
	"github.com/anggerdeni/tempelna/unsplash"
)

func main() {
	reader, err := os.Open("assets/Amiri-Regular.ttf")
	if err != nil {
		log.Fatal(err.Error())
	}

	imageManipulator, err := image_manipulator.New(72, reader)
	if err != nil {
		log.Fatal(err.Error())
	}

	unsplash := unsplash.New()

	srcImg, err := unsplash.GetImage()
	if err != nil {
		log.Fatal(err.Error())
	}

	text := `Lorem ipsum dolor sit amet, consectetur adipiscing elit.
Nam ultrices leo at mollis blandit.
Integer finibus orci ac eros scelerisque vulputate.
Pellentesque eu neque molestie, lobortis felis id, convallis sem.
Aliquam eu est vulputate, interdum arcu quis, vestibulum ligula.
Fusce sodales metus ac quam consequat eleifend.`

	err = imageManipulator.AddTextToCenterOfImage(srcImg, text, "assets/result.jpg")
	if err != nil {
		log.Fatal(err.Error())
	}
}
