package asciiconvertor

import (
	"image"
	"strings"

	"github.com/muesli/termenv"
	"golang.org/x/image/draw"
)

const dict string = ".,:;i1tfLCG08@"

func pixel2Ascii(r, g, b, a uint32) byte {
	value := (r + g + b) * a / 255
	precision := 255 * 3 / (len(dict) - 1)
	rawChar := dict[int(value)/precision]
	return rawChar
}

func ConverImage2Ascii(img image.Image, targetW, targetH int, ratio float64) string {
	oldBounds := img.Bounds()
	oldWidth := oldBounds.Dx()
	oldHeight := oldBounds.Dy()
	newWidth := targetW
	newHeight := int(float64(newWidth*oldHeight) / float64(oldWidth) * ratio)
	if newHeight > targetH {
		newHeight = targetH
		newWidth = int(float64(newHeight*oldWidth) / float64(oldHeight) / ratio)
	}
	newImg := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))
	draw.CatmullRom.Scale(newImg, newImg.Bounds(), img, oldBounds, draw.Over, nil)
	p := termenv.ColorProfile()

	var asciiArt strings.Builder

	for y := 0; y < newHeight; y++ {
		for x := 0; x < newWidth; x++ {
			r, g, b, a := newImg.At(x, y).RGBA()
			asciiChar := string(pixel2Ascii(r>>8, g>>8, b>>8, a>>8))
			color := p.FromColor(newImg.At(x, y))
			coloredChar := termenv.String(asciiChar).Foreground(color).String()
			asciiArt.WriteString(coloredChar)
		}
		asciiArt.WriteString("\n")
	}
	return asciiArt.String()
}
