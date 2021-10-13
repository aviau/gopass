//    Copyright (C) 2018-2021 Alexandre Viau <alexandre@alexandreviau.net>
//
//    This file is part of gopass.
//
//    gopass is free software: you can redistribute it and/or modify
//    it under the terms of the GNU General Public License as published by
//    the Free Software Foundation, either version 3 of the License, or
//    (at your option) any later version.
//
//    gopass is distributed in the hope that it will be useful,
//    but WITHOUT ANY WARRANTY; without even the implied warranty of
//    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//    GNU General Public License for more details.
//
//    You should have received a copy of the GNU General Public License
//    along with gopass.  If not, see <http://www.gnu.org/licenses/>.

package cli

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"time"

	"github.com/aviau/gopass/pkg/store"
)

// CommandConfig contains everything that commands needs to run.
type CommandConfig interface {
	WriterOutput() io.Writer
	WriterError() io.Writer
	ReaderInput() io.Reader
	Edit(string) (string, error)
	PasswordStoreDir() string
	PasswordStore() *store.PasswordStore
	Now() time.Time
}

// DefaultConfig is a default CommandConfig implementation.
type DefaultConfig struct {
	writerOutput io.Writer // The writer to use for output
	writerError  io.Writer // The writer to use for errors
	readerInput  io.Reader // The reader to use for input
}

// NewCommandConfig creates a CommandConfig.
func NewCommandConfig(writerOutput, writerError io.Writer, readerInput io.Reader) CommandConfig {
	cfg := DefaultConfig{
		writerOutput: writerOutput,
		writerError:  writerError,
		readerInput:  readerInput,
	}
	return &cfg
}

// WriterOutput returns the output writer
func (cfg *DefaultConfig) WriterOutput() io.Writer {
	return cfg.writerOutput
}

// WriterError returns the output writer
func (cfg *DefaultConfig) WriterError() io.Writer {
	return cfg.writerError
}

// ReaderInput returns the input reader
func (cfg *DefaultConfig) ReaderInput() io.Reader {
	return cfg.readerInput
}

// PasswordStoreDir returns the password store directory.
func (cfg *DefaultConfig) PasswordStoreDir() string {
	storePath := os.Getenv("PASSWORD_STORE_DIR")
	if storePath == "" {
		storePath = path.Join(os.Getenv("HOME"), ".password-store")
	}
	return storePath
}

// PasswordStore returns the PasswordStore.
func (cfg *DefaultConfig) PasswordStore() *store.PasswordStore {
	storePath := cfg.PasswordStoreDir()
	s := store.NewPasswordStore(storePath)
	return s
}

// editor returns the configured editor.
func (cfg *DefaultConfig) editor() string {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "editor"
	}
	return editor
}

// Edit is a callback for asking the user to edit text.
func (cfg *DefaultConfig) Edit(content string) (string, error) {
	file, err := ioutil.TempFile(os.TempDir(), "gopass")
	if err != nil {
		return "", fmt.Errorf("can't create tempfile: %w", err)
	}
	defer file.Close()
	defer os.Remove(file.Name())

	if _, err := file.WriteString(content); err != nil {
		return "", fmt.Errorf("could not write content to tempfile: %w", err)
	}

	editor := exec.Command(cfg.editor(), file.Name())
	editor.Stdout = cfg.WriterOutput()
	editor.Stderr = cfg.WriterError()
	editor.Stdin = cfg.ReaderInput()
	editor.Run()

	if _, err := file.Seek(0, 0); err != nil {
		return "", fmt.Errorf("could not seek to file start: %w", err)
	}

	editedContentBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("could not read edited file: %w", err)
	}

	editedContent := string(editedContentBytes)

	return editedContent, nil
}

func (cfg *DefaultConfig) Now() time.Time {
	return time.Now()
}
