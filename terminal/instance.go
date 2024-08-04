package terminal

import (
	"errors"
	"os"
	"runtime"

	"github.com/mattn/go-isatty"
	"golang.org/x/term"
)

const charWidth = 0.5

type TerminalIF interface {
	GetCharWidth() float64
	GetScreenSize() (width, height int, err error)
	IsWindows() bool
}

type TerminalInstance struct {
	width  int
	height int
}

func (instance TerminalInstance) GetCharWidth() float64 {
	return charWidth
}

func (instance TerminalInstance) IsWindows() bool {
	return runtime.GOOS == "windows"
}

func (instance TerminalInstance) GetScreenSize() (width, height int, err error) {
	if instance.width != -1 && instance.height != -1 {
		return instance.width, instance.height, err
	}
	if !isatty.IsCygwinTerminal(os.Stdout.Fd()) && !isatty.IsTerminal(os.Stdout.Fd()) {
		return -1, -1, errors.New("this is not vaild terminal")
	}
	w, h, err := term.GetSize(int(os.Stdout.Fd()))
	instance.width = w
	instance.height = h
	return w, h, err
}

func NewTerminal() TerminalIF {
	return TerminalInstance{
		width:  -1,
		height: -1,
	}
}
