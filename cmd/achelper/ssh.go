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

// ReadHostsFile reads the hosts file and returns a list of hosts
func ReadHostsFile() ([]string, error) {
	hostsLines, err := readFileAsList("/etc/hosts")
	if err != nil {
		return nil, err
	}
	return parseHostsFile(hostsLines), nil
}

// maxIncludeDepth bounds Include recursion, as ssh itself does.
const maxIncludeDepth = 16

// includePattern matches an "Include" keyword in an ssh config file.
var includePattern = regexp.MustCompile(`(?i)^include(?:\s+|\s*=\s*)(.+)$`)

// sshDir returns ~/.ssh.
func sshDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".ssh"), nil
}

// readSshConfigFile returns the hosts in fpath, following Include lines.
// Relative include paths are resolved against basedir (~/.ssh) and may
// contain globs. Unreadable included files are skipped, like ssh does.
func readSshConfigFile(fpath, basedir string, depth int) ([]string, error) {
	if depth > maxIncludeDepth {
		return nil, fmt.Errorf("%s: too many nested includes", fpath)
	}

	lines, err := readFileAsList(fpath)
	if err != nil {
		return nil, err
	}

	retv := parseSshConfig(lines)
	for _, line := range lines {
		matches := includePattern.FindStringSubmatch(line)
		if len(matches) != 2 {
			continue
		}
		for _, pattern := range strings.Fields(matches[1]) {
			if expanded, err := Expand(pattern); err == nil {
				pattern = expanded
			}
			if !filepath.IsAbs(pattern) {
				pattern = filepath.Join(basedir, pattern)
			}
			files, _ := filepath.Glob(pattern)
			for _, file := range files {
				if hosts, err := readSshConfigFile(file, basedir, depth+1); err == nil {
					retv = append(retv, hosts...)
				}
			}
		}
	}
	return retv, nil
}

// ReadSshConfig reads ~/.ssh/config (and its includes) and returns a list of hosts
func ReadSshConfig() ([]string, error) {
	dir, err := sshDir()
	if err != nil {
		return nil, err
	}
	fpath := filepath.Join(dir, "config")

	if _, err := os.Stat(fpath); os.IsNotExist(err) {
		return nil, fmt.Errorf("file %s does not exist", fpath)
	}

	return readSshConfigFile(fpath, dir, 0)
}

// parseKnownHosts returns the plain host names from known_hosts lines.
// Hashed entries cannot be recovered and are skipped; "[host]:port" is
// reduced to host.
func parseKnownHosts(lines []string) []string {
	retv := []string{}
	for _, line := range lines {
		columns := strings.Fields(line)
		if len(columns) > 0 && strings.HasPrefix(columns[0], "@") {
			columns = columns[1:] // @cert-authority / @revoked marker
		}
		if len(columns) < 2 || strings.HasPrefix(columns[0], "|") {
			continue
		}
		for _, host := range strings.Split(columns[0], ",") {
			if strings.HasPrefix(host, "[") {
				if end := strings.Index(host, "]"); end > 0 {
					host = host[1:end]
				}
			}
			if host != "" && !isHostPattern(host) {
				retv = append(retv, host)
			}
		}
	}
	return retv
}

// ReadKnownHosts reads ~/.ssh/known_hosts and returns a list of hosts
func ReadKnownHosts() ([]string, error) {
	dir, err := sshDir()
	if err != nil {
		return nil, err
	}
	lines, err := readFileAsList(filepath.Join(dir, "known_hosts"))
	if err != nil {
		return nil, err
	}
	return parseKnownHosts(lines), nil
}

// completionScript returns the bash complete commands for hosts. scp and
// rsync fall back to file name completion when no host matches.
func completionScript(hosts string) string {
	return fmt.Sprintf("complete -W \"%s\" ssh sftp\ncomplete -o default -W \"%s\" scp rsync\n", hosts, hosts)
}

// SecureShellHelper prints bash completion commands for ssh, sftp, scp and
// rsync with the hosts from the hosts file and ssh config, and optionally
// known_hosts.
func SecureShellHelper(knownHosts bool) error {
	terms := &Terms{terms: make(map[string]bool)}

	sources := []func() ([]string, error){ReadHostsFile, ReadSshConfig}
	if knownHosts {
		sources = append(sources, ReadKnownHosts)
	}
	for _, source := range sources {
		if hosts, err := source(); err == nil {
			for _, host := range hosts {
				terms.Add(host)
			}
		}
	}

	fmt.Print(completionScript(terms.String()))

	return nil
}
