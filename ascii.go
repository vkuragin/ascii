package ascii

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math"
	"os"
	"time"

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
	// measure execution time
	start := time.Now()
	defer func() {
		res := time.Since(start)
		log.Printf("Processing took: %v", res)
	}()

	// validate bounds
	bounds := img.src.Bounds()
	if w > bounds.Max.X {
		log.Printf("W exceds image bounds: %d -> %d\n", w, bounds.Max.X)
		w = bounds.Max.X
	}
	if h > bounds.Max.Y {
		log.Printf("H exceds image bounds: %d -> %d\n", h, bounds.Max.Y)
		h = bounds.Max.Y
	}

	// process
	res := ""
	dx, dy := delta(bounds.Max.X, w), delta(bounds.Max.Y, h)
	y, maxY, maxX := float64(0), float64(bounds.Max.Y), float64(bounds.Max.X)
	for y < maxY {
		nextY := math.Min(math.Round(y+dy), maxY)
		c, x := 0, float64(0)

		for x < maxX {
			nextX := math.Min(math.Round(x+dx), maxX)
			if x >= nextX || y >= nextY {
				log.Printf("empty set, skipping: x=%f, x2=%f, y=%f, y2=%f, maxX=%f, maxY=%f\n", x, nextX, y, nextY, maxX, maxY)
				x = nextX
				continue
			}
			point := downsample(&img.src, int(x), int(nextX), int(y), int(nextY))
			res += fmt.Sprintf("%c", colorToAscii(point))
			c++
			x = nextX
		}
		y = nextY
		res += fmt.Sprintf(" | %d\n", c)
	}

	img.asc = res
	return nil
}

func delta(max, steps int) float64 {
	res := float64(max) / float64(steps)
	log.Printf("delta: %d / %d = %f\n", max, steps, res)
	return res
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
		}
	}
	res := color.RGBA{uint8(sumR / count), uint8(sumG / count), uint8(sumB / count), uint8(sumA / count)}

	//log.Printf("downsampling is done for {%d,%d}-{%d,%d}, count=%d\n", x1, y1, x2, y2, count)
	//log.Printf("R=%d, G=%d, B=%d, A=%d -> %+v\n", sumR, sumG, sumB, sumA, res)
	//
	//r, g, b, a := res.RGBA()
	//log.Printf("RGBA() -> %v,%v,%v,%v\n", r, g, b, a)
	//log.Printf("RGBA -> %v,%v,%v,%v\n", res.R, res.G, res.B, res.A)

	return res
}

func colorToAscii(c color.Color) byte {
	charMin := byte(' ')
	charMax := byte('~')
	charRange := int(charMax - charMin)

	r, g, b, _ := c.RGBA()
	tmp := uint32(uint8(r)) + uint32(uint8(g)) + uint32(uint8(b))
	tmp /= uint32(charRange)
	val := tmp%uint32(charRange) + uint32(charMin)
	return byte(val)
}
