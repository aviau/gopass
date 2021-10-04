//    Copyright (C) 2018 Alexandre Viau <alexandre@alexandreviau.net>
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
	"io"
	"os"
	"path"

	"github.com/aviau/gopass"
)

// DefaultConfig is a default Config implementation.
type DefaultConfig struct {
	writerOutput io.Writer // The writer to use for output
	writerError  io.Writer // The writer to use for errors
	readerInput  io.Reader // The reader to use for input
}

// NewConfig creates a Config.
func NewConfig(writerOutput, writerError io.Writer, readerInput io.Reader) *DefaultConfig {
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

// Editor returns the configured editor.
func (cfg *DefaultConfig) Editor() string {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "editor"
	}
	return editor
}

// PasswordStore returns the PasswordStore.
func (cfg *DefaultConfig) PasswordStore() *gopass.PasswordStore {
	storePath := cfg.PasswordStoreDir()
	s := gopass.NewPasswordStore(storePath)
	return s
}
