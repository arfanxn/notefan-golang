package factories

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math/rand"
)

// FakeImageBuffer will make a random fake image buffer
func FakeImageBuffer() *bytes.Buffer {
	// Create an 100 x 50 image
	img := image.NewRGBA(image.Rect(0, 0, 400, 400))

	// Image color
	rgba := color.RGBA{
		uint8(rand.Intn(250)),
		uint8(rand.Intn(250)),
		uint8(rand.Intn(250)),
		255,
	}

	// Draw a random colored image
	draw.Draw(img, img.Bounds(), &image.Uniform{rgba}, image.Point{}, draw.Src)

	buffer := new(bytes.Buffer)
	png.Encode(buffer, img)

	return buffer
}
