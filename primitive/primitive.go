package primitive

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

// Mode defines the shapes used when transforming images.
type Mode int

// Modes supported by the primitive package.
const (
	ModeCombo Mode = iota
	ModeTriangle
	ModeRect
	ModeEllipse
	ModeCircle
	ModeRotatedRect
	ModeBeziers
	ModeRotatedEllipse
	ModePolygon
)

// WithMode is an option for the Transform function that will define the
// mode you want to use. By default, ModeTriangle will be used.
func WithMode(mode Mode) func() []string {
	return func() []string {
		return []string{"-m", fmt.Sprintf("%d", mode)}
	}
}

// Transform will take the provided image and apply a primitive
// transformation to it, then return a reader to the resulting image.
func Transform(image io.Reader, numShapes int, opts ...func() []string) (io.Reader, error) {
	in, err := tempfile("in_", "png")
	if err != nil {
		return nil, errors.New("primitive: failed to create temporary input file")
	}
	defer os.Remove(in.Name())
	out, err := tempfile("in_", "png")
	if err != nil {
		return nil, errors.New("primitive: failed to create temporary output file")
	}
	defer os.Remove(out.Name())

	// Read image into in file
	_, err = io.Copy(in, image)
	if err != nil {
		return nil, errors.New("primitive: failed to copy image into temp input file")
	}

	// Run primitive w/ -i in.Name() -o out.Name()
	stdCombo, err := primitive(in.Name(), out.Name(), numShapes, ModeCombo)
	if err != nil {
		return nil, fmt.Errorf("primitive: failed to run the primitive command. stdcombo=%s", stdCombo)
	}
	fmt.Println(stdCombo)

	// read out into a reader, return reader, delete out
	b := bytes.NewBuffer(nil)
	_, err = io.Copy(b, out)
	if err != nil {
		return nil, errors.New("primitive: Failed to copy output file into byte buffer")
	}
	return b, nil
}

func primitive(inputFile, outputFile string, numShapes int, mode Mode) (string, error) {
	argStr := fmt.Sprintf("-i %s -o %s -n %d -m %d", inputFile, outputFile, numShapes, mode)
	cmd := exec.Command("primitive", strings.Fields(argStr)...)
	b, err := cmd.CombinedOutput()
	return string(b), err
}

func tempfile(prefix, ext string) (*os.File, error) {
	in, err := ioutil.TempFile("", "in_")
	if err != nil {
		return nil, errors.New("primitive: failed to create temporary file")
	}
	defer os.Remove(in.Name())
	return os.Create(fmt.Sprintf("%s.%s", in.Name(), ext))
}
