package imgutils

import (
	"image"
	_ "image/jpeg" // Initialization for image.Decode()
	_ "image/png"  //
	"os"

	"github.com/nfnt/resize"
)

// OpenImg reads image a from file and returns Image interface
func OpenImg(filename string) (image.Image, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return img, nil
}

// Resize resizes image using resize library
func Resize(w, h uint, img image.Image) image.Image {
	res := resize.Resize(w, h, img, resize.Lanczos3)
	return res
}

// ToGray converts image to grayscale
func ToGray(img image.Image) *image.Gray {
	bounds := img.Bounds()
	gray := image.NewGray(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			gray.Set(x, y, img.At(x, y))
		}
	}
	return gray
}

// Threshold does automatic image segmentation
// using Otsu's Method
//
// http://www.labbookpages.co.uk/software/imgProc/otsuThreshold.html
func Threshold(img *image.Gray) *image.Gray {
	hist := histogramGray(img)
	sum := 0
	for i := 0; i < 256; i++ {
		sum += i * hist[i]
	}
	bVarMax := 0.0
	sumB, sumF := 0, sum
	wB, wF := 0, len(img.Pix)
	threshold := uint8(0)
	for t := 0; t < 256; t++ {
		wB += hist[t]
		wF -= hist[t]
		if wB == 0 {
			continue
		}
		if wF == 0 {
			break
		}
		sumB += t * hist[t]
		sumF = sum - sumB
		mB := float64(sumB) / float64(wB)
		mF := float64(sumF) / float64(wF)
		bVar := float64(wB*wF) * (mB - mF) * (mB - mF)
		if bVar > bVarMax {
			bVarMax = bVar
			threshold = uint8(t)
		}
	}

	// Segment image
	binImg := image.NewGray(img.Bounds())
	for i := 0; i < len(binImg.Pix); i++ {
		if img.Pix[i] > threshold {
			binImg.Pix[i] = 255
		} else {
			binImg.Pix[i] = 0
		}
	}
	return binImg
}

func histogramGray(gray *image.Gray) []int {
	hist := make([]int, 256)
	for i := 0; i < len(gray.Pix); i++ {
		hist[gray.Pix[i]]++
	}
	return hist
}
