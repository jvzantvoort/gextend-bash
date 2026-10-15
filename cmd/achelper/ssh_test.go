package main

import (
	"reflect"
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
