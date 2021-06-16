package state

import "os"

type State struct {
	HasPassword   bool
	SudoEntries   Sudoers
	Hostname      string
	DistroID      DistributionID
	DistroVersion string
}

func New() *State {
	return &State{}
}

func (s *State) Assess() {

	s.Hostname, _ = os.Hostname()

	s.processSudoers(s.Hostname)

	s.processDistro()

	// check existing backdoors
	// list users
	// list current user + groups
	// sudo -l
	// os info:
	//   (cat /proc/version || uname -a ) 2>/dev/null
	//   lsb_release -a 2>/dev/null
	// env vars
	// disks
	// printers
	// app armor/selinux etc.
	// installed packages
	// processes
	// cron
	// services
	// network
}
