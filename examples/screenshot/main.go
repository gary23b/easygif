package main

import "github.com/gary23b/easygif"

func main() {
	img, _ := easygif.Screenshot()
	_ = easygif.SaveImageToPNG(img, "./examples/screenshot/screenshot.png")

	// trimmed
	img, _ = easygif.ScreenshotTrimmed(200, 200, 200, 600)
	_ = easygif.SaveImageToPNG(img, "./examples/screenshot/screenshotTrimmed.png")
}
