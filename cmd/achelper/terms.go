package main

import (
	"sort"
	"strings"
)

type Terms struct {
	terms map[string]bool
}

func (t *Terms) Add(term string) {
	t.terms[term] = true
}

func (t *Terms) List() []string {
	var terms []string
	for term := range t.terms {
		terms = append(terms, term)
	}
	sort.Strings(terms)
	return terms
}

func (t Terms) String() string {
	terms := t.List()
	return strings.TrimSpace(strings.Join(terms, " "))
}
