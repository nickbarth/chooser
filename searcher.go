package chooser

import (
	_ "fmt"
	"regexp"
	"sort"
	"strings"
)

// Fuzzy Search
// https://en.wikipedia.org/wiki/Approximate_string_matching
type FuzzySearcher struct{}

type sortables struct {
	value string
	score int
}

func (f FuzzySearcher) getMatches(search string, searchable []string) []sortables {
	var matches []sortables

	pattern := strings.Join(strings.Split(search, ""), ".*?")
	r, _ := regexp.Compile(pattern)

	for _, search := range searchable {
		if r.MatchString(search) {
			matches = append(matches, sortables{value: search, score: len(search)})
		}
	}

	return matches
}

func (f FuzzySearcher) sortMatches(matches []sortables) []string {
	var sorted []string

	sort.Slice(matches, func(n, m int) bool {
		return matches[n].score > matches[m].score
	})

	for _, s := range matches {
		sorted = append(sorted, s.value)
	}

	return sorted
}

func (f FuzzySearcher) Search(search string, searchable []string) []string {
	matches := f.getMatches(search, searchable)
	return f.sortMatches(matches)
}

/*
func main() {
	f := FuzzySearcher{}
	ss := f.Search("abc", []string{"ahello btown cthere", "abc", "aabbcc"})

	fmt.Println(ss)
}
*/
