package easygif

import (
	"image"
	"image/color"
	"math"
	"sort"
	"sync"
)

type colorAvg struct {
	count int

	r int
	g int
	b int

	rSum int
	gSum int
	bSum int
}

// Compute the distance squared between two colors
func newColorAvg(c color.RGBA, usageCount int) colorAvg {
	newItem := colorAvg{
		count: usageCount,

		r: int(c.R),
		g: int(c.G),
		b: int(c.B),

		rSum: usageCount * int(c.R),
		gSum: usageCount * int(c.G),
		bSum: usageCount * int(c.B),
	}
	return newItem
}

// Compute the distance squared between two colors
func (s *colorAvg) distanceSqrd(c2 *colorAvg) float64 {
	dR := float64(s.r - c2.r)
	dG := float64(s.g - c2.g)
	dB := float64(s.b - c2.b)
	distSqrd := dR*dR + dG*dG + dB*dB
	return distSqrd
}

// effectively merges the two colors. c2 gets zeroed out.
func (s *colorAvg) consumeOtherColor(c2 *colorAvg) {
	s.rSum += c2.rSum
	s.gSum += c2.gSum
	s.bSum += c2.bSum
	s.count += c2.count

	invCount := 1.0 / float64(s.count)
	s.r = int(float64(s.rSum)*invCount + .5)
	s.g = int(float64(s.gSum)*invCount + .5)
	s.b = int(float64(s.bSum)*invCount + .5)

	*c2 = colorAvg{}
}

// re-compute internal state
func (s *colorAvg) getColor() color.RGBA {
	ret := color.RGBA{
		R: uint8(s.r),
		G: uint8(s.g),
		B: uint8(s.b),
		A: 0xff,
	}

	return ret
}

func FindMostCommonColors(frames []image.Image) []color.Color {
	colorCount := getColorHistogram(frames)

	// Convert to a list
	foundColors := make([]colorAvg, 0, len(colorCount))
	for i, j := range colorCount {
		foundColors = append(foundColors, newColorAvg(i, j))
	}

	if len(foundColors) > 256 {
		foundColors = trimDownCommonColorList(foundColors)
	}

	ret := make([]color.Color, 0, 256)

	for _, c := range foundColors {
		ret = append(ret, c.getColor())
	}
	return ret
}

func trimDownCommonColorList(sortedColors []colorAvg) []colorAvg {
	// First sort.
	sort.SliceStable(sortedColors, func(i, j int) bool {
		return sortedColors[i].count > sortedColors[j].count
	})

	// Merge near colors, first nearest, then expanding to a distance of 4
	sortedColors = combineNearColors(sortedColors, 1.0)
	sortedColors = combineNearColors(sortedColors, 1.5)
	sortedColors = combineNearColors(sortedColors, 2.0)
	sortedColors = combineNearColors(sortedColors, 3.0)
	sortedColors = combineNearColors(sortedColors, 4.0)

	// Sort again.
	sort.SliceStable(sortedColors, func(i, j int) bool {
		return sortedColors[i].count > sortedColors[j].count
	})

	// Trim off all the zero counts:
	trimIndex := len(sortedColors) - 1
	for ; trimIndex > 1; trimIndex-- {
		c := &sortedColors[trimIndex]
		if c.count > 0 {
			break
		}
	}
	sortedColors = sortedColors[:trimIndex+1]
	// log.Println("Number of colors after trimming", len(sortedColors))

	// Now walk through from least common color to most common color.
	// If there is a color within a distance of 50, merge it.
	const maxDistSqrd = 50 * 50
	sortEveryN := int(float64(len(sortedColors)) * .01)
	IterationsAfterSort := 0
	outlierColors := []colorAvg{}

	for i := len(sortedColors) - 1; len(sortedColors)+len(outlierColors) > 256; i-- {
		c := &sortedColors[i]

		bestIndex := 0
		bestDistanceSqrd := math.MaxFloat64
		for j := 0; j < i; j++ {
			c2 := sortedColors[j]
			distSqrd := c.distanceSqrd(&c2)
			if distSqrd < bestDistanceSqrd {
				bestDistanceSqrd = distSqrd
				bestIndex = j
			}
		}

		if bestDistanceSqrd < maxDistSqrd {
			sortedColors[bestIndex].consumeOtherColor(c)
		} else {
			outlierColors = append(outlierColors, *c)
		}

		sortedColors = sortedColors[:i]

		IterationsAfterSort++

		if IterationsAfterSort >= sortEveryN {
			IterationsAfterSort = 0

			sort.SliceStable(sortedColors, func(i, j int) bool {
				return sortedColors[i].count > sortedColors[j].count
			})
		}
	}

	// Add the outlier colors back into the array.
	// log.Println("outlier color count", len(outlierColors))
	sortedColors = append(sortedColors, outlierColors...)

	// final Sort.
	sort.SliceStable(sortedColors, func(i, j int) bool {
		return sortedColors[i].count > sortedColors[j].count
	})

	return sortedColors
}

func combineNearColors(sortedColors []colorAvg, mergeDistance float64) []colorAvg {
	// log.Println("initial Colors", len(sortedColors))
	mergeDistSqrd := mergeDistance * mergeDistance
	mergedColors := 0
	for i := range sortedColors {
		c := &sortedColors[i]
		if c.count == 0 {
			continue
		}

		for j := i + 1; j < len(sortedColors); j++ {
			c2 := &sortedColors[j]
			if c2.count == 0 {
				continue
			}
			dist := c.distanceSqrd(c2)
			if dist <= mergeDistSqrd {
				c.consumeOtherColor(c2)
				mergedColors++
			}
		}
	}
	// log.Println("Merged Colors", mergedColors)

	return sortedColors
}

// Use a worker group to search through all the frames and count up all the colors.
func getColorHistogram(frames []image.Image) map[color.RGBA]int {
	workerFunc := func(wg *sync.WaitGroup, frames chan image.Image, ret chan map[color.RGBA]int) {
		colorCount := make(map[color.RGBA]int, 10000)
		for frame := range frames {
			bounds := frame.Bounds()
			for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
				for x := bounds.Min.X; x < bounds.Max.X; x++ {
					r, g, b, a := frame.At(x, y).RGBA()
					c := color.RGBA{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), uint8(a >> 8)}
					count := colorCount[c]
					count++
					colorCount[c] = count
				}
			}
		}
		ret <- colorCount
		wg.Done()
	}

	workerCount := 8
	requestChan := make(chan image.Image, 100)
	resultChan := make(chan map[color.RGBA]int, workerCount)
	wg := &sync.WaitGroup{}
	wg.Add(workerCount)
	for i := 0; i < workerCount; i++ {
		go workerFunc(wg, requestChan, resultChan)
	}

	for _, frame := range frames {
		requestChan <- frame
	}
	close(requestChan)
	wg.Wait()
	close(resultChan)

	colorCount := make(map[color.RGBA]int, 10000)
	for ret := range resultChan {
		for c, val := range ret {
			count := colorCount[c]
			count += val
			colorCount[c] = count
		}
	}

	return colorCount
}
