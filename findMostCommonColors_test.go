package easygif

import (
	"bytes"
	"image"
	"image/jpeg"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFindMostCommonColors(t *testing.T) {
	t.SkipNow()
	fileData, _ := os.ReadFile("./examples/compare/OneDoesNotSimply_Template.jpg")
	img, _ := jpeg.Decode(bytes.NewReader(fileData))
	frames := []image.Image{img}

	a := FindMostCommonColors(frames)
	b := FindMostCommonColors(frames)

	require.Equal(t, a[:5], b[:5])
}
