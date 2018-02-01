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

package config

import (
	"github.com/aviau/gopass"
	"io"
	"os"
	"path"
)

//CliConfig holds options from the main parser.
type CliConfig struct {
	Path         string    //Path to the password store
	Editor       string    //Text editor to use
	WriterOutput io.Writer //The writer to use for output
	WriterError  io.Writer //The writer to use for errors
	ReaderInput  io.Reader //The reader to use for input
}

//NewCliConfig creates a CliConfig.
func NewCliConfig(path, editor string, writerOutput, writerError io.Writer, readerInput io.Reader) *CliConfig {
	cfg := CliConfig{
		Path:         path,
		Editor:       editor,
		WriterOutput: writerOutput,
		WriterError:  writerError,
		ReaderInput:  readerInput,
	}
	return &cfg
}

//GetDefaultPasswordStoreDir returns the default password store directory.
func (cfg *CliConfig) GetDefaultPasswordStoreDir() string {
	//Look for the store path in the commandLine,
	// env var, or default to $HOME/.password-store
	storePath := cfg.Path
	if storePath == "" {
		storePath = os.Getenv("PASSWORD_STORE_DIR")
		if storePath == "" {
			storePath = path.Join(os.Getenv("HOME"), ".password-store")
		}
	}
	return storePath
}

//GetEditor returns the configured editor.
func (cfg *CliConfig) GetEditor() string {
	// Look for the editor to use in the commandLine,
	// env var, or default to editor.
	editor := cfg.Editor
	if editor == "" {
		editor = os.Getenv("EDITOR")
		if editor == "" {
			editor = "editor"
		}
	}
	return editor
}

//GetStore finds and returns the PasswordStore.
func (cfg *CliConfig) GetStore() *gopass.PasswordStore {
	storePath := cfg.GetDefaultPasswordStoreDir()
	s := gopass.NewPasswordStore(storePath)
	return s
}
