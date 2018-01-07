package main

import (
	_ "fmt"
	"github.com/renstrom/fuzzysearch/fuzzy"
	"golang.org/x/crypto/ssh/terminal"
	"os"
	"path/filepath"
	"sort"
)

func enter(term *terminal.Terminal) {
	// term.Write([]byte{'\n'})         // ENTER
	term.Write([]byte{27, '[', 'B'}) // DOWN
	term.Write([]byte{27, '[', 'G'}) // START OF LINE
}

func clear(term *terminal.Terminal) {
	term.Write([]byte{'\n'})
	term.Write([]byte{'\n'})
	term.Write([]byte{'\n'})
	term.Write([]byte{'\n'})
	term.Write([]byte{'\n'})
}

func out(term *terminal.Terminal) {
	term.Write([]byte{27, '[', 'B'}) // DOWN
	term.Write([]byte{27, '[', 'B'}) // DOWN
	term.Write([]byte{27, '[', 'B'}) // DOWN
	term.Write([]byte{27, '[', 'B'}) // DOWN
	term.Write([]byte{27, '[', 'B'}) // DOWN
	term.Write([]byte{27, '[', 'G'}) // START OF LINE
}

func top(term *terminal.Terminal) {
	term.Write([]byte{27, '[', 'A'})      // UP
	term.Write([]byte{27, '[', '2', 'K'}) // CLEAR LINE
	term.Write([]byte{27, '[', 'A'})      // UP
	term.Write([]byte{27, '[', '2', 'K'}) // CLEAR LINE
	term.Write([]byte{27, '[', 'A'})      // UP
	term.Write([]byte{27, '[', '2', 'K'}) // CLEAR LINE
	term.Write([]byte{27, '[', 'A'})      // UP
	term.Write([]byte{27, '[', '2', 'K'}) // CLEAR LINE
	term.Write([]byte{27, '[', 'A'})      // UP
	term.Write([]byte{27, '[', '2', 'K'}) // CLEAR LINE
	term.Write([]byte{27, '[', 'G'})      // START OF LINE
}

func bottom(term *terminal.Terminal) {
	term.Write([]byte{27, '[', 'B'}) // UP
	term.Write([]byte{27, '[', 'B'}) // UP
	term.Write([]byte{27, '[', 'B'}) // UP
	term.Write([]byte{27, '[', 'B'}) // UP
	term.Write([]byte{27, '[', 'B'}) // UP
	term.Write([]byte{27, '[', 'G'}) // START OF LINE
}

func printOptions(term *terminal.Terminal, find string, files []string) {
	matches := fuzzy.RankFind(find, files)
	sort.Sort(matches)
	sort.Sort(sort.Reverse(matches))

	for i := 0; i < 5-len(matches); i++ {
		enter(term)
	}

	for n, match := range matches {
		if n > 5 {
			break
		}
		term.Write([]byte(match.Target))
		enter(term)
	}

	term.Write([]byte{27, '[', '2', 'K'}) // CLEAR LINE
	term.Write([]byte("> "))              // ENTER
}

func main() {
	// grab files
	files, err := filepath.Glob("*")
	if err != nil {
		panic(err)
	}

	// create terminal
	oldstate, _ := terminal.MakeRaw(0)
	// width, height, _ := terminal.GetSize(0)
	defer terminal.Restore(0, oldstate)

	term := terminal.NewTerminal(os.Stdin, "")
	// enter(term)
	clear(term)
	top(term)
	printOptions(term, "", files)

	var resultBuf []byte
	var f = os.Stdin
	for {
		var buf [1]byte

		n, _ := f.Read(buf[:])
		// if err != nil && err != io.EOF { return "", err }

		if n == 0 || buf[0] == '\n' || buf[0] == '\r' {
			break
		}

		if buf[0] == byte(127) {
			if len(resultBuf) > 1 {
				resultBuf = resultBuf[:len(resultBuf)-1]
			} else if len(resultBuf) == 1 {
				resultBuf = []byte{}
			}
		} else if buf[0] == byte(21) || buf[0] == byte(11) || buf[0] == byte(23) || buf[0] == byte(12) {
			resultBuf = []byte{}
		} else {
			resultBuf = append(resultBuf, buf[0])
		}

		top(term)
		printOptions(term, string(resultBuf), files)
		term.Write([]byte(resultBuf))
	}

	term.Write([]byte("\n"))
	enter(term)
}
