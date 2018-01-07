package main

import (
	"golang.org/x/crypto/ssh/terminal"
	"io"
	"math/rand"
	"os"
)

type Chooser struct {
	height  int
	width   int
	options []string
	term    *terminal.Terminal
	r       io.Reader
}

func NewChooser() *Chooser {
	width, _, _ := terminal.GetSize(0)
	term := terminal.NewTerminal(os.Stdin, "")

	return &Chooser{
		r:       os.Stdin,
		term:    term,
		height:  5,
		width:   width,
		options: []string{"one", "two", "three", "four", "five", "six", "seven"},
	}
}

var (
	tcEscape    = byte(27)
	tcUp        = []byte{tcEscape, '[', 'A'}
	tcDown      = []byte{tcEscape, '[', 'B'}
	tcClearLine = []byte{tcEscape, '[', '2', 'K'}
	tcLineStart = []byte{tcEscape, '[', 'G'}
	tcReturn    = []byte{'\n'}
)

func (c Chooser) Choose() {
	for n := 0; n < c.height-1; n++ {
		c.term.Write(tcReturn)
	}
	c.clear()
	c.draw()
	c.readInput()
}

func (c Chooser) writeln(ln string) {
	c.writeRaw(tcLineStart)
	c.write(ln)
	c.writeRaw(tcDown)
}

func (c Chooser) write(str string) {
	c.term.Write([]byte(str))
}

func (c Chooser) writeRaw(str []byte) {
	c.term.Write([]byte(str))
}

func (c Chooser) clear() {
	for n := 0; n < c.height-1; n++ {
		c.writeRaw(tcUp)
		c.writeRaw(tcClearLine)
	}
}

func (c Chooser) drawPrompt() {
	c.writeRaw(tcLineStart)
	c.write("> ")
}

func (c Chooser) draw() {
	n := 0
	for n = 0; n < len(c.options); n++ {
		if n > (c.height - 2) {
			break
		}
		c.writeln(c.options[n])
	}
	c.drawPrompt()
}

func (c *Chooser) sortOptions() {
	for n := range c.options {
		m := rand.Intn(n + 1)
		c.options[n], c.options[m] = c.options[m], c.options[n]
	}
}

func (c Chooser) readInput() {
	var result []byte
	var f = os.Stdin

	for {
		var buf [1]byte
		n, _ := f.Read(buf[:])

		if n == 0 || buf[0] == '\n' || buf[0] == '\r' {
			break
		}

		result = append(result, buf[0])
		c.clear()
		c.sortOptions()
		c.draw()
	}

	c.write("\n")
}

func main() {
	oldstate, _ := terminal.MakeRaw(0)
	defer terminal.Restore(0, oldstate)

	chooser := NewChooser()
	chooser.Choose()
}
