package main

import (
	"golang.org/x/crypto/ssh/terminal"
	"io"
	"os"
)

type Searcher interface {
	Search(search string, sortable []string) (sorted []string)
}

type Chooser struct {
	height  int
	width   int
	choices []string
	matches []string
	term    *terminal.Terminal
	search  Searcher
	r       io.Reader
}

func NewChooser(height int, choices []string) *Chooser {
	width, _, _ := terminal.GetSize(0)
	term := terminal.NewTerminal(os.Stdin, "")

	return &Chooser{
		r:       os.Stdin,
		search:  FuzzySearcher{},
		term:    term,
		height:  height,
		width:   width,
		choices: choices,
		matches: choices,
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
	c.printChoices()
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

func (c Chooser) printPrompt() {
	c.writeRaw(tcLineStart)
	c.write("> ")
}

func (c Chooser) printChoices() {
	for i := 0; i < c.height-len(c.matches)-1; i++ {
		c.writeln("")
	}

	for n, match := range c.matches {
		if n > c.height-2 {
			break
		}
		c.writeln(match)
	}

	c.printPrompt()
}

func (c *Chooser) searchOptions(search string) {
	c.matches = c.search.Search(search, c.choices)
}

func (c Chooser) readInput() {
	var search []byte
	var f = os.Stdin

	for {
		var buf [1]byte
		n, _ := f.Read(buf[:])

		if n == 0 || buf[0] == '\n' || buf[0] == '\r' {
			break
		}

		search = append(search, buf[0])
		c.clear()
		c.searchOptions(string(search))
		c.printChoices()
	}

	c.write("\n")
}

func main() {
	oldstate, _ := terminal.MakeRaw(0)
	defer terminal.Restore(0, oldstate)

	chooser := NewChooser(3, []string{"one", "two", "three", "four", "five", "six", "seven"})
	chooser.Choose()
}
