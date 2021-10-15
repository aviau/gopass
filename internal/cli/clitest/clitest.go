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

// Package clitest provides utilities for testing the CLI.
package clitest

import (
	"bytes"
	"context"
	"io"
	"os"
	"testing"
	"time"

	"github.com/aviau/gopass/internal/cli"
	"github.com/aviau/gopass/internal/storetest"
	"github.com/aviau/gopass/pkg/store"
)

// testConfig is a CommandConfig for testing.
type testConfig struct {
	passwordStore *store.PasswordStore
	editFunc      func(string) (string, error)
	nowFunc       func() time.Time
	writerOutput  io.Writer
	writerError   io.Writer
	readerInput   io.Reader
}

func (cfg *testConfig) PasswordStore() *store.PasswordStore {
	return cfg.passwordStore
}

func (cfg *testConfig) PasswordStoreDir() string {
	return cfg.passwordStore.Path
}

func (cfg *testConfig) Edit(content string) (string, error) {
	return cfg.editFunc(content)
}

func (cfg *testConfig) WriterOutput() io.Writer {
	return cfg.writerOutput
}

func (cfg *testConfig) WriterError() io.Writer {
	return cfg.writerError
}

func (cfg *testConfig) ReaderInput() io.Reader {
	return cfg.readerInput
}

func (cfg *testConfig) Now() time.Time {
	return cfg.nowFunc()
}

// cliTest allows for testing the CLI without a TTY.
type cliTest struct {
	t                 *testing.T
	passwordStoreTest *storetest.PasswordStoreTest
	EditFunc          func(string) (string, error)
	NowFunc           func() time.Time
}

func NewCliTest(t *testing.T) *cliTest {
	passwordStoreTest := storetest.NewPasswordStoreTest(t)

	cliTest := cliTest{
		t:                 t,
		passwordStoreTest: passwordStoreTest,
		EditFunc: func(content string) (string, error) {
			return content, nil
		},
		NowFunc: func() time.Time {
			return time.Now()
		},
	}

	return &cliTest
}

func (cliTest *cliTest) PasswordStore() *store.PasswordStore {
	return cliTest.passwordStoreTest.PasswordStore
}

type runResult struct {
	Stdout *bytes.Buffer
	Stderr *bytes.Buffer
}

func (cliTest *cliTest) Run(args []string) (*runResult, error) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	testConfig := &testConfig{
		passwordStore: cliTest.PasswordStore(),
		editFunc:      cliTest.EditFunc,
		nowFunc:       cliTest.NowFunc,
		writerOutput:  stdout,
		writerError:   stderr,
		readerInput:   os.Stdin,
	}

	err := cli.Run(context.TODO(), testConfig, args)

	runResult := &runResult{
		Stdout: stdout,
		Stderr: stderr,
	}

	return runResult, err
}

func (cliTest *cliTest) Close() {
	cliTest.passwordStoreTest.Close()
}
