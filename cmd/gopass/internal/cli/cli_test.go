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
	"bytes"
	"testing"

	"github.com/aviau/gopass"
	"github.com/aviau/gopass/cmd/gopass/internal/cli/command"
	"github.com/aviau/gopass/internal/gopasstest"
)

// testConfig is a fake command.Config, it does not use env variables.
type testConfig struct {
	command.Config
	passwordStoreTest *gopasstest.PasswordStoreTest
}

func (cfg *testConfig) PasswordStore() *gopass.PasswordStore {
	return cfg.passwordStoreTest.PasswordStore
}

type cliTest struct {
	OutputWriter      bytes.Buffer
	ErrorWriter       bytes.Buffer
	t                 *testing.T
	passwordStoreTest *gopasstest.PasswordStoreTest
}

func newCliTest(t *testing.T) *cliTest {
	passwordStoreTest := gopasstest.NewPasswordStoreTest(t)

	cliTest := cliTest{
		t:                 t,
		passwordStoreTest: passwordStoreTest,
	}
	return &cliTest
}

func (cliTest *cliTest) PasswordStore() *gopass.PasswordStore {
	return cliTest.passwordStoreTest.PasswordStore
}

func (cliTest *cliTest) Run(args []string) error {
	baseConfig := command.NewConfig(&cliTest.OutputWriter, &cliTest.ErrorWriter, nil)

	testConfig := testConfig{
		Config:            baseConfig,
		passwordStoreTest: cliTest.passwordStoreTest,
	}

	return runCommand(&testConfig, args)
}

func (cliTest *cliTest) Close() {
	cliTest.passwordStoreTest.Close()
}
