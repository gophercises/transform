package main

import (
	"io"
	"os"

	"github.com/gophercises/transform/primitive"
)

func main() {
	inFile, err := os.Open("tmp/jump.png")
	if err != nil {
		panic(err)
	}
	defer inFile.Close()
	out, err := primitive.Transform(inFile, 33)
	if err != nil {
		panic(err)
	}
	os.Remove("out.png")
	outFile, err := os.Create("out.png")
	if err != nil {
		panic(err)
	}
	io.Copy(outFile, out)
}
