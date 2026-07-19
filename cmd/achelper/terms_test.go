package main

import "testing"

func TestTermsAddListString(t *testing.T) {
	terms := &Terms{terms: make(map[string]bool)}

	terms.Add("beta")
	terms.Add("alpha")
	terms.Add("gamma")
	terms.Add("alpha") // duplicate, should not appear twice

	got := terms.List()
	want := []string{"alpha", "beta", "gamma"}
	if len(got) != len(want) {
		t.Fatalf("List() = %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("List()[%d] = %q, want %q", i, got[i], want[i])
		}
	}

	if gotStr, want := terms.String(), "alpha beta gamma"; gotStr != want {
		t.Errorf("String() = %q, want %q", gotStr, want)
	}
}

func TestTermsStringEmpty(t *testing.T) {
	terms := &Terms{terms: make(map[string]bool)}
	if got := terms.String(); got != "" {
		t.Errorf("String() = %q, want empty string", got)
	}
}
