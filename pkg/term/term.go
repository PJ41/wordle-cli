package term

import (
    "fmt"
	"golang.org/x/term"
	"os"
    "bufio"
)

type Color int

const (
    DEFAULT Color = iota
    RED
    GREEN
    YELLOW
)

var colorMap = map[Color]string {
    DEFAULT : "\033[0m",
    RED : "\033[31m",
    GREEN : "\033[32m",
    YELLOW : "\033[33m",
}
var writer *bufio.Writer
var oldState *term.State

func Configure() (err error) {
    fd := int(os.Stdin.Fd())
	oldState, err = term.MakeRaw(fd)

    if err != nil {
        return err
    }

    writer = bufio.NewWriter(os.Stdout)

    return err
}

func Restore() {
    fd := int(os.Stdin.Fd())
	term.Restore(fd, oldState)
}

func Println(format string, a ...interface{}) {
    Print(format + "\n\r", a...)
}

func PrintColored(color Color, format string, a ...interface{}) {
    Print(colorMap[color])
    Print(format, a...)
    Print(colorMap[DEFAULT])
}

func Print(format string, a ...interface{}) {
    fmt.Fprintf(writer, format, a...)
}

func Flush() {
    writer.Flush()
}

func MoveCursorUp(rows int) {
    Print("\033[%dA", rows)
}

func MoveCursorDown(rows int) {
    Print("\033[%dB", rows)
}

func ClearRows(rows int) {
    MoveCursorUp(rows)
    for i := 0; i < rows; i++ {
        Println("\033[K")
    }
    MoveCursorUp(rows)
}
