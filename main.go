package main

import (
	"fluorescence/geometry"
	"fmt"
	"image"
	"image/png"
	"math/rand"
	"os"
	"time"
)

func main() {

	// maxThreads := int64(runtime.NumCPU())
	// maxThreads := int64(runtime.NumCPU() * 10)
	// maxThreads := int64(runtime.NumCPU() * 1000)
	// maxThreads := int64(1)
	// get parameters
	parametersFileName := "./config/parameters.json"
	camerasFileName := "./config/cameras.json"
	objectsFileName := "./config/objects.json"
	materialsFileName := "./config/materials.json"
	fmt.Printf("Loading Config files...\n")
	parameters, err := LoadConfigs(parametersFileName, camerasFileName, objectsFileName, materialsFileName)
	if err != nil {
		fmt.Printf("Error loading parameters data: %s\n", err.Error())
		return
	}

	// create image
	fmt.Printf("Creating in-mem image...\n")
	img := image.NewRGBA64(image.Rect(0, 0, parameters.ImageWidth, parameters.ImageHeight))

	// fill image
	fmt.Printf("Filling in-mem image...\n")

	// spew.Dump(parameters.Scene.Objects)

	// wg := sync.WaitGroup{}
	pixelCount := parameters.ImageWidth * parameters.ImageHeight
	// sem := semaphore.NewWeighted(maxThreads)
	doneChan := make(chan int, pixelCount)
	// rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	startTime := time.Now()
	for y := 0; y < parameters.ImageHeight; y++ {
		// sem.Acquire(context.Background(), 1)
		// wg.Add(1)
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		go func(y int, rng *rand.Rand, dc chan<- int) {
			// defer wg.Done()
			for x := 0; x < parameters.ImageWidth; x++ {
				// defer sem.Release(1)
				colorAccumulator := geometry.VECTOR_ZERO
				for s := 0; s < parameters.SampleCount; s++ {
					u := (float64(x) + rng.Float64()) / float64(parameters.ImageWidth)
					v := (float64(y) + rng.Float64()) / float64(parameters.ImageHeight)

					ray := parameters.Scene.Camera.GetRay(u, v, rng)

					tempColor := colorOf(parameters, ray, rng, 0)
					colorAccumulator = colorAccumulator.Add(tempColor)
				}
				colorAccumulator = colorAccumulator.DivScalar(float64(parameters.SampleCount)).Clamp(0, 1).Pow(1.0 / float64(parameters.GammaCorrection))
				color := colorAccumulator.ToColor()
				// pixelsChan <- geometry.Pixel{x, parameters.ImageHeight - y - 1, *color}

				img.SetRGBA64(x, parameters.ImageHeight-y-1, color.ToRGBA64())
				dc <- 1
			}
			// fmt.Printf("ok\n")
		}(y, r, doneChan)
		// if y%10 == 0 {
		// 	fmt.Printf("\t\t%3.4f%%\n", 100*float64(y)/float64(parameters.ImageHeight))
		// }
	}
	// var p geometry.Pixel
	// for i := 0; i < parameters.ImageWidth*parameters.ImageHeight; i++ {
	// 	p = <-pixelsChan
	// 	img.SetRGBA64(p.X, p.Y, *p.Color.ToRGBA64())
	// }
	// fmt.Printf("Waiting on threads...\n")
	doneCount := 0
	printInterval := pixelCount / 1000
	for i := 0; i < pixelCount; i++ {
		<-doneChan
		doneCount++
		if pixelCount > 1000 && doneCount%printInterval == 0 {
			elapsedTime := time.Since(startTime)
			estimatedTime := time.Duration(float64(elapsedTime) * (float64(pixelCount) / float64(doneCount)))
			remainingTime := estimatedTime - elapsedTime
			fmt.Printf("\t\t%5.1f%% - Est. Rem: ~%v,\tTotal: ~%v\n", 100*float64(doneCount)/float64(pixelCount), remainingTime, estimatedTime)
		}
	}
	// wg.Wait()
	// sem.Release(0)
	totalDuration := time.Since(startTime)
	fmt.Printf("\tTotal time: %v\n", totalDuration)

	// create file
	fmt.Printf("Creating image file...\n")
	file, err := getImageFile(parameters)
	if err != nil {
		fmt.Printf("Error creating image file: %s\n", err.Error())
		return
	}
	defer file.Close()

	// encode image to file
	fmt.Printf("Writing in-mem image to image file...\n")
	err = png.Encode(file, img)
	if err != nil {
		fmt.Printf("Error encoding to image file: %s\n", err.Error())
		return
	}
	fmt.Printf("Done!\n")
	return
}

func colorOf(parameters *Parameters, r geometry.Ray, rng *rand.Rand, depth int) geometry.Vector {

	backgroundColor := geometry.Vector{
		X: parameters.BackgroundColor.Red,
		Y: parameters.BackgroundColor.Green,
		Z: parameters.BackgroundColor.Blue,
	}

	if depth > parameters.MaxBounces {
		return backgroundColor
	}
	rayHit, hitSomething := parameters.Scene.Objects.Intersection(r, parameters.TMin, parameters.TMax)
	if !hitSomething {
		return backgroundColor
	}

	mat := rayHit.Material

	if mat.Reflectance() == geometry.VECTOR_ZERO {
		return mat.Emittance()
	}

	scatteredRay, wasScattered := rayHit.Material.Scatter(*rayHit, rng)
	if !wasScattered {
		return backgroundColor
	}
	incomingColor := colorOf(parameters, scatteredRay, rng, depth+1)
	return mat.Emittance().Add(mat.Reflectance().MultVector(incomingColor))
}

func getImageFile(parameters *Parameters) (*os.File, error) {
	filename := fmt.Sprintf(
		"%s%s_%ds_%s.%s",
		parameters.FileDirectory,
		parameters.FileName,
		parameters.SampleCount,
		time.Now().Format("2006-01-02_T150405"),
		parameters.FileType)
	os.MkdirAll(parameters.FileDirectory, os.ModePerm)
	return os.Create(filename)
}
