package main

import (
	"fmt"
	"image/png"
	"os"
	"strconv"
	"time"

	asciiconvertor "heartalborada.fun/VideoAsciiPlayer/asciiConvertor"
	"heartalborada.fun/VideoAsciiPlayer/terminal"
)

func single() {
	terminal := terminal.NewTerminal()
	w, h, _ := terminal.GetScreenSize()
	for i := 1; i <= 6300; i++ {
		startMicroSec := time.Now().UnixMicro()
		inputPath := "1/" + strconv.Itoa(i) + ".png"
		inputFile, err := os.Open(inputPath)
		if err != nil {
			panic(err)
		}
		img, err := png.Decode(inputFile)
		if err != nil {
			panic(err)
		}
		s := asciiconvertor.ConverImage2Ascii(img, w, h-1, terminal.GetCharWidth())
		//fmt.Print(s)
		os.Stdout.WriteString(s)
		frameSpeed := 1000000 / float64(time.Now().UnixMicro()-startMicroSec)
		os.Stdout.WriteString(fmt.Sprintf("Frame Speed: %.2f fps, Frame Counter: %d", frameSpeed, i))
		//time.Sleep(5 * time.Millisecond)
		inputFile.Close()
		//fmt.Print("\033[H\033[2J")
		os.Stdout.WriteString("\033[H")
	}
}
