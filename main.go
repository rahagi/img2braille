package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"log"
	"strconv"

	"github.com/cytopz/img2braille/imgutils"
)

func toBraille(img *image.Gray) string {
	var res string
	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y += 4 {
		for x := bounds.Min.X; x < bounds.Max.X; x += 2 {
			braillePix := []color.Color{
				img.At(x, y),     // Dot 1
				img.At(x, y+1),   // Dot 2
				img.At(x, y+2),   // Dot 3
				img.At(x+1, y),   // Dot 4
				img.At(x+1, y+1), // Dot 5
				img.At(x+1, y+2), // Dot 6
				img.At(x, y+3),   // Dot 7
				img.At(x+1, y+3), // Dot 8
			}
			braille, err := binToBraille(pixToBin(braillePix))
			if err != nil {
				log.Fatal("error converting to braille: ", err)
			}
			res += braille
		}
		res += "\n"
	}
	return res
}

func binToBraille(bin string) (string, error) {
	hexVal, err := strconv.ParseUint(bin, 2, 64)
	if err != nil {
		return "", err
	}
	unicodeVal := "28"
	unicodeVal += fmt.Sprintf("%x", hexVal)
	if len(unicodeVal) == 3 {
		unicodeVal += "0"
	}
	res, err := strconv.Unquote(`"\u` + unicodeVal + `"`)
	if err != nil {
		return "", err
	}
	return res, nil
}

func pixToBin(pix []color.Color) string {
	var res string
	for i := len(pix) - 1; i >= 0; i-- {
		y, _, _, _ := pix[i].RGBA()
		if uint8(y) == 255 {
			res += "1"
		} else {
			res += "0"
		}
	}
	return res
}

func main() {
	fileName := flag.String("i", "", "input file")
	w := flag.Uint("w", 50, "width")
	h := flag.Uint("h", 0, "height")
	flag.Parse()

	img, err := imgutils.OpenImg(*fileName)
	if err != nil {
		log.Fatal("error opening file: ", err)
	}
	img = imgutils.Resize(*w, *h, img)
	binImg := imgutils.Threshold(imgutils.ToGray(img))
	result := toBraille(binImg)
	fmt.Printf("%s\n", result)
}
