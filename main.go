package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"math/rand"
	"os"
	"os/user"
	"sync"
	"time"
)

const (
	// Position and size
	px   = -0.5557506
	py   = -0.55560
	size = 0.000000001
	//px   = -2
	//py   = -1.2
	//size = 2.5

	// Quality
	imgWidth = 1920
	imgHigh  = 1080
	maxIter  = 1000
	samples  = 25

	showProgress = true
)

func main() {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	filename := usr.HomeDir + "/fractal.png" // Default output filename
	if len(os.Args) == 2 {
		filename = usr.HomeDir + "/" + os.Args[1] // Change filename if given as argument
	}

	log.Println("Allocating image...")
	img := image.NewRGBA(image.Rect(0, 0, imgWidth, imgHigh))

	log.Println("Rendering...")
	start := time.Now()
	render(img)

	log.Println("Done rendering in", time.Since(start))

	log.Println("Encoding image...")
	f, err := os.Create(filename)
	defer f.Close()
	if err != nil {
		panic(err)
	}
	err = png.Encode(f, img)
	if err != nil {
		panic(err)
	}
	log.Println("Done! in", time.Since(start))
}

func render(img *image.RGBA) {
	progress := make(chan struct{})

	// Progress
	if showProgress {
		go func() {
			for i := 1; ; i++ {
				if _, k := <-progress; !k {
					break
				}
				fmt.Printf("\r%d/%d (%d%%)", i, imgWidth, int(100*(float64(i)/float64(imgHigh))))
			}
			fmt.Println()
		}()
	}

	var wg sync.WaitGroup
	for y := 0; y < imgHigh; y++ {
		wg.Add(1)
		go func(y int) {
			for x := 0; x < imgWidth; x++ {
				var sampledColours [samples]color.RGBA
				for i := 0; i < samples; i++ {
					nx := size*((float64(x)+rand.Float64())/float64(imgWidth)) + px
					ny := size*((float64(y)+rand.Float64())/float64(imgHigh)) + py
					sampledColours[i] = paint(mandelbrotIter(nx, ny, maxIter))
				}
				var r, g, b int
				for _, colour := range sampledColours {
					r += int(colour.R)
					g += int(colour.G)
					b += int(colour.B)
				}
				img.SetRGBA(x, y, color.RGBA{
					R: uint8(float64(r) / float64(samples)),
					G: uint8(float64(g) / float64(samples)),
					B: uint8(float64(b) / float64(samples)),
					A: 255,
				})
			}
			if showProgress {
				progress <- struct{}{}
			}
			wg.Done()
		}(y)
	}
	wg.Wait()
	close(progress)
}

func paint(r float64, n int) color.RGBA {
	var insideSet = color.RGBA{R: 255, G: 255, B: 255, A: 255}

	if r > 4 {
		c := hslToRGB(float64(n)/800*r, 1, 0.5)
		return c
	} else {
		return insideSet
	}
}

func mandelbrotIter(px, py float64, maxIter int) (float64, int) {
	var x, y, xx, yy, xy float64

	for i := 0; i < maxIter; i++ {
		xx, yy, xy = x*x, y*y, x*y
		if xx+yy > 4 {
			return xx + yy, i
		}
		x = xx - yy + px
		y = 2*xy + py
	}

	return xx + yy, maxIter
}

// by u/Boraini
//func mandelbrotIterComplex(px, py float64, maxIter int) (float64, int) {
//	var current complex128
//	pxpy := complex(px, py)
//
//	for i := 0; i < maxIter; i++ {
//		magnitude := cmplx.Abs(current)
//		if magnitude > 2 {
//			return magnitude * magnitude, i
//		}
//		current = current * current + pxpy
//	}
//
//	magnitude := cmplx.Abs(current)
//	return magnitude * magnitude, maxIter
//}
