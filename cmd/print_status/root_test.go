package main

import "testing"

func TestRootCmdUse(t *testing.T) {
	if rootCmd.Use != "print_status" {
		t.Errorf("rootCmd.Use = %q, want %q", rootCmd.Use, "print_status")
	}
}

func TestRootCmdHasVerboseFlag(t *testing.T) {
	flag := rootCmd.PersistentFlags().Lookup("verbose")
	if flag == nil {
		t.Fatal("expected a persistent --verbose flag")
	}
	if flag.Shorthand != "v" {
		t.Errorf("verbose flag shorthand = %q, want %q", flag.Shorthand, "v")
	}
}
