package main

import (
	"flag"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/muesli/termenv"
	asciiconvertor "heartalborada.fun/VideoAsciiPlayer/asciiConvertor"
	"heartalborada.fun/VideoAsciiPlayer/terminal"
)

func main() {
	var frameBufferSize = flag.Int("buffer", 120, "The Size of frame buffer [Recommended above 120]")
	var videoImgPath = flag.String("path", "", "The Path of Stroage Video Every Frame Images")
	var videoImgSize = flag.Int64("count", -1, "The Count of Frame Image [From 0 to Your setting]")
	var renderPoolSize = flag.Int("pool", 16, "Render Pool Size")
	var targetFramePreSecond = flag.Int("fps", 30, "Target FPS")
	flag.Parse()
	if *videoImgPath == "" || *videoImgSize <= 0 || *frameBufferSize <= 1 || *renderPoolSize <= 0 || *targetFramePreSecond <= 0 {
		flag.PrintDefaults()
		return
	}
	term := terminal.NewTerminal()
	restoreConsole, err := termenv.EnableVirtualTerminalProcessing(termenv.DefaultOutput())
	if err != nil {
		panic(err)
	}
	defer restoreConsole()
	w, h, _ := term.GetScreenSize()
	var wg sync.WaitGroup
	frameCh := make(chan struct {
		index int
		frame string
	}, *frameBufferSize)
	quitCh := make(chan struct{})
	wg.Add(1)
	go func() {
		defer wg.Done()
		frameMap := make(map[int]string)
		nextFrame := 1
		start := time.Now()
		frameInterval := time.Second / time.Duration(*targetFramePreSecond)
		accumulatedDelay := time.Duration(0)
		for {
			select {
			case f := <-frameCh:
				frameMap[f.index] = f.frame
				for {
					if frame, ok := frameMap[nextFrame]; ok {
						os.Stdout.WriteString("\033[H")
						os.Stdout.WriteString(frame)
						delete(frameMap, nextFrame)
						frameElapsed := time.Since(start) + accumulatedDelay
						if sleepDuration := frameInterval - frameElapsed; sleepDuration > 0 {
							end := time.Now().Add(sleepDuration)
							for time.Now().Before(end) {
								// busy wait
							}
							accumulatedDelay = 0
						} else {
							accumulatedDelay = -sleepDuration
						}
						frameSpeed := 1e6 / float64(time.Since(start).Microseconds())
						os.Stdout.WriteString("Frame Speed: ")
						os.Stdout.WriteString(strconv.FormatFloat(frameSpeed, 'f', 2, 64))
						os.Stdout.WriteString(" fps, Current Frame: ")
						os.Stdout.WriteString(strconv.FormatInt(int64(nextFrame), 10))
						start = time.Now()
						nextFrame++
					} else {
						if int64(nextFrame) >= *videoImgSize+1 {
							close(quitCh)
						}
						break
					}
				}
			case <-quitCh:
				return
			}
		}
	}()

	sem := make(chan struct{}, *renderPoolSize)

	var charWidth float64 = term.GetCharWidth()
	for i := 1; int64(i) <= *videoImgSize; i++ {
		wg.Add(1)
		sem <- struct{}{} // Acquire a slot
		go func(i int) {
			defer wg.Done()
			defer func() { <-sem }() // Release the slot
			inputPath := *videoImgPath + "/" + strconv.Itoa(i) + ".png"
			inputFile, err := os.Open(inputPath)
			if err != nil {
				panic(fmt.Sprintf("Error opening file: %v\n", err))
			}

			img, _, err := image.Decode(inputFile)
			inputFile.Close()
			if err != nil {
				panic(fmt.Sprintf("Error decoding image: %v\n", err))
			}
			s := asciiconvertor.ConverImage2Ascii(img, w, h-2, charWidth)

			frameCh <- struct {
				index int
				frame string
			}{i, s}
		}(i)
	}
	wg.Wait()
}
