package group

import (
	"bytes"
	// "fmt"
	"image/color"
	"image/png"

	// "os"
	// "time"

	"github.com/jizizr/wordclouds"
)

var DefaultColors = []color.RGBA{
	{0x1b, 0x1b, 0x1b, 0xff},
	{0x48, 0x48, 0x4B, 0xff},
	{0x59, 0x3a, 0xee, 0xff},
	{0x65, 0xCD, 0xFA, 0xff},
	{0x70, 0xD6, 0xBF, 0xff},
}

var boxes = wordclouds.Mask(
	"source/mask.png",
	600,
	600,
	color.RGBA{
		R: 0,
		G: 0,
		B: 0,
		A: 0,
	})

var oarr []wordclouds.Option
var colors = make([]color.Color, 0)

func init() {
	for _, c := range DefaultColors {
		colors = append(colors, c)
	}
	oarr = []wordclouds.Option{
		wordclouds.FontFile("source/font.ttf"),
		wordclouds.FontMaxSize(100),
		wordclouds.FontMinSize(10),
		wordclouds.Colors(colors),
		wordclouds.MaskBoxes(boxes),
		wordclouds.Height(600),
		wordclouds.Width(600),
		wordclouds.RandomPlacement(false),
		wordclouds.WordSizeFunction("linear"),
	}
}

func Rank(inputWords map[string]int) []byte {

	// Load words
	// inputWords := map[string]int{"消息": 42, "是啊": 30, "中文": 15, "也是": 10, "而我": 5, "撒旦": 11, "落后": 11}

	// Load config

	// start := time.Now()
	w := wordclouds.NewWordcloud(inputWords,
		oarr...)

	// outputFile, _ := os.Create("test.png")
	buf := new(bytes.Buffer)
	// Encode takes a writer interface and an image interface
	// We pass it the File and the RGBA
	png.Encode(buf, w.Draw())

	// Don't forget to close files
	// outputFile.Close()
	// fmt.Printf("Done in %v\n", time.Since(start))
	return buf.Bytes()
}
