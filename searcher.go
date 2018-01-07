package main

import (
	"github.com/renstrom/fuzzysearch/fuzzy"
	"sort"
)

type FuzzySearcher struct{}

func (s FuzzySearcher) Search(search string, sortable []string) (sorted []string) {
	matches := fuzzy.RankFind(search, sortable)
	sort.Sort(matches)
	sort.Sort(sort.Reverse(matches))

	sorted = make([]string, len(matches))
	for n, match := range matches {
		sorted[n] = match.Target
	}

	return sorted
}
