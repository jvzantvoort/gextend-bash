package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

// ReadSshConfig reads the ssh config file and returns a list of hosts
func ReadSshConfig() ([]string, error) {
	retv := []string{}                                 // hosts
	hostsPattern := regexp.MustCompile(`Host\s+(\S+)`) // ssh config host pattern
	fpath, err := Expand("$HOME/.ssh/config")          // ssh config file

	if err != nil {
		return retv, err
	}

	if _, err := os.Stat(fpath); os.IsNotExist(err) {
		return retv, fmt.Errorf("file %s does not exist", fpath)
	}

	sshConfigLines, err := readFileAsList(fpath)
	if err != nil {
		return retv, err
	}

	for _, line := range sshConfigLines {
		matches := hostsPattern.FindStringSubmatch(line)
		if len(matches) == 2 {
			retv = append(retv, matches[1])
		}
	}

	return retv, nil
}

// ReadHostsFile reads the hosts file and returns a list of hosts
func ReadHostsFile() ([]string, error) {
	retv := []string{}
	// open /etc/hosts and get hosts
	hostsLines, err := readFileAsList("/etc/hosts")
	if err != nil {
		return retv, err
	}

	for _, line := range hostsLines {
		columns := strings.Split(line, " ")
		columns = columns[:len(columns)-1]
		retv = append(retv, columns...)
	}
	return retv, nil
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
