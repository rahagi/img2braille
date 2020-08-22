package imgutils_test

import (
	"image/png"
	"os"
	"testing"

	. "github.com/cytopz/img2braille/imgutils"
)

func TestThresh(t *testing.T) {
	fileName := "test.png"
	img, err := OpenImg(fileName)
	if err != nil {
		t.Errorf("error: %v", err)
	}
	imgResize := Resize(50, 0, img)
	imgGray := ToGray(imgResize)
	file, err := os.Create("test_gray.png")
	if err != nil {
		t.Errorf("error: %v", err)
	}
	png.Encode(file, imgGray)
	defer file.Close()
	binImg := Threshold(imgGray)
	file, err = os.Create("test_bin.png")
	if err != nil {
		t.Errorf("error: %v", err)
	}
	defer file.Close()
	png.Encode(file, binImg)
}
