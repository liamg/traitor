package state

type State struct {
}

func New() *State {
	return &State{}
}

func (s *State) Assess() {
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
