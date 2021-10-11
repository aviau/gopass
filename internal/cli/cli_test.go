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
	"context"
	"testing"

	"github.com/aviau/gopass/internal/gopasstest"
	"github.com/aviau/gopass/pkg/store"
)

// testConfig is a fake CommandConfig.
type testConfig struct {
	CommandConfig
	passwordStoreTest *gopasstest.PasswordStoreTest
	editFunc          func(string) (string, error)
}

func (cfg *testConfig) PasswordStore() *store.PasswordStore {
	return cfg.passwordStoreTest.PasswordStore
}

func (cfg *testConfig) Edit(content string) (string, error) {
	return cfg.editFunc(content)
}

// cliTest allows for testing the CLI without a TTY.
type cliTest struct {
	OutputWriter      bytes.Buffer
	ErrorWriter       bytes.Buffer
	t                 *testing.T
	passwordStoreTest *gopasstest.PasswordStoreTest
	EditFunc          func(string) (string, error)
}

func newCliTest(t *testing.T) *cliTest {
	passwordStoreTest := gopasstest.NewPasswordStoreTest(t)

	cliTest := cliTest{
		t:                 t,
		passwordStoreTest: passwordStoreTest,
		EditFunc: func(content string) (string, error) {
			return content, nil
		},
	}
	return &cliTest
}

func (cliTest *cliTest) PasswordStore() *store.PasswordStore {
	return cliTest.passwordStoreTest.PasswordStore
}

func (cliTest *cliTest) Run(args []string) error {
	baseConfig := NewCommandConfig(&cliTest.OutputWriter, &cliTest.ErrorWriter, nil)

	testConfig := testConfig{
		CommandConfig:     baseConfig,
		passwordStoreTest: cliTest.passwordStoreTest,
		editFunc:          cliTest.EditFunc,
	}

	return Run(context.TODO(), &testConfig, args)
}

func (cliTest *cliTest) Close() {
	cliTest.passwordStoreTest.Close()
}
