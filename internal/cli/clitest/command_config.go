//    Copyright (C) 2021 Alexandre Viau <alexandre@alexandreviau.net>
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

package clitest

import (
	"io"
	"time"

	"github.com/aviau/gopass/pkg/store"
)

// testCommandConfig is a CommandConfig for testing.
type testCommandConfig struct {
	passwordStore *store.PasswordStore
	runOptions    *runOptions
	writerOutput  io.Writer
	writerError   io.Writer
	readerInput   io.Reader
}

func (cfg *testCommandConfig) PasswordStore() *store.PasswordStore {
	return cfg.passwordStore
}

func (cfg *testCommandConfig) PasswordStoreDir() string {
	return cfg.passwordStore.Path
}

func (cfg *testCommandConfig) Edit(content string) (string, error) {
	if cfg.runOptions.editFunc != nil {
		return cfg.runOptions.editFunc(content)
	}
	return content, nil
}

func (cfg *testCommandConfig) WriterOutput() io.Writer {
	return cfg.writerOutput
}

func (cfg *testCommandConfig) WriterError() io.Writer {
	return cfg.writerError
}

func (cfg *testCommandConfig) ReaderInput() io.Reader {
	return cfg.readerInput
}

func (cfg *testCommandConfig) Now() time.Time {
	if cfg.runOptions.nowFunc != nil {
		return cfg.runOptions.nowFunc()
	}
	return time.Now()
}
