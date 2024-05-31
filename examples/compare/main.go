package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
	"path"
	"time"

	"github.com/fogleman/gg"
	"github.com/gary23b/easygif"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font/gofont/goregular"
)

func main() {
	chooseNearest()
	chooseNearestPlan9()
	useDithering()
	useDitheringPlan9()
	ditheringArtifactExample()
	useMostCommonColors()
	easyGifExample()
	EasyGifExample2()

	ExampleWithOneColorCountEach()
	CheckThatOutliersAreKept()
}

// https://www.w3schools.com/colors/colors_names.asp
var (
	// Black to white
	Black     color.RGBA = color.RGBA{0x00, 0x00, 0x00, 0xFF} // #000000
	DarkGray  color.RGBA = color.RGBA{0x26, 0x26, 0x26, 0xFF} // #262626
	Gray      color.RGBA = color.RGBA{0x80, 0x80, 0x80, 0xFF} // #808080
	LightGray color.RGBA = color.RGBA{0xD3, 0xD3, 0xD3, 0xFF} // #D3D3D3
	White     color.RGBA = color.RGBA{0xFF, 0xFF, 0xFF, 0xFF} // #FFFFFF

	// Primary Colors
	Red  color.RGBA = color.RGBA{0xFF, 0x00, 0x00, 0xFF} // #FF0000
	Lime color.RGBA = color.RGBA{0x00, 0xFF, 0x00, 0xFF} // #00FF00
	Blue color.RGBA = color.RGBA{0x00, 0x00, 0xFF, 0xFF} // #0000FF

	// half strength primary colors
	Maroon   color.RGBA = color.RGBA{0x80, 0x00, 0x00, 0xFF} // #800000
	Green    color.RGBA = color.RGBA{0x00, 0x80, 0x00, 0xFF} // #008000
	NavyBlue color.RGBA = color.RGBA{0x00, 0x00, 0x80, 0xFF} // #000080

	// full strength primary mixes
	Yellow  color.RGBA = color.RGBA{0xFF, 0xFF, 0x00, 0xFF} // #FFFF00
	Aqua    color.RGBA = color.RGBA{0x00, 0xFF, 0xFF, 0xFF} // #00FFFF
	Magenta color.RGBA = color.RGBA{0xFF, 0x00, 0xFF, 0xFF} // #FF00FF

	// half strength primary mixes
	Olive  color.RGBA = color.RGBA{0x80, 0x80, 0x00, 0xFF} // #808000
	Purple color.RGBA = color.RGBA{0x80, 0x00, 0x80, 0xFF} // #800080
	Teal   color.RGBA = color.RGBA{0x00, 0x80, 0x80, 0xFF} // #008080

)

var simplePalette color.Palette = color.Palette{
	Black,
	DarkGray,
	Gray,
	LightGray,
	White,

	Red,
	Lime,
	Blue,

	Maroon,
	Green,
	NavyBlue,

	Yellow,
	Aqua,
	Magenta,

	Olive,
	Purple,
	Teal,
}

func chooseNearest() {
	fileData, _ := os.ReadFile("OneDoesNotSimply_Template.jpg")
	img, _ := jpeg.Decode(bytes.NewReader(fileData))

	bound := img.Bounds()
	palettedImg := image.NewPaletted(bound, simplePalette)
	draw.Draw(palettedImg, bound, img, image.Point{}, draw.Src)

	anim := gif.GIF{}
	anim.Image = append(anim.Image, palettedImg)
	anim.Delay = append(anim.Delay, 100)

	file, _ := os.Create("OneDoesNotSimply_ChooseNearestColor.gif")
	defer file.Close()
	_ = gif.EncodeAll(file, &anim)
}

func chooseNearestPlan9() {
	fileData, _ := os.ReadFile("OneDoesNotSimply_Template.jpg")
	img, _ := jpeg.Decode(bytes.NewReader(fileData))

	bound := img.Bounds()
	palettedImg := image.NewPaletted(bound, palette.Plan9)
	draw.Draw(palettedImg, bound, img, image.Point{}, draw.Src)

	anim := gif.GIF{}
	anim.Image = append(anim.Image, palettedImg)
	anim.Delay = append(anim.Delay, 100)

	file, _ := os.Create("OneDoesNotSimply_ChooseNearestColorPlan9.gif")
	defer file.Close()
	_ = gif.EncodeAll(file, &anim)
}

func useDithering() {
	fileData, _ := os.ReadFile("OneDoesNotSimply_Template.jpg")
	img, _ := jpeg.Decode(bytes.NewReader(fileData))

	bound := img.Bounds()
	palettedImg := image.NewPaletted(bound, simplePalette)
	drawer := draw.FloydSteinberg
	drawer.Draw(palettedImg, bound, img, image.Point{})

	anim := gif.GIF{}
	anim.Image = append(anim.Image, palettedImg)
	anim.Delay = append(anim.Delay, 100)

	file, _ := os.Create("OneDoesNotSimply_UseDithering.gif")
	defer file.Close()
	_ = gif.EncodeAll(file, &anim)
}

func useDitheringPlan9() {
	fileData, _ := os.ReadFile("OneDoesNotSimply_Template.jpg")
	img, _ := jpeg.Decode(bytes.NewReader(fileData))

	bound := img.Bounds()
	palettedImg := image.NewPaletted(bound, palette.Plan9)
	drawer := draw.FloydSteinberg
	drawer.Draw(palettedImg, bound, img, image.Point{})

	anim := gif.GIF{}
	anim.Image = append(anim.Image, palettedImg)
	anim.Delay = append(anim.Delay, 100)

	file, _ := os.Create("OneDoesNotSimply_UseDitheringPlan9.gif")
	defer file.Close()
	_ = gif.EncodeAll(file, &anim)
}

func ditheringArtifactExample() {
	drawCircle := func(c color.Color, size float64) image.Image {
		dc := gg.NewContext(100, 100)
		dc.SetColor(c)
		dc.DrawCircle(50, 50, size)
		dc.Fill()
		return dc.Image()
	}

	frames := []image.Image{}

	blue := color.RGBA{0x00, 0x00, 0xe0, 0xFF}

	for i := 25.0; i < 50; i += .5 {
		frames = append(frames, drawCircle(blue, i))
	}

	for i := 50.0; i > 25; i -= .5 {
		frames = append(frames, drawCircle(blue, i))
	}

	p := color.Palette{easygif.Blue, easygif.NavyBlue, easygif.Black}

	g := easygif.DitheredOptions(frames, time.Millisecond*150, p)
	file, _ := os.Create("ditheringArtifactExample.gif")
	defer file.Close()
	_ = gif.EncodeAll(file, g)
}

func useMostCommonColors() {
	fileData, _ := os.ReadFile("OneDoesNotSimply_Template.jpg")
	img, _ := jpeg.Decode(bytes.NewReader(fileData))
	frames := []image.Image{}
	for i := 0; i < 1; i++ {
		frames = append(frames, img)
	}
	startTime := time.Now()
	_ = easygif.MostCommonColorsWrite(frames, time.Second, "OneDoesNotSimply_UseMostCommonColors.gif")
	fmt.Println("OneDoesNotSimply_UseMostCommonColors.gif", time.Since(startTime))
}

func addTextToCenterOfImage(img image.Image, text string, c color.Color, fontSize float64) image.Image {
	font, err := truetype.Parse(goregular.TTF)
	if err != nil {
		panic("")
	}
	face := truetype.NewFace(font, &truetype.Options{
		Size: fontSize,
	})

	dc := gg.NewContextForImage(img)
	bound := img.Bounds()
	dc.SetFontFace(face)
	dc.SetColor(c)
	dc.DrawStringAnchored(text, float64(bound.Dx())/2, float64(bound.Dy())*.45, 0.5, 0.5)
	// dc.DrawStringAnchored(text, 0, 0, 0.5, 0.5)
	return dc.Image()
}

func easyGifExample() {
	var img image.Image = image.NewRGBA(image.Rect(0, 0, 100, 100))
	imgA := addTextToCenterOfImage(img, "A", easygif.Red, 60)
	imgB := addTextToCenterOfImage(img, "B", easygif.Green, 60)
	imgC := addTextToCenterOfImage(img, "C", easygif.Blue, 60)
	_ = easygif.SaveImageToPNG(imgA, "./imageDirectory/A.png")
	_ = easygif.SaveImageToPNG(imgB, "./imageDirectory/B.png")
	_ = easygif.SaveImageToPNG(imgC, "./imageDirectory/C.png")

	imageDirectory := "./imageDirectory"
	files, _ := os.ReadDir(imageDirectory)
	frames := []image.Image{}
	for _, file := range files {
		fileData, _ := os.ReadFile(path.Join(imageDirectory, file.Name()))
		img, _ := png.Decode(bytes.NewReader(fileData))
		frames = append(frames, img)
	}
	_ = easygif.NearestWrite(frames, time.Millisecond*1000, "easyGif.gif")
}

func EasyGifExample2() {
	frames, _ := easygif.ScreenshotVideo(30, time.Millisecond*100)
	_ = easygif.DitheredWrite(frames, time.Millisecond*100, "easyGif.gif")
}

func ExampleWithOneColorCountEach() {
	img := image.NewRGBA(image.Rect(0, 0, 256, 512))

	var c color.RGBA
	c.A = 0xFF
	for y := 0; y < 512; y++ {
		for x := 0; x < 256; x++ {
			c.B = uint8(x)
			if y <= 255 {
				c.R = uint8(255 - y)
				c.G = 0
			} else {
				c.R = 0
				c.G = uint8(y - 255)
			}
			img.SetRGBA(x, y, c)
		}
	}

	_ = easygif.MostCommonColorsWrite([]image.Image{img}, time.Second, "everyColorOnce.gif")
}

func CheckThatOutliersAreKept() {
	img := image.NewRGBA(image.Rect(0, 0, 512, 512))

	var c color.RGBA
	c.A = 0xFF
	for y := 0; y < 512; y++ {
		for x := 0; x < 512; x++ {
			c.G = uint8(y / 2)
			c.B = uint8(x / 2)
			c.R = 0
			img.SetRGBA(x, y, c)
		}
	}

	img.SetRGBA(255, 255, Red)
	img.SetRGBA(255, 256, Yellow)
	img.SetRGBA(256, 255, Magenta)
	img.SetRGBA(256, 256, Olive)

	_ = easygif.MostCommonColorsWrite([]image.Image{img}, time.Second, "Keep4OutlierColors.gif")
}
