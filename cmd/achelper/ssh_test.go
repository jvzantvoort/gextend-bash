package main

import (
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"testing"
)

func TestParseSshConfig(t *testing.T) {
	// lines as returned by readFileAsList: trimmed, comments removed
	lines := []string{
		"Host web1",
		"HostName web1.example.com",
		"host db1 db2",
		"Host=jump",
		"Host *",
		"Host *.example.com !bastion",
		"Match host foo",
		"User admin",
	}

	got := parseSshConfig(lines)
	want := []string{"web1", "db1", "db2", "jump"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("parseSshConfig() = %v, want %v", got, want)
	}
}

func TestParseHostsFile(t *testing.T) {
	lines := []string{
		"127.0.0.1\tlocalhost localhost.localdomain",
		"::1             ip6-localhost ip6-loopback",
		"192.168.1.10 server",
		"10.0.0.1",
	}

	got := parseHostsFile(lines)
	want := []string{"localhost", "localhost.localdomain", "ip6-localhost", "ip6-loopback", "server"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("parseHostsFile() = %v, want %v", got, want)
	}
}

func TestReadSshConfigFileIncludes(t *testing.T) {
	dir := t.TempDir()
	write := func(name, content string) {
		t.Helper()
		if err := os.MkdirAll(filepath.Dir(filepath.Join(dir, name)), 0755); err != nil {
			t.Fatalf("os.MkdirAll() error = %v", err)
		}
		if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0644); err != nil {
			t.Fatalf("os.WriteFile() error = %v", err)
		}
	}
	write("config", "Host main\nInclude conf.d/*.conf missing.conf\n")
	write("conf.d/a.conf", "Host alpha\ninclude = ../nested\n")
	write("conf.d/b.conf", "Host beta\n")
	write("nested", "Host gamma\n")
	write("loop", "Host loop\nInclude loop\n")

	got, err := readSshConfigFile(filepath.Join(dir, "config"), dir, 0)
	if err != nil {
		t.Fatalf("readSshConfigFile() error = %v", err)
	}
	sort.Strings(got)
	want := []string{"alpha", "beta", "main"}
	if !reflect.DeepEqual(got, want) {
		// "../nested" resolves relative to the ssh dir, not the including file
		t.Errorf("readSshConfigFile() = %v, want %v", got, want)
	}

	t.Run("include loops are bounded", func(t *testing.T) {
		got, err := readSshConfigFile(filepath.Join(dir, "loop"), dir, 0)
		if err != nil {
			t.Fatalf("readSshConfigFile() error = %v", err)
		}
		if len(got) == 0 || got[0] != "loop" {
			t.Errorf("readSshConfigFile() = %v, want it to start with loop", got)
		}
	})
}

func TestParseKnownHosts(t *testing.T) {
	lines := []string{
		"web1,192.168.1.10 ssh-ed25519 AAAA",
		"[git.example.com]:2222 ssh-rsa AAAA",
		"|1|abc=|def= ssh-ed25519 AAAA",
		"@cert-authority *.example.com ssh-rsa AAAA",
		"@revoked old.example.com ssh-rsa AAAA",
		"broken-line",
	}

	got := parseKnownHosts(lines)
	want := []string{"web1", "192.168.1.10", "git.example.com", "old.example.com"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("parseKnownHosts() = %v, want %v", got, want)
	}
}

func TestCompletionScript(t *testing.T) {
	got := completionScript("a b")
	want := "complete -W \"a b\" ssh sftp\ncomplete -o default -W \"a b\" scp rsync\n"
	if got != want {
		t.Errorf("completionScript() = %q, want %q", got, want)
	}
}
