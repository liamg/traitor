package state

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
)

func (s *State) processSudoers(hostname string) {
	args := []string{"-l"}
	if !s.HasPassword {
		args = append(args, "-n")
	}
	cmd := exec.Command("sudo", args...)
	cmd.Env = append(os.Environ(), "LANG=C", "LC_MESSAGES=C")
	data, err := cmd.Output()
	if err != nil {
		return
	}

	lines := strings.Split(string(data), "\n")
	commandSection := false
	for _, line := range lines {

		if commandSection {
			item, err := parseSudoLine(line, hostname)
			if err != nil {
				continue
			}
			s.SudoEntries = append(s.SudoEntries, item)
		}

		if strings.Contains(line, "may run the following") {
			commandSection = true
			continue
		}
	}
}

type SudoEntry struct {
	AllUsers        bool
	UserName        string
	AllHosts        bool
	Hostname        string
	AllCommands     bool
	Command         string
	NoPasswd        bool
	BinaryName      string
	HostnameMatches bool
}

type Sudoers []*SudoEntry

func (s Sudoers) GetEntryForBinary(binary string, hasPasswd bool) (*SudoEntry, error) {
	for _, entry := range s {
		if (entry.BinaryName == binary || entry.AllCommands) && entry.HostnameMatches {
			if (entry.UserName == "root" || entry.AllUsers) && (hasPasswd || entry.NoPasswd) {
				return entry, nil
			}
		}
	}
	return nil, fmt.Errorf("nothing found")
}

func parseSudoLine(line string, hostname string) (*SudoEntry, error) {
	line = strings.TrimSpace(line)
	if line == "" {
		return nil, fmt.Errorf("empty line")
	}
	if line[0] != '(' {
		return nil, fmt.Errorf("invalid line: no bracketed section")
	}
	line = line[1:]
	parts := strings.SplitN(line, ")", 2)
	bracketed := parts[0]

	if len(parts) == 1 {
		return nil, fmt.Errorf("invalid line: no command")
	}

	var entry SudoEntry

	bracketedParts := strings.SplitN(bracketed, ":", 2)
	entry.UserName = strings.TrimSpace(bracketedParts[0])
	entry.AllUsers = entry.UserName == "ALL"
	entry.AllHosts = true
	if len(bracketedParts) == 2 {
		entry.Hostname = strings.TrimSpace(bracketedParts[1])
		entry.AllHosts = entry.Hostname == "ALL"
	}

	specification := parts[1]
	parts = strings.Split(specification, ":")

	for i := 0; i < len(parts)-1; i++ {
		switch strings.TrimSpace(parts[i]) {
		case "NOPASSWD":
			entry.NoPasswd = true
		}
	}

	entry.Command = strings.TrimSpace(parts[len(parts)-1])
	entry.AllCommands = entry.Command == "ALL"

	if !entry.AllCommands {
		commandParts := strings.Split(entry.Command, " ")
		entry.BinaryName = path.Base(commandParts[0])
	}

	entry.HostnameMatches = entry.AllHosts || entry.Hostname == hostname

	return &entry, nil
}
