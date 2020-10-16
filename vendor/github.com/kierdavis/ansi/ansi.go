package ansi

import (
	"fmt"
	"io"
	"sync"
)

func init() {
	go func() {
		// On program exit, ensure that we reset the screen.
		defer func() {
			fmt.Print("\x1b[0m")
		}()

		<-make(chan struct{})
	}()
}

var Mutex sync.Mutex
var UseMutex bool

type Color uint

type Attribute struct {
	Attr uint
	FG   Color
	BG   Color
}

const CSI = "\x1b["

const (
	Bold uint = 1 << iota
	Underline
	Blink
	Inverse
)

const (
	ColorNone    Color = 0 // This should be black, but the default (0) needs to be "no color"
	ColorRed     Color = 1
	ColorGreen   Color = 2
	ColorYellow  Color = 3
	ColorBlue    Color = 4
	ColorMagenta Color = 5
	ColorCyan    Color = 6
	ColorWhite   Color = 7
	ColorBlack   Color = 8 // Luckily, color % 8 will turn this 8 into a 0 correctly
)

var (
	Red     = Attribute{FG: ColorRed}
	Green   = Attribute{FG: ColorGreen}
	Yellow  = Attribute{FG: ColorYellow}
	Blue    = Attribute{FG: ColorBlue}
	Magenta = Attribute{FG: ColorMagenta}
	Cyan    = Attribute{FG: ColorCyan}
	White   = Attribute{FG: ColorWhite}
	Black   = Attribute{FG: ColorBlack}

	RedBold     = Attribute{FG: ColorRed, Attr: Bold}
	GreenBold   = Attribute{FG: ColorGreen, Attr: Bold}
	YellowBold  = Attribute{FG: ColorYellow, Attr: Bold}
	BlueBold    = Attribute{FG: ColorBlue, Attr: Bold}
	MagentaBold = Attribute{FG: ColorMagenta, Attr: Bold}
	CyanBold    = Attribute{FG: ColorCyan, Attr: Bold}
	WhiteBold   = Attribute{FG: ColorWhite, Attr: Bold}
	BlackBold   = Attribute{FG: ColorBlack, Attr: Bold}
)

func SAttrOn(attr Attribute) (s string) {
	if attr.FG != ColorNone {
		s += fmt.Sprintf("\x1b[3%dm", attr.FG%8)
	}

	if attr.BG != ColorNone {
		s += fmt.Sprintf("\x1b[4%dm", attr.BG%8)
	}

	if attr.Attr&Bold != 0 {
		s += "\x1b[1m"
	}

	if attr.Attr&Underline != 0 {
		s += "\x1b[4m"
	}

	if attr.Attr&Blink != 0 {
		s += "\x1b[5m"
	}

	if attr.Attr&Inverse != 0 {
		s += "\x1b[7m"
	}

	return s
}

func FAttrOn(w io.Writer, attr Attribute) (n int, err error) {
	if UseMutex {
		Mutex.Lock()
		defer Mutex.Unlock()
	}

	return fmt.Fprint(w, SAttrOn(attr))
}

func AttrOn(attr Attribute) (n int, err error) {
	if UseMutex {
		Mutex.Lock()
		defer Mutex.Unlock()
	}

	return fmt.Print(SAttrOn(attr))
}

func SAttrOff(attr Attribute) (s string) {
	if attr.FG != ColorNone {
		s += CSI + "39m"
	}

	if attr.BG != ColorNone {
		s += CSI + "49m"
	}

	if attr.Attr&Bold != 0 {
		s += CSI + "22m"
	}

	if attr.Attr&Underline != 0 {
		s += CSI + "24m"
	}

	if attr.Attr&Blink != 0 {
		s += CSI + "25m"
	}

	if attr.Attr&Inverse != 0 {
		s += CSI + "27m"
	}

	return s
}

func FAttrOff(w io.Writer, attr Attribute) (n int, err error) {
	if UseMutex {
		Mutex.Lock()
		defer Mutex.Unlock()
	}

	return fmt.Fprint(w, SAttrOff(attr))
}

func AttrOff(attr Attribute) (n int, err error) {
	if UseMutex {
		Mutex.Lock()
		defer Mutex.Unlock()
	}

	return fmt.Print(SAttrOff(attr))
}

func Sprint(attr Attribute, a ...interface{}) (s string) {
	return SAttrOn(attr) + fmt.Sprint(a...) + SAttrOff(attr)
}

func Sprintln(attr Attribute, a ...interface{}) (s string) {
	return SAttrOn(attr) + fmt.Sprintln(a...) + SAttrOff(attr)
}

func Sprintf(attr Attribute, format string, a ...interface{}) (s string) {
	return SAttrOn(attr) + fmt.Sprintf(format, a...) + SAttrOff(attr)
}

func Print(attr Attribute, a ...interface{}) (n int, err error) {
	if UseMutex {
		Mutex.Lock()
		defer Mutex.Unlock()
	}

	return fmt.Print(Sprint(attr, a...))
}

func Println(attr Attribute, a ...interface{}) (n int, err error) {
	if UseMutex {
		Mutex.Lock()
		defer Mutex.Unlock()
	}

	return fmt.Print(Sprintln(attr, a...))
}

func Printf(attr Attribute, format string, a ...interface{}) (n int, err error) {
	if UseMutex {
		Mutex.Lock()
		defer Mutex.Unlock()
	}

	return fmt.Print(Sprintf(attr, format, a...))
}

func Fprint(w io.Writer, attr Attribute, a ...interface{}) (n int, err error) {
	if UseMutex {
		Mutex.Lock()
		defer Mutex.Unlock()
	}

	return fmt.Fprint(w, Sprint(attr, a...))
}

func Fprintln(w io.Writer, attr Attribute, a ...interface{}) (n int, err error) {
	if UseMutex {
		Mutex.Lock()
		defer Mutex.Unlock()
	}

	return fmt.Fprint(w, Sprintln(attr, a...))
}

func Fprintf(w io.Writer, attr Attribute, format string, a ...interface{}) (n int, err error) {
	if UseMutex {
		Mutex.Lock()
		defer Mutex.Unlock()
	}

	return fmt.Fprint(w, Sprintf(attr, format, a...))
}

func CursorUp(lines int) {
	fmt.Printf("%s%dA", CSI, lines)
}

func CursorDown(lines int) {
	fmt.Printf("%s%dB", CSI, lines)
}

func CursorForward(cols int) {
	fmt.Printf("%s%dC", CSI, cols)
}

func CursorBack(cols int) {
	fmt.Printf("%s%dD", CSI, cols)
}

func CursorNextLine(lines int) {
	fmt.Printf("%s%dE", CSI, lines)
}

func CursorPrevLine(lines int) {
	fmt.Printf("%s%dF", CSI, lines)
}

func CursorHozPosition(column int) {
	fmt.Printf("%s%dG", CSI, column)
}

func CursorPosition(row int, column int) {
	fmt.Printf("%s%d;%dH", CSI, row, column)
}

func ClearToEndOfScreen() {
	fmt.Printf("%s0J", CSI)
}

func ClearToStartOfScreen() {
	fmt.Printf("%s1J", CSI)
}

func ClearScreen() {
	fmt.Printf("%s2J", CSI)
}

func ClearToEndOfLine() {
	fmt.Printf("%s0K", CSI)
}

func ClearToStartOfLine() {
	fmt.Printf("%s1K", CSI)
}

func ClearLine() {
	fmt.Printf("%s2K", CSI)
}

func ScrollUp(lines int) {
	fmt.Printf("%s%dS", CSI, lines)
}

func ScrollDown(lines int) {
	fmt.Printf("%s%dT", CSI, lines)
}

func SavePosition() {
	fmt.Printf("%ss", CSI)
}

func RestorePosition() {
	fmt.Printf("%su", CSI)
}

func HideCursor() {
	fmt.Printf("%s?25l", CSI)
}

func ShowCursor() {
	fmt.Printf("%s?25h", CSI)
}
