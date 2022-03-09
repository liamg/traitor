package shell

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/creack/pty"
	"github.com/liamg/traitor/internal/pipe"
	"github.com/liamg/traitor/pkg/logger"
	"golang.org/x/crypto/ssh/terminal"
)

func WithPassword(username, password string, log logger.Logger) error {
	log.Printf("Setting up tty...")

	cmd := exec.Command("sh", "-c", fmt.Sprintf("su - %s", username))

	// Start the command with a pty.
	ptmx, err := pty.Start(cmd)
	if err != nil {
		return err
	}
	// Make sure to close the pty at the end.
	defer func() { _ = ptmx.Close() }() // Best effort.

	// Handle pty size.
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGWINCH)
	go func() {
		for range ch {
			_ = pty.InheritSize(os.Stdin, ptmx)
		}
	}()
	ch <- syscall.SIGWINCH // Initial resize.

	// Set stdin in raw mode.
	oldState, err := terminal.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return err
	}
	defer func() { _ = terminal.Restore(int(os.Stdin.Fd()), oldState) }() // Best effort.

	lockable := pipe.NewLockable(ptmx)

	log.Printf("Attempting authentication as %s...", username)
	if err := lockable.WaitForString("Password:", time.Second*2); err == nil {
		_ = lockable.Flush()
		if _, err := ptmx.Write([]byte(fmt.Sprintf("%s\n", password))); err != nil {
			return err
		}
		if err := lockable.WaitForString("#", time.Second*8); err != nil {
			_ = lockable.Flush()
			return fmt.Errorf("invalid password")
		}
		time.Sleep(time.Millisecond * 100)
	} else {
		return err
	}
	log.Printf("Authenticated as %s!", username)

	// Copy stdin to the pty and the pty to stdout.
	go func() { _, _ = io.Copy(ptmx, os.Stdin) }()
	_, _ = io.Copy(os.Stdout, lockable)
	return nil
}
