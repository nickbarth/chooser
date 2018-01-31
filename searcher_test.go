package chooser

import (
	"reflect"
	"testing"
)

var searchTests = []struct {
	search   string
	terms    []string
	expected []string
}{
	{"abc", []string{"abc"}, []string{"abc"}},
	{"abc", []string{"abc", "aabbcc", "ahello btown cthere"}, []string{"ahello btown cthere", "aabbcc", "abc"}},
	{"abc", []string{"123", "456", "789"}, []string{}},
}

func TestSearch(t *testing.T) {
	f := FuzzySearcher{}

	for _, tt := range searchTests {
		actual := f.Search(tt.search, tt.terms)
		if !reflect.DeepEqual(actual, tt.expected) && !(len(actual) == 0 && len(tt.expected) == 0) {
			t.Errorf(".Search(%v): expected %v, actual %v", tt.search, tt.expected, actual)
		}
	}
}
