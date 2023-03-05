package factories

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"math/rand"
)

// FakeImageBuffer will make a random fake image buffer
func FakeImageBuffer() *bytes.Buffer {
	// Create an 100 x 50 image
	img := image.NewRGBA(image.Rect(0, 0, 100, 50))

	// Draw a random colored image
	img.Set(2, 3, color.RGBA{
		uint8(rand.Intn(250)),
		uint8(rand.Intn(250)),
		uint8(rand.Intn(250)),
		uint8(rand.Intn(250)),
	})

	buffer := new(bytes.Buffer)
	png.Encode(buffer, img)

	return buffer
}
