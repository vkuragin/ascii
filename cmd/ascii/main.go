package main

import (
	"flag"
	"fmt"
	"github.com/vkuragin/ascii"
	"log"
)

func main() {

	inFile := flag.String("in", "example.png", "path to image file (jpeg/png)")
	outFile := flag.String("out", "ascii.txt", "output txt file")
	w := flag.Int("w", 80, "output pic width in chars")
	h := flag.Int("h", 50, "output pic height in chars")
	p := flag.Bool("p", false, "print the result to console or not")
	c := flag.Bool("c", false, "compute concurrently or not")
	flag.Parse()

	img, err := ascii.Load(*inFile)
	if err != nil {
		log.Fatalf("Failed to load image from file: %s.\nError: %v\n", *inFile, err)
		panic(err)
	}

	err = img.Process(*w, *h, *c)
	if err != nil {
		log.Fatalf("Failed to process image from file: %s\nError: %v\n", *inFile, err)
		panic(err)
	}

	err = img.WriteToFile(*outFile)
	if err != nil {
		log.Fatalf("Failed to save image to file: %s\nError: %v\n", *outFile, err)
		panic(err)
	}

	if *p {
		fmt.Printf("%s\n", img.Result())
	}
}
