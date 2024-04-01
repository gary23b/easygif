package easygif

import (
	"fmt"
	"image"
	"time"

	"github.com/vova616/screenshot"
)

func Screenshot() (image.Image, error) {
	img, err := screenshot.CaptureScreen()
	if err != nil {
		return nil, fmt.Errorf("TakeScreenshot(), error while getting screenshot: %w", err)
	}
	return img, nil
}

func ScreenshotTrimmed(leftTrim, rightTrim, topTrim, bottomTrim int) (image.Image, error) {
	captureRec, err := getTrimmedRectangle(leftTrim, rightTrim, topTrim, bottomTrim)
	if err != nil {
		return nil, fmt.Errorf("TakeScreenshot(): %w", err)
	}

	img, err := screenshot.CaptureRect(captureRec)
	if err != nil {
		return nil, fmt.Errorf("TakeScreenshot(), error while getting screenshot: %w", err)
	}
	return img, nil
}

func getTrimmedRectangle(leftTrim, rightTrim, topTrim, bottomTrim int) (image.Rectangle, error) {
	if leftTrim < 0 || rightTrim < 0 || topTrim < 0 || bottomTrim < 0 {
		return image.Rectangle{}, fmt.Errorf("getTrimmedRectangle(), all trim values must be >= 0")
	}

	screenRec, err := screenshot.ScreenRect()
	if err != nil {
		return image.Rectangle{}, fmt.Errorf("getTrimmedRectangle(), error while getting ScreenRect: %w", err)
	}

	totalTrimX := leftTrim + rightTrim
	totalTrimY := topTrim + bottomTrim
	if totalTrimX >= screenRec.Dx() || totalTrimY >= screenRec.Dy() {
		return image.Rectangle{}, fmt.Errorf("getTrimmedRectangle(), total trim is larger than screen size")
	}

	captureRec := screenRec
	captureRec.Min.X += leftTrim
	captureRec.Max.X -= rightTrim
	captureRec.Min.Y += topTrim
	captureRec.Max.Y -= bottomTrim

	return captureRec, nil
}

func ScreenshotVideo(
	frameCount int,
	delayBetweenScreenshots time.Duration,
) ([]image.Image, error) {
	return ScreenshotVideoTrimmed(frameCount, delayBetweenScreenshots, 0, 0, 0, 0)
}

func ScreenshotVideoTrimmed(
	frameCount int,
	delayBetweenScreenshots time.Duration,
	leftTrim, rightTrim, topTrim, bottomTrim int,
) ([]image.Image, error) {
	//
	captureRec, err := getTrimmedRectangle(leftTrim, rightTrim, topTrim, bottomTrim)
	if err != nil {
		return nil, fmt.Errorf("ScreenshotVideoTrimmed(): %w", err)
	}

	frames := make([]image.Image, frameCount)

	nextTime := time.Now()
	for frameIndex := 0; frameIndex < frameCount; frameIndex++ {
		frames[frameIndex], err = screenshot.CaptureRect(captureRec)
		if err != nil {
			return nil, fmt.Errorf("TakeScreenshot(), error while getting screenshot: %w", err)
		}
		nextTime = nextTime.Add(delayBetweenScreenshots)
		time.Sleep(time.Until(nextTime))
	}

	return frames, nil
}
