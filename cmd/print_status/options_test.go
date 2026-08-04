package main

import "testing"

func TestOptionsRegisteredAsSubcommands(t *testing.T) {
	if len(Options) == 0 {
		t.Fatal("expected at least one option")
	}

	for _, name := range Options {
		t.Run(name, func(t *testing.T) {
			cmd, _, err := rootCmd.Find([]string{name})
			if err != nil {
				t.Fatalf("rootCmd.Find(%q) error = %v", name, err)
			}
			if cmd.Use != name {
				t.Errorf("cmd.Use = %q, want %q", cmd.Use, name)
			}
			if cmd.Run == nil {
				t.Errorf("expected subcommand %q to have a Run function", name)
			}
		})
	}
}
