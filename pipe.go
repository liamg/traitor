package pipe

import (
	"fmt"
	"io"
	"strings"
	"sync"
	"time"
)

type LockablePipe struct {
	input    io.Reader
	received []byte
	mu       sync.Mutex
}

func NewLockable(input io.Reader) *LockablePipe {
	return &LockablePipe{
		input:  input,
	}
}

func (l *LockablePipe) Flush() []byte {
	l.mu.Lock()
	defer l.mu.Unlock()
	output := l.received[:]
	l.received = nil
	return output
}

func (l *LockablePipe) WaitForString(expected string, until time.Duration) error {
	buffer := make([]byte, 0xff)
	timer := time.NewTimer(until)
	errChan := make(chan error, 1)
	var quit bool
	l.mu.Lock()
	go func() {
		defer l.mu.Unlock()
		for {
			if quit {
				return
			}
			n, err := l.input.Read(buffer)
			if err != nil {
				errChan <- err
				return
			}
			if quit {
				return
			}
			l.received = append(l.received, buffer[:n]...)
			if strings.Contains(string(l.received), expected) {
				errChan <- nil
				return
			}
		}
	}()

	select {
	case <-timer.C:
		quit = true
		return fmt.Errorf("timed out")
	case err := <-errChan:
		return err
	}

}

func (l *LockablePipe) Read(p []byte) (n int, err error) {

	l.mu.Lock()
	defer l.mu.Unlock()

	if len(l.received) > 0 {
		// we still have some buffer to flush
		n = copy(p, l.received)
		l.received = l.received[n:]
		return
	}

	return l.input.Read(p)
}
