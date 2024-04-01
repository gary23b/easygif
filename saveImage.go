package easygif

import (
	"encoding/gob"
	"image"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
)

func SaveImageToPNG(img image.Image, outputPath string) error {
	out, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer out.Close()

	err = png.Encode(out, img)
	if err != nil {
		return err
	}
	return nil
}

func SaveImageToJPEG(img image.Image, outputPath string) error {
	out, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer out.Close()

	var opt jpeg.Options
	opt.Quality = 100
	err = jpeg.Encode(out, img, &opt)
	if err != nil {
		return err
	}
	return nil
}

func writeGif(g *gif.GIF, filePath string) error {
	// Write the file
	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0o600)
	if err != nil {
		return err
	}
	defer f.Close()

	err = gif.EncodeAll(f, g)
	if err != nil {
		return err
	}
	return nil
}

func SaveFramesToFile(frames []image.Image, filePath string) error {
	rgbaFrames := []image.RGBA{}

	for _, frame := range frames {
		b := frame.Bounds()
		m := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
		draw.Draw(m, m.Bounds(), frame, b.Min, draw.Src)
		rgbaFrames = append(rgbaFrames, *m)
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	enc := gob.NewEncoder(file)
	err = enc.Encode(&rgbaFrames)
	if err != nil {
		return err
	}

	return nil
}

func LoadFramesToFile(filePath string) ([]image.Image, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	frames := []image.RGBA{}

	dec := gob.NewDecoder(file)
	err = dec.Decode(&frames)
	if err != nil {
		return nil, err
	}

	var ret []image.Image
	for i := range frames {
		ret = append(ret, &frames[i])
	}

	return ret, nil
}
