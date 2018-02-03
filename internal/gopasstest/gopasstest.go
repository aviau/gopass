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

package gopasstest

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/aviau/gopass"
)

//PasswordStoreTest allows for testing password stores.
type PasswordStoreTest struct {
	PasswordStore *gopass.PasswordStore
	StorePath     string
	t             *testing.T
}

//NewPasswordStoreTest creates a password store for testing
func NewPasswordStoreTest(t *testing.T) *PasswordStoreTest {
	storePath, err := ioutil.TempDir("", "gopass")
	if err != nil {
		t.Fatal(err)
	}

	passwordStore := gopass.NewPasswordStore(storePath)
	passwordStore.UsesGit = false

	if err := passwordStore.Init("test"); err != nil {
		t.Fatal(err)
	}

	passwordStoreTest := PasswordStoreTest{
		PasswordStore: passwordStore,
		StorePath:     storePath,
	}

	return &passwordStoreTest
}

//Close removes the password store
func (test *PasswordStoreTest) Close() {
	if err := os.RemoveAll(test.StorePath); err != nil {
		test.t.Fatal(err)
	}
}
