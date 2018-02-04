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

package command

import (
	"io"
	"os"
	"path"

	"github.com/aviau/gopass"
)

//Config contains everything that commands need to run.
type Config struct {
	WriterOutput io.Writer //The writer to use for output
	WriterError  io.Writer //The writer to use for errors
	ReaderInput  io.Reader //The reader to use for input
}

//NewConfig creates a Config.
func NewConfig(writerOutput, writerError io.Writer, readerInput io.Reader) *Config {
	cfg := Config{
		WriterOutput: writerOutput,
		WriterError:  writerError,
		ReaderInput:  readerInput,
	}
	return &cfg
}

//GetDefaultPasswordStoreDir returns the default password store directory.
func (cfg *Config) GetDefaultPasswordStoreDir() string {
	storePath := os.Getenv("PASSWORD_STORE_DIR")
	if storePath == "" {
		storePath = path.Join(os.Getenv("HOME"), ".password-store")
	}
	return storePath
}

//GetEditor returns the configured editor.
func (cfg *Config) GetEditor() string {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "editor"
	}
	return editor
}

//GetStore finds and returns the PasswordStore.
func (cfg *Config) GetStore() *gopass.PasswordStore {
	storePath := cfg.GetDefaultPasswordStoreDir()
	s := gopass.NewPasswordStore(storePath)
	return s
}
