package image_manipulator

import (
	"bufio"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"os"
	"strings"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
)

type ImageManipulator interface {
	AddTextToCenterOfImage(r io.Reader, text string, outfile string) error
}

type imageManipulator struct {
	freeTypeContext *freetype.Context
	fontSize        float64
	fontType        *truetype.Font
}

func New(fontSize float64, font io.Reader) (ImageManipulator, error) {
	fontBytes, err := io.ReadAll(font)
	if err != nil {
		return nil, err
	}

	fontType, err := freetype.ParseFont(fontBytes)
	if err != nil {
		return nil, err
	}

	c := freetype.NewContext()
	c.SetFontSize(fontSize)
	c.SetFont(fontType)
	return &imageManipulator{
		freeTypeContext: c,
		fontSize:        fontSize,
		fontType:        fontType,
	}, nil
}

func (im *imageManipulator) AddTextToCenterOfImage(r io.Reader, text string, outfile string) error {
	srcImg, _, err := image.Decode(r)
	if err != nil {
		return err
	}

	bounds := srcImg.Bounds()
	textSeparated := strings.Split(text, "\n")
	textBoundingBoxes := []image.Rectangle{}

	for _, txt := range textSeparated {
		boundingBox, err := im.findBoundingBox(bounds, im.fontSize, txt, im.fontType)
		if err != nil {
			log.Fatal(err.Error())
		}

		textBoundingBoxes = append(textBoundingBoxes, boundingBox)
	}

	resultImg := image.NewRGBA(bounds)
	draw.Draw(resultImg, bounds, srcImg, image.Point{}, draw.Src)

	complimentColor := im.findComplimentColor(srcImg, textBoundingBoxes)

	im.freeTypeContext.SetClip(bounds)
	im.freeTypeContext.SetDst(resultImg)
	im.freeTypeContext.SetSrc(image.NewUniform(complimentColor))

	for i, txt := range textSeparated {
		x := im.findStartingXPoint(bounds, textBoundingBoxes[i])
		y := im.findStartingYPoint(bounds, textBoundingBoxes[i], len(textBoundingBoxes), i)

		pt := freetype.Pt(x, y)
		_, err = im.freeTypeContext.DrawString(txt, pt)
		if err != nil {
			return err
		}
	}

	outFile, err := os.Create(outfile)
	if err != nil {
		log.Fatal(err)
	}
	defer outFile.Close()

	b := bufio.NewWriter(outFile)
	err = jpeg.Encode(b, resultImg, &jpeg.Options{Quality: 90})
	if err != nil {
		log.Fatal(err)
	}
	err = b.Flush()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Image saved.")
	return nil
}

func (im *imageManipulator) findBoundingBox(originalBound image.Rectangle, fontSize float64, s string, font *truetype.Font) (image.Rectangle, error) {
	tmpImg := image.NewRGBA(originalBound)
	im.freeTypeContext.SetClip(originalBound)
	im.freeTypeContext.SetDst(tmpImg)
	im.freeTypeContext.SetSrc(image.White)

	pt := freetype.Pt(0, im.freeTypeContext.PointToFixed(fontSize).Round())
	_, err := im.freeTypeContext.DrawString(s, pt)
	if err != nil {
		return image.Rectangle{}, err
	}

	var box image.Rectangle
	for y := 0; y < tmpImg.Bounds().Dy(); y++ {
		for x := 0; x < tmpImg.Bounds().Dx(); x++ {
			if r, _, _, _ := tmpImg.At(x, y).RGBA(); r != 0 {
				box = box.Union(image.Rect(x, y, x+1, y+1))
			}
		}
	}

	return box, nil
}

func (im *imageManipulator) findStartingXPoint(originalBound image.Rectangle, fontBound image.Rectangle) int {
	middlePoint := image.Point{
		X: originalBound.Dx() / 2,
		Y: originalBound.Dy() / 2,
	}

	pxFontWidth := fontBound.Max.X - fontBound.Min.X

	return middlePoint.X - (pxFontWidth / 2)
}

func (im *imageManipulator) findStartingYPoint(originalBound image.Rectangle, fontBound image.Rectangle, countLines int, currentLineNumber int) int {
	spacing := int(0.5 * im.fontSize)

	middlePoint := image.Point{
		X: originalBound.Dx() / 2,
		Y: originalBound.Dy() / 2,
	}

	pxFontHeight := fontBound.Max.Y - fontBound.Min.Y
	totalPxFontHeight := (countLines * pxFontHeight) + (spacing * (countLines - 1))
	start := middlePoint.Y - (totalPxFontHeight / 2)

	return start + (currentLineNumber * pxFontHeight) + (currentLineNumber * spacing)

}

func (im *imageManipulator) findComplimentColor(srcImg image.Image, boundingBoxes []image.Rectangle) color.Color {
	minX := boundingBoxes[0].Min.X
	minY := boundingBoxes[0].Min.Y
	maxX := boundingBoxes[0].Max.X
	maxY := boundingBoxes[0].Max.Y

	for _, b := range boundingBoxes {
		if b.Min.X < minX {
			minX = b.Min.X
		}
		if b.Min.Y < minY {
			minY = b.Min.Y
		}

		if b.Max.X > maxX {
			maxX = b.Max.X
		}
		if b.Max.Y > maxY {
			maxY = b.Max.Y
		}
	}

	impactedArea := image.Rectangle{
		Min: image.Point{minX, minY},
		Max: image.Point{maxX, maxY},
	}
	// initialize variables to hold total color and total pixels
	var totalColorR, totalColorG, totalColorB, totalPixels uint32

	// iterate over each pixel in the impacted area
	for y := impactedArea.Min.Y; y < impactedArea.Max.Y; y++ {
		for x := impactedArea.Min.X; x < impactedArea.Max.X; x++ {
			// get the color of the pixel
			r, g, b, _ := srcImg.At(x, y).RGBA()

			// add the color to the total
			totalColorR += r
			totalColorG += g
			totalColorB += b

			// increment the total pixel count
			totalPixels++
		}
	}

	// calculate the average color
	averageColorR := totalColorR / totalPixels
	averageColorG := totalColorG / totalPixels
	averageColorB := totalColorB / totalPixels

	// calculate the complimentary color
	complimentaryColorR := uint8(255 - (averageColorR >> 8))
	complimentaryColorG := uint8(255 - (averageColorG >> 8))
	complimentaryColorB := uint8(255 - (averageColorB >> 8))

	// return the complimentary color
	return color.RGBA{complimentaryColorR, complimentaryColorG, complimentaryColorB, 255}
}
