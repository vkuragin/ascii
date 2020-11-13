package ascii

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"os"

	_ "image/jpeg"
	_ "image/png"
)

type Img struct {
	src image.Image
	asc string
}

func Load(filePath string) (*Img, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	img, str, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	err = file.Close()
	if err != nil {
		return nil, err
	}

	bounds := img.Bounds()
	log.Printf("image type: %v, bounds=%v\n", str, bounds)

	return &Img{src: img}, nil
}

func (img *Img) Process(w, h int) error {
	bounds := img.src.Bounds()

	res := ""
	for i := 0; i < bounds.Max.Y; i += bounds.Max.Y / h {
		c := 0
		for j := 0; j < bounds.Max.X; j += bounds.Max.X / w {
			point := downsample(&img.src, j, j+bounds.Max.X/w, i, i+bounds.Max.Y/h)
			res += fmt.Sprintf("%c", colorToAscii(point))
			c++
		}
		res += fmt.Sprintf(" | %d\n", c)

	}
	img.asc = res
	return nil
}

func (img *Img) WriteToFile(filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}

	n, err := file.Write([]byte(img.asc))
	if err != nil {
		return err
	}

	log.Printf("Succesfully written %d bytes to file %s\n", n, filePath)

	err = file.Close()
	if err != nil {
		return err
	}

	return nil
}

func (img *Img) Result() string {
	return img.asc
}

func downsample(img *image.Image, x1, x2, y1, y2 int) color.Color {
	sumR, sumG, sumB, sumA := uint32(0), uint32(0), uint32(0), uint32(0)
	count := uint32(0)

	for row := y1; row < y2; row++ {
		for col := x1; col < x2; col++ {
			count++
			r, g, b, a := (*img).At(col, row).RGBA()
			sumR, sumG, sumB, sumA = sumR+uint32(uint8(r)), sumG+uint32(uint8(g)), sumB+uint32(uint8(b)), sumA+uint32(uint8(a))
			//sumR, sumG, sumB, sumA = sumR+r, sumG+g, sumB+b, sumA+a
		}
	}
	res := color.RGBA{uint8(sumR / count), uint8(sumG / count), uint8(sumB / count), uint8(sumA / count)}

	log.Printf("downsampling is done for {%d,%d}-{%d,%d}, count=%d\n", x1, y1, x2, y2, count)
	log.Printf("R=%d, G=%d, B=%d, A=%d -> %+v\n", sumR, sumG, sumB, sumA, res)

	r, g, b, a := res.RGBA()
	log.Printf("RGBA() -> %v,%v,%v,%v\n", r, g, b, a)
	log.Printf("RGBA -> %v,%v,%v,%v\n", res.R, res.G, res.B, res.A)

	return res
}

func colorToAscii(c color.Color) byte {
	charMin := byte(' ')
	charMax := byte('~')
	charRange := int(charMax - charMin)

	r, g, b, _ := c.RGBA()
	tmp := (uint32(uint8(r)) + uint32(uint8(g)) + uint32(uint8(b)))
	s := fmt.Sprintf("tmp=%d ", tmp)
	tmp /= uint32(charRange)
	s += fmt.Sprintf("tmp/charRange=%d ", tmp)
	val := tmp%uint32(charRange) + uint32(charMin)
	s += fmt.Sprintf("val=%d\n", val)

	log.Printf("%s\n------\n", s)
	return byte(val)
}
