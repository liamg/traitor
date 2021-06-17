package state

import "os/exec"

func (s *State) IsPackageInstalled(name string) bool {

	switch s.DistroID {
	case Debian, Ubuntu, Kali, Parrot:
		return exec.Command("dpkg", "-s", name).Run() == nil
	case Arch:
		return exec.Command("pacman", "-Qi", name).Run() == nil
	case Fedora, RHEL, CentOS, OpenSUSE:
		return exec.Command("rpm", "-q", name).Run() == nil
	case Alpine:
		return exec.Command("apk", "-e", "info", name).Run() == nil
	default:
		return false
	}

}
