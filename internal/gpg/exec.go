// Package gpg is a store.GPGBackend implementation using the GPG binary.
package gpg

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

type gpg struct {
	path string
	env  []string
}

func New(path string) *gpg {
	return &gpg{
		path: path,
	}
}

func (gpg *gpg) cmd(args ...string) *exec.Cmd {
	cmd := exec.Command(gpg.path, args...)
	cmd.Env = append(os.Environ(), gpg.env...)
	return cmd
}

func (gpg *gpg) Encrypt(content []byte, recipients []string) ([]byte, error) {
	gpgArgs := []string{
		"--encrypt",
		"--batch",
		"--use-agent",
		"--no-tty",
		"--quiet",
		"--yes",
		"--output", "-",
	}

	for _, recipient := range recipients {
		gpgArgs = append(gpgArgs, "--recipient", recipient)
	}

	var stdout bytes.Buffer

	cmd := gpg.cmd(gpgArgs...)
	cmd.Stdout = &stdout

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("could not get gpg's stdin: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("could not start gpg: %w", err)
	}

	if _, err := stdin.Write(content); err != nil {
		return nil, fmt.Errorf("could not write to gpg's stdin: %w", err)
	}

	if err := stdin.Close(); err != nil {
		return nil, fmt.Errorf("could not close gpg's stdin: %w", err)
	}

	if err := cmd.Wait(); err != nil {
		return nil, fmt.Errorf("%w: %s", err, stdout.Bytes())
	}

	return stdout.Bytes(), nil
}

func (gpg *gpg) Decrypt(content []byte) ([]byte, error) {
	gpgArgs := []string{
		"--decrypt",
		"--batch",
		"--use-agent",
		"--no-tty",
		"--quiet",
		"--output", "-",
	}

	var stdout bytes.Buffer

	cmd := gpg.cmd(gpgArgs...)
	cmd.Stdout = &stdout

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("could not get gpg's stdin: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("could not start gpg: %w", err)
	}

	if _, err := stdin.Write(content); err != nil {
		return nil, fmt.Errorf("could not write to gpg's stdin: %w", err)
	}

	if err := stdin.Close(); err != nil {
		return nil, fmt.Errorf("could not close gpg's stdin: %w", err)
	}

	if err := cmd.Wait(); err != nil {
		return nil, fmt.Errorf("%w: %s", err, stdout.Bytes())
	}

	return stdout.Bytes(), err
}
