package colors

import "testing"

func TestColornameToColorvalue(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{"black", ColorBlack},
		{"blue", ColorBlue},
		{"brown", ColorBrown},
		{"cyan", ColorCyan},
		{"darkgray", ColorDarkGray},
		{"gray", ColorGray},
		{"green", ColorGreen},
		{"lightblue", ColorLightBlue},
		{"lightcyan", ColorLightCyan},
		{"lightgray", ColorLightGray},
		{"lightgreen", ColorLightGreen},
		{"lightpurple", ColorLightPurpl},
		{"lightred", ColorLightRed},
		{"purple", ColorPurple},
		{"red", ColorRed},
		{"white", ColorWhite},
		{"yellow", ColorYellow},
		{"end", ColorEnd},
		{"unknown-color-name", ColorWhite},
		{"", ColorWhite},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			want := printc(tt.want)
			got := ColornameToColorvalue(tt.name)
			if got != want {
				t.Errorf("ColornameToColorvalue(%q) = %q, want %q", tt.name, got, want)
			}
		})
	}
}

func TestPrintc(t *testing.T) {
	got := printc(ColorRed)
	want := "\\[\033[0;31m\\]"
	if got != want {
		t.Errorf("printc(%q) = %q, want %q", ColorRed, got, want)
	}
}
