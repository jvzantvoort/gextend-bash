package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// hostsPattern matches a "Host" keyword in an ssh config file; keywords are
// case-insensitive and may be separated from the value by '=' (see ssh_config(5)).
var hostsPattern = regexp.MustCompile(`(?i)^host(?:\s+|\s*=\s*)(.+)$`)

// isHostPattern reports whether an ssh config host entry is a pattern
// (wildcard or negation) rather than a name that can be completed.
func isHostPattern(host string) bool {
	return strings.ContainsAny(host, "*?!")
}

// parseSshConfig returns the host names from the Host lines of an ssh config.
func parseSshConfig(lines []string) []string {
	retv := []string{}
	for _, line := range lines {
		matches := hostsPattern.FindStringSubmatch(line)
		if len(matches) != 2 {
			continue
		}
		for _, host := range strings.Fields(matches[1]) {
			if !isHostPattern(host) {
				retv = append(retv, host)
			}
		}
	}
	return retv
}

// parseHostsFile returns the host names (all columns after the address)
// from the lines of a hosts file.
func parseHostsFile(lines []string) []string {
	retv := []string{}
	for _, line := range lines {
		columns := strings.Fields(line)
		if len(columns) < 2 {
			continue
		}
		retv = append(retv, columns[1:]...)
	}
	return retv
}

// ReadSshConfig reads the ssh config file and returns a list of hosts
func ReadSshConfig() ([]string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	fpath := filepath.Join(home, ".ssh", "config")

	if _, err := os.Stat(fpath); os.IsNotExist(err) {
		return nil, fmt.Errorf("file %s does not exist", fpath)
	}

	sshConfigLines, err := readFileAsList(fpath)
	if err != nil {
		return nil, err
	}

	return parseSshConfig(sshConfigLines), nil
}

// ReadHostsFile reads the hosts file and returns a list of hosts
func ReadHostsFile() ([]string, error) {
	hostsLines, err := readFileAsList("/etc/hosts")
	if err != nil {
		return nil, err
	}
	return parseHostsFile(hostsLines), nil
}

// SecureShellHelper prints the hosts in the hosts file and ssh config
func SecureShellHelper() error {
	terms := &Terms{terms: make(map[string]bool)}

	if hosts, err := ReadHostsFile(); err == nil {
		for _, host := range hosts {
			terms.Add(host)
		}
	}

	if hosts, err := ReadSshConfig(); err == nil {
		for _, host := range hosts {
			terms.Add(host)
		}
	}

	fmt.Printf("complete -W \"%s\" ssh\n", terms.String())

	return nil

}
