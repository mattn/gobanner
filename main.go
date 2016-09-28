package main

import (
	"flag"
	"fmt"
	"image"
	"image/draw"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/golang/freetype"
)

var (
	font  = flag.String("f", os.Getenv("GOBANNER_FONTFILE"), "path to TTF file")
	width = flag.Int("w", 79, "screen width")
	size  = flag.Int("s", 15, "font size")
)

func main() {
	flag.Parse()

	var args []string

	if *font == "" {
		flag.Usage()
		os.Exit(1)
	}

	if flag.NArg() == 0 {
		b, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			log.Fatal(err)
		}
		args = strings.Split(string(b), "\n")
	} else {
		args = flag.Args()
	}

	b, err := ioutil.ReadFile(*font)
	if err != nil {
		log.Fatal(err)
	}
	f, err := freetype.ParseFont(b)
	if err != nil {
		log.Fatal(err)
	}
	fc := freetype.NewContext()
	fc.SetDPI(72)
	fc.SetFont(f)
	fc.SetFontSize(float64(*size))

	fw := int(fc.PointToFixed(float64(*size))/72) + 1

	rgba := image.NewRGBA(image.Rect(0, 0, *width, fw*len(args)+1))
	draw.Draw(rgba, rgba.Bounds(), image.Black, image.ZP, draw.Src)

	fc.SetClip(rgba.Bounds())
	fc.SetDst(rgba)
	fc.SetSrc(image.White)

	for i := 0; i < len(args); i++ {
		pt := freetype.Pt(0, fw*(i+1))
		fc.DrawString(args[i], pt)
	}

	for y := 0; y < rgba.Bounds().Dy(); y++ {
		for x := 0; x < rgba.Bounds().Dx(); x++ {
			r, g, b, _ := rgba.At(x, y).RGBA()
			if r == 0 && g == 0 && b == 0 {
				fmt.Print(" ")
			} else {
				fmt.Print("#")
			}
		}
		fmt.Println()
	}
}
