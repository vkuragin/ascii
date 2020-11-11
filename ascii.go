package ascii

import (
	"fmt"
	"image"
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

	span := 126 - 32
	res := ""
	for i := 0; i < bounds.Max.Y; i += bounds.Max.Y / h {
		c := 0
		for j := 0; j < bounds.Max.X; j += bounds.Max.X / w {
			color := img.src.At(j, i)
			r, g, b, _ := color.RGBA()
			val := int(r+g+b)%span + 32
			res += fmt.Sprintf("%c", byte(val))
			c++
		}
		res += fmt.Sprintf("|%d\n", c)

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
