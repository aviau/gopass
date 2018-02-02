//    Copyright (C) 2017 Alexandre Viau <alexandre@alexandreviau.net>
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

package gopass_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/aviau/gopass"
)

type passwordStoreTest struct {
	PasswordStore *gopass.PasswordStore
	StorePath     string
	t             *testing.T
}

func newPasswordStoreTest(t *testing.T) *passwordStoreTest {
	storePath, err := ioutil.TempDir("", "gopass")
	if err != nil {
		t.Fatal(err)
	}

	passwordStore := gopass.NewPasswordStore(storePath)
	passwordStore.UsesGit = false

	if err := passwordStore.Init("test"); err != nil {
		t.Fatal(err)
	}

	passwordStoreTest := passwordStoreTest{
		PasswordStore: passwordStore,
		StorePath:     storePath,
	}

	return &passwordStoreTest
}

func (test *passwordStoreTest) Close() {
	if err := os.RemoveAll(test.StorePath); err != nil {
		test.t.Fatal(err)
	}
}
