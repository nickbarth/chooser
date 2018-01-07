package main

import (
	"fmt"
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

func NewChooser(height int, width int) *Chooser {
	// width, height, _ := terminal.GetSize(0)
	term := terminal.NewTerminal(os.Stdin, "")

	return &Chooser{
		r:       os.Stdin,
		search:  FuzzySearcher{},
		term:    term,
		height:  height,
		width:   width,
		choices: []string{},
		matches: []string{},
	}
}

const (
	tcEscape    = byte(27)
	tcNewline   = byte(10)
	tcReturn    = byte(13)
	tcBackspace = byte(127)
	tcCtrlC     = byte(23)
	tcCtrlW     = byte(23)
	tcCtrlU     = byte(21)
	tcTab       = byte(9)
)

var (
	tcUp        = []byte{tcEscape, '[', 'A'}
	tcDown      = []byte{tcEscape, '[', 'B'}
	tcClearLine = []byte{tcEscape, '[', '2', 'K'}
	tcLineStart = []byte{tcEscape, '[', 'G'}
)

func (c Chooser) Choose(choices []string) string {
	oldstate, _ := terminal.MakeRaw(0)

	c.choices = choices
	c.matches = choices

	for n := 0; n < c.height-1; n++ {
		c.term.Write([]byte{tcNewline})
	}

	c.clear()
	c.printChoices()
	c.printPrompt()
	c.readInput()

	terminal.Restore(0, oldstate)

	return "one"
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
	c.writeRaw(tcClearLine)
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

		if n == 0 || buf[0] == tcReturn {
			break
		}

		switch buf[0] {
		case tcTab:
		case tcBackspace:
			if len(search) > 1 {
				search = search[:len(search)-1]
			} else if len(search) == 1 {
				search = []byte{}
			}
		case tcCtrlW, tcCtrlU:
			search = []byte{}
		default:
			search = append(search, buf[0])
		}

		c.clear()
		c.searchOptions(string(search))
		c.printChoices()
		c.printPrompt()
		c.write(string(search))
	}

	c.write("\n")
}

func main() {
	chooser := NewChooser(5, 20)
	choice := chooser.Choose([]string{"one", "two", "three", "four", "five", "six", "seven"})

	fmt.Println("You Chose:", choice)
}
