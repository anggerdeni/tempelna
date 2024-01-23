package main

import (
	"encoding/json"
	"fmt"
	_ "image/jpeg"
	"io"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/abdullahdiaa/garabic"
	"github.com/anggerdeni/tempelna/image_manipulator"
	"github.com/anggerdeni/tempelna/unsplash"
)

func main() {
	reader, err := os.Open("assets/Amiri-Regular.ttf")
	if err != nil {
		log.Fatal(err.Error())
	}

	imageManipulator, err := image_manipulator.New(128, reader)
	if err != nil {
		log.Fatal(err.Error())
	}

	unsplash := unsplash.New(os.Getenv("UNSPLASH_ACCESS_KEY"))

	srcImg, err := unsplash.GetImage()
	if err != nil {
		log.Fatal(err.Error())
	}

	texts := getTexts()

	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)
	text := texts[r.Intn(len(texts))]
	log.Println(text)

	err = imageManipulator.AddTextToCenterOfImage(srcImg, text, "assets/result.jpg")
	if err != nil {
		log.Fatal(err.Error())
	}
}

func getTexts() []string {
	type Src struct {
		Arabic  string `json:"arabic"`
		English string `json:"english"`
	}

	jsonFile, err := os.Open("assets/texts.json")
	if err != nil {
		log.Fatalf("Failed to open JSON file: %s", err)
	}
	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)

	var srcs []Src
	json.Unmarshal(byteValue, &srcs)

	texts := make([]string, 0)

	for _, src := range srcs {
		texts = append(texts, fmt.Sprintf("%s\n%s", garabic.Shape(src.Arabic), src.English))
	}
	return texts
}
