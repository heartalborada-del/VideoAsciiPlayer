package main

import (
	"fmt"
	"image/png"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/muesli/termenv"
	asciiconvertor "heartalborada.fun/VideoAsciiPlayer/asciiConvertor"
	"heartalborada.fun/VideoAsciiPlayer/terminal"
)

func main() {
	term := terminal.NewTerminal()
	if term.IsWindows() {
		restoreConsole, err := termenv.EnableVirtualTerminalProcessing(termenv.DefaultOutput())
		if err != nil {
			panic(err)
		}
		defer restoreConsole()
	}
	w, h, _ := term.GetScreenSize()

	var wg sync.WaitGroup
	frameCh := make(chan struct {
		index int
		frame string
	}, 120)
	quitCh := make(chan struct{})
	wg.Add(1)
	go func() {
		defer wg.Done()
		frameMap := make(map[int]string)
		nextFrame := 1
		start := time.Now()
		frameInterval := time.Second / 30
		const offset time.Duration = time.Microsecond * 500 //500 us
		for {
			select {
			case f := <-frameCh:
				frameMap[f.index] = f.frame
				for {
					if frame, ok := frameMap[nextFrame]; ok {
						os.Stdout.WriteString("\033[H")
						os.Stdout.WriteString(frame)
						delete(frameMap, nextFrame)
						frameElapsed := time.Since(start)
						if sleepDuration := frameInterval - frameElapsed; sleepDuration > 0 {
							//time.Sleep(sleepDuration)
							//IDK WHY time.Sleep CANT to SLEEP Accrate TIME
							end := time.Now().Add(sleepDuration - offset)
							for time.Now().Before(end) {
								// busy wait
							}
						}
						frameSpeed := 1e6 / float64(time.Since(start).Microseconds())
						var builder strings.Builder
						builder.WriteString("Frame Speed: ")
						builder.WriteString(strconv.FormatFloat(frameSpeed, 'f', 2, 64))
						builder.WriteString(" fps, Current Frame: ")
						builder.WriteString(strconv.FormatInt(int64(nextFrame), 10))
						os.Stdout.WriteString(builder.String())
						start = time.Now()
						nextFrame++
					} else {
						if nextFrame >= 6301 {
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

	sem := make(chan struct{}, 16)

	var charWidth float64 = term.GetCharWidth()
	for i := 1; i <= 6300; i++ {
		wg.Add(1)
		sem <- struct{}{} // Acquire a slot
		go func(i int) {
			defer wg.Done()
			defer func() { <-sem }() // Release the slot
			inputPath := "1/" + strconv.Itoa(i) + ".png"
			inputFile, err := os.Open(inputPath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error opening file: %v\n", err)
				return
			}

			img, err := png.Decode(inputFile)
			inputFile.Close()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error decoding image: %v\n", err)
				return
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
