package main

import (
	"fmt"
	"time"

	"github.com/gary23b/easygif"
)

func main() {
	frames, _ := easygif.ScreenshotVideoTrimmed(30, time.Millisecond*50, 150, 1050, 380, 1270)
	_ = easygif.NearestWrite(frames, time.Millisecond*100, "./examples/globsave/globsave1.gif")

	err := easygif.SaveFramesToFile(frames, "./examples/globsave/save.bin")
	if err != nil {
		panic(err)
	}

	time.Sleep(time.Second)
	frames, err = easygif.LoadFramesToFile("./examples/globsave/save.bin")
	// frames, err := easygif.LoadFramesToFile("./save.bin")
	if err != nil {
		panic(err)
	}
	fmt.Println("frames", len(frames))

	err = easygif.NearestWrite(frames, time.Millisecond*100, "./examples/globsave/globsave2.gif")
	if err != nil {
		panic(err)
	}
}
