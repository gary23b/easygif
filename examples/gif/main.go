package main

import (
	"fmt"
	"image"
	"image/color"
	"time"

	"github.com/fogleman/gg"
	"github.com/gary23b/easygif"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font/gofont/goregular"
)

func main() {

}

func UseEasyGif() {
	// frames, _ := easygif.ScreenshotVideo(100, time.Millisecond*50)
	// _ = easygif.EasyGifWrite(frames, time.Millisecond*100, "./examples/gif/screenshot.gif")

	// time.Sleep(time.Second * 3)
	// fmt.Println("Set?")
	time.Sleep(time.Second * 3)
	fmt.Println("GO!")
	frames, _ := easygif.ScreenshotVideoTrimmed(30, time.Millisecond*50, 150, 1050, 380, 1270)
	fmt.Println("Collection Done.")
	fmt.Println("Adding Text.")
	s1 := "ONE DOES NOT SIMPLY"
	s2 := "MAKE A GIF"
	AddMemeText(frames, s1, s2, easygif.Crimson)

	fmt.Println("Encoding GIF.")

	startTime := time.Now()
	_ = easygif.NearestWrite(frames, time.Millisecond*100, "./examples/gif/OneDoesNotSimplyMakeAGIF_Nearest.gif")
	fmt.Println("Took: ", time.Since(startTime), " to encode Nearest.")

	startTime = time.Now()
	_ = easygif.DitheredWrite(frames, time.Millisecond*100, "./examples/gif/OneDoesNotSimplyMakeAGIF_Dithered.gif")
	fmt.Println("Took: ", time.Since(startTime), " to encode Dithered.")

	startTime = time.Now()
	_ = easygif.MostCommonColorsWrite(frames, time.Millisecond*100, "./examples/gif/OneDoesNotSimplyMakeAGIF_MostCommon.gif")
	fmt.Println("Took: ", time.Since(startTime), " to encode MostCommon.")
}

func AddMemeText(frames []image.Image, s1, s2 string, c color.Color) {
	fontSize := 60.0
	font, err := truetype.Parse(goregular.TTF)
	if err != nil {
		panic("")
	}
	face := truetype.NewFace(font, &truetype.Options{
		Size: fontSize,
	})

	for i := range frames {
		frame := frames[i]
		dc := gg.NewContextForImage(frame)
		bound := frame.Bounds()
		dc.SetFontFace(face)
		dc.SetColor(c)
		dc.DrawStringAnchored(s1, float64(bound.Dx())/2, float64(bound.Dy())*.10, 0.5, 0.5)
		dc.DrawStringAnchored(s2, float64(bound.Dx())/2, float64(bound.Dy())*.90, 0.5, 0.5)

		frames[i] = dc.Image()
	}
}

// func makeGif(frames []image.Image) {
// 	s1 := "ONE DOES NOT SIMPLY"
// 	s2 := "MAKE A GIF"
// 	AddMemeText(frames, s1, s2, easygif.Crimson)

// 	hundredthOfSecondDelay := 10

// 	// Process the images.
// 	imagesPal := make([]*image.Paletted, 0, len(frames))
// 	delays := make([]int, 0, len(frames))

// 	// Fill the request channel with images to convert
// 	for frameIndex := range frames {
// 		screenShot := frames[frameIndex]
// 		bounds := screenShot.Bounds()
// 		ssPaletted := image.NewPaletted(bounds, palette.Plan9)
// 		imagesPal = append(imagesPal, ssPaletted)
// 		delays = append(delays, hundredthOfSecondDelay)
// 		draw.FloydSteinberg.Draw(ssPaletted, bounds, screenShot, image.Point{})
// 	}

// 	// Create the GIF struct and write it to a file.
// 	g := &gif.GIF{
// 		Image: imagesPal,
// 		Delay: delays,
// 	}
// 	file, _ := os.Create("OneDoesNotSimply_try5.gif")
// 	defer file.Close()
// 	_ = gif.EncodeAll(file, g)
// }

// func main() {
// 	time.Sleep(time.Second * 3)
// 	fmt.Println("GO!")
// 	frames, _ := easygif.ScreenshotVideoTrimmed(30, time.Millisecond*50, 150, 1050, 380, 1270)
// 	fmt.Println("Collection Done.")

// 	makeGif(frames)
// }
