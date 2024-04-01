package easygif

import (
	"image"
	"image/color"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"log"
	"sync"
	"time"
)

// //////////////////////////////////////////////
// //////////////////////////////////////////////
// //////////////////////////////////////////////

func MostCommonColors(
	frames []image.Image,
	timeBetweenFrames time.Duration,
) *gif.GIF {
	mostCommonColors := FindMostCommonColors(frames)

	return NearestOptions(frames, timeBetweenFrames, mostCommonColors)
	// return CreateDitheredGif(frames, timeBetweenFrames, mostCommonColors)
}

func MostCommonColorsWrite(
	frames []image.Image,
	timeBetweenFrames time.Duration,
	outputGifFilePath string,
) error {
	g := MostCommonColors(frames, timeBetweenFrames)
	return writeGif(g, outputGifFilePath)
}

// //////////////////////////////////////////////
// //////////////////////////////////////////////
// //////////////////////////////////////////////

// Create a GIF with the Plan9 palette of colors. If an exact match is not found, use the nearest available color.
func Nearest(
	frames []image.Image,
	timeBetweenFrames time.Duration,
) *gif.GIF {
	return NearestOptions(frames, timeBetweenFrames, palette.Plan9)
}

// Create a GIF with the Plan9 palette of colors. If an exact match is not found, use the nearest available color.
func NearestWrite(
	frames []image.Image,
	timeBetweenFrames time.Duration,
	outputGifFilePath string,
) error {
	g := Nearest(frames, timeBetweenFrames)
	return writeGif(g, outputGifFilePath)
}

// Create a GIF with the given palette of colors. If an exact match is not found, use the nearest available color.
func NearestOptions(
	frames []image.Image,
	timeBetweenFrames time.Duration,
	colorPalette color.Palette,
) *gif.GIF {
	//
	hundredthOfSecondDelay := int(timeBetweenFrames.Seconds() * 100)

	// Process the images.
	imagesPal := make([]*image.Paletted, 0, len(frames))
	delays := make([]int, 0, len(frames))

	// Create 10 workers
	requestChan := make(chan palettedImageProcessorRequest, 100)
	wg := &sync.WaitGroup{}
	wg.Add(10)
	go func() {
		for i := 0; i < 10; i++ {
			go gifPalettedImageProcessor(wg, requestChan)
		}
	}()

	// Fill the request channel with images to convert
	for frameIndex := range frames {
		screenShot := frames[frameIndex]
		bounds := screenShot.Bounds()
		ssPaletted := image.NewPaletted(bounds, colorPalette)
		imagesPal = append(imagesPal, ssPaletted)
		delays = append(delays, hundredthOfSecondDelay)

		// // All this additional logic to speed up the following commented lines. which takes a couple seconds per frame
		// for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		// 	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		// 		ssPaletted.Set(x, y, screenShot.At(x, y))
		// 	}
		// }

		// // calling convertToPalettedWithCacheRGBA takes 2.1s for 95 images. vs 684ms with 10 workers. 3.12 times slower
		// srcRGBA, _ := screenShot.(*image.RGBA)
		// convertToPalettedWithCacheRGBA(palettedCacheRGBA, srcRGBA, ssPaletted)

		newRequest := palettedImageProcessorRequest{src: screenShot, dest: ssPaletted}
		requestChan <- newRequest
	}
	// Close the channel and wait for all workers to finish.
	close(requestChan)
	wg.Wait()

	ret := &gif.GIF{
		Image: imagesPal,
		Delay: delays,
	}

	return ret
}

// takes 0.12s on average.
func convertToPalettedWithCache(palettedCache map[color.Color]uint8, src image.Image, dest *image.Paletted) {
	// startTime := time.Now()

	if src.Bounds() != dest.Bounds() {
		log.Println("src and dest do not have the same rectangle")
		return
	}

	bounds := src.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := src.At(x, y)

			// Get the palette color index for this RGBA color
			palettedIndex, ok := palettedCache[c]
			if !ok {
				palettedIndex = uint8(dest.Palette.Index(c))
				palettedCache[c] = palettedIndex
			}
			// dest.Set(x, y, c)
			i := dest.PixOffset(x, y)
			dest.Pix[i] = palettedIndex
		}
	}

	// deltaTime := time.Since(startTime)
	// fmt.Println(deltaTime.Seconds())
}

// takes 0.065s on average.
func convertToPalettedWithCacheRGBA(palettedCache map[color.RGBA]uint8, src *image.RGBA, dest *image.Paletted) {
	// startTime := time.Now()

	if src.Bounds() != dest.Bounds() {
		log.Println("src and dest do not have the same rectangle")
		return
	}

	bounds := src.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			// Unroll the operations of: dest.Set(x, y, src.At(x, y))
			// first, get the src color: srcImage.At(x, y)
			i := (y-src.Rect.Min.Y)*src.Stride + (x-src.Rect.Min.X)*4
			// This appears to be called "Full slice expressions". a[low : high : max]. It sets the new slice capacity to max-low
			s := src.Pix[i : i+4 : i+4]
			c := color.RGBA{s[0], s[1], s[2], s[3]}

			// Get the palette color index for this RGBA color
			palettedIndex, ok := palettedCache[c]
			if !ok {
				palettedIndex = uint8(dest.Palette.Index(c))
				palettedCache[c] = palettedIndex
			}
			// dest.Set(x, y, c)
			i = (y-dest.Rect.Min.Y)*dest.Stride + (x-dest.Rect.Min.X)*1
			dest.Pix[i] = palettedIndex
		}
	}

	// deltaTime := time.Since(startTime)
	// fmt.Println(deltaTime.Seconds())
}

type palettedImageProcessorRequest struct {
	src  image.Image
	dest *image.Paletted
}

// Has a 2x speed improvement when the src image is an image.RGBA
// The chosen palette color for a given src color is saved in a cache.
func gifPalettedImageProcessor(
	wg *sync.WaitGroup,
	requestChan chan palettedImageProcessorRequest,
) {
	palettedCacheColor := make(map[color.Color]uint8, 100)
	palettedCacheRGBA := make(map[color.RGBA]uint8, 100)
	for {
		request, ok := <-requestChan
		if !ok {
			break
		}

		srcRGBA, ok := request.src.(*image.RGBA)
		if ok {
			convertToPalettedWithCacheRGBA(palettedCacheRGBA, srcRGBA, request.dest)
		} else {
			convertToPalettedWithCache(palettedCacheColor, request.src, request.dest)
		}
	}

	wg.Done()
}

// //////////////////////////////////////////////
// //////////////////////////////////////////////
// //////////////////////////////////////////////

// uses draw.FloydSteinberg.Draw() and the Plan9 Palette
// Much slower, but is able to approximate colors much better that just a nearest fit match.
func Dithered(
	frames []image.Image,
	timeBetweenFrames time.Duration,
) *gif.GIF {
	return DitheredOptions(frames, timeBetweenFrames, palette.Plan9)
}

// uses draw.FloydSteinberg.Draw() and the Plan9 Palette
// Much slower, but is able to approximate colors much better that just a nearest fit match.
func DitheredWrite(
	frames []image.Image,
	timeBetweenFrames time.Duration,
	outputGifFilePath string,
) error {
	g := Dithered(frames, timeBetweenFrames)

	// Write the file
	return writeGif(g, outputGifFilePath)
}

// uses draw.FloydSteinberg.Draw()
// Much slower, but is able to approximate colors much better that just a nearest fit match.
func DitheredOptions(
	frames []image.Image,
	timeBetweenFrames time.Duration,
	colorPalette color.Palette,
) *gif.GIF {
	//
	hundredthOfSecondDelay := int(timeBetweenFrames.Seconds() * 100)

	// Process the images.
	imagesPal := make([]*image.Paletted, 0, len(frames))
	delays := make([]int, 0, len(frames))

	wp := newWorkerPool(10)

	// Fill the request channel with images to convert
	for frameIndex := range frames {
		screenShot := frames[frameIndex]
		bounds := screenShot.Bounds()
		ssPaletted := image.NewPaletted(bounds, colorPalette)
		imagesPal = append(imagesPal, ssPaletted)
		delays = append(delays, hundredthOfSecondDelay)

		wp.addTask(func() {
			draw.FloydSteinberg.Draw(ssPaletted, bounds, screenShot, image.Point{})
		})
	}

	wp.waitForCompletion()

	ret := &gif.GIF{
		Image: imagesPal,
		Delay: delays,
	}

	return ret
}

///////////////////////////////////
///////////////////////////////////
