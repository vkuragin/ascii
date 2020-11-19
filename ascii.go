package ascii

import (
	"image"
	"image/color"
	"log"
	"math"
	"os"
	"time"

	_ "image/jpeg"
	_ "image/png"
)

// Img contains original binary image and its ascii representation
type Img struct {
	src image.Image
	asc []byte
}

// Load loads image from file and returns new Img object
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

// Process converts image from source to ascii representation (w - ascii width, h - ascii height)
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

	// compute steps
	dx, dy := delta(bounds.Max.X, w), delta(bounds.Max.Y, h)
	xSteps, ySteps := steps(bounds.Max.X, dx), steps(bounds.Max.Y, dy)
	log.Printf("steps: w=%d, xSteps=%d, h=%d, ySteps=%d\n", w, xSteps, h, ySteps)

	// process image
	img.asc = make([]byte, ySteps*xSteps+ySteps)
	lastX, lastY, index := 0, 0, 0
	for i := 0; i < ySteps; i++ {
		y := lastY + int(dy)
		lastX = 0
		for j := 0; j < xSteps; j++ {
			x := lastX + int(dx)
			processArea(img, index, lastX, x, lastY, y)
			index++
			lastX = x
		}
		lastY = y
		img.asc[index] = 10
		index++
	}
	return nil
}

func steps(max int, delta float64) int {
	res := float64(max) / delta
	log.Printf("steps: %d/%f = %f\n", max, delta, res)
	return int(math.Ceil(res))
}

func processArea(img *Img, index, x1, x2, y1, y2 int) {
	point := downsample(&img.src, x1, x2, y1, y2)
	ascii := colorToAscii(point)
	img.asc[index] = ascii
	//log.Printf("{%d, %d} - {%d, %d} = %q\n", x1, y1, x2, y2, ascii)
}

func delta(max, steps int) float64 {
	res := float64(max) / float64(steps)
	return math.Floor(res)
}

// WriteToFile writes ascii representation to the output file
func (img *Img) WriteToFile(filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}

	n, err := file.Write(img.asc)
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

// Result return ascii representation
func (img *Img) Result() string {
	return string(img.asc)
}

// downsample computes new colored point from the image area (x1,y1)-(x2,y2): average values for each layer R,G,B,A
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

	return res
}

// colorToAscii converts colored point to an ascii character
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
