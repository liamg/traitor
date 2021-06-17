package state

import (
	"os/exec"
	"strings"
)

type DistributionID string

const (
	UnknownLinux DistributionID = "linux" // default
	Ubuntu       DistributionID = "ubuntu"
	Debian       DistributionID = "debian"
	Arch         DistributionID = "arch"
	RHEL         DistributionID = "rhel"
	Fedora       DistributionID = "fedora"
	CentOS       DistributionID = "centos"
	Kali         DistributionID = "kali"
	Parrot       DistributionID = "parrot"
	Alpine       DistributionID = "alpine"
	OpenSUSE     DistributionID = "opensuse"
)

func (s *State) processDistro() {
	s.DistroID = UnknownLinux
	data, err := exec.Command("sh", "-c", "cat /etc/*-release").Output()
	if err != nil {
		return
	}
	for _, line := range strings.Split(string(data), "\n") {
		if strings.HasPrefix(line, "ID=") {
			s.DistroID = DistributionID(strings.TrimSpace(line[3:]))
		}
		if strings.HasPrefix(line, "VERSION=") {
			s.DistroVersion = strings.TrimSpace(line[8:])
		}
	}
}

func (s *State) IsDebianLike() bool {
	switch s.DistroID {
	case Debian, Ubuntu, Kali, Parrot:
		return true
	}

	return false
}
