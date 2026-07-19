package main

import (
	"strings"
	"testing"
)

func TestIndentStr(t *testing.T) {
	tests := []struct {
		num  int
		want string
	}{
		{0, ""},
		{1, "  "},
		{3, "      "},
	}

	for _, tt := range tests {
		got := IndentStr(tt.num)
		if got != tt.want {
			t.Errorf("IndentStr(%d) = %q, want %q", tt.num, got, tt.want)
		}
	}
}

func TestStartEndFunc(t *testing.T) {
	before := indent
	defer func() { indent = before }()

	startMsg := StartFunc("myfunc")
	if !strings.Contains(startMsg, "myfunc") || !strings.HasSuffix(startMsg, "START") {
		t.Errorf("StartFunc() = %q, want it to contain the function name and end with START", startMsg)
	}
	if indent != before+1 {
		t.Errorf("indent after StartFunc() = %d, want %d", indent, before+1)
	}

	endMsg := EndFunc("myfunc")
	if !strings.Contains(endMsg, "myfunc") || !strings.HasSuffix(endMsg, "END") {
		t.Errorf("EndFunc() = %q, want it to contain the function name and end with END", endMsg)
	}
	if indent != before {
		t.Errorf("indent after EndFunc() = %d, want %d", indent, before)
	}
}

func TestStatMsg(t *testing.T) {
	t.Run("short message is preserved", func(t *testing.T) {
		got := StatMsg("hello", "SUCCESS")
		if !strings.Contains(got, "hello") {
			t.Errorf("StatMsg() = %q, want it to contain %q", got, "hello")
		}
		if !strings.Contains(got, "SUCCESS") {
			t.Errorf("StatMsg() = %q, want it to contain %q", got, "SUCCESS")
		}
	})

	t.Run("long message is truncated", func(t *testing.T) {
		long := strings.Repeat("x", 100)
		got := StatMsg(long, "FAILURE")
		if strings.Contains(got, long) {
			t.Error("expected the message to be truncated")
		}
		if !strings.Contains(got, "...") {
			t.Errorf("StatMsg() = %q, want it to contain an ellipsis", got)
		}
	})
}

func TestCurFunc(t *testing.T) {
	got := CurFunc()
	if !strings.HasSuffix(got, "TestCurFunc") {
		t.Errorf("CurFunc() = %q, want it to end with %q", got, "TestCurFunc")
	}
	if strings.Contains(got, "main.") {
		t.Errorf("CurFunc() = %q, want the main. prefix to be stripped", got)
	}
}
