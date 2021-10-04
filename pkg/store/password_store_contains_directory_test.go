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

package store_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/aviau/gopass/pkg/internal/gopasstest"
)

func TestContainsDirectory(t *testing.T) {
	st := gopasstest.NewPasswordStoreTest(t)
	defer st.Close()

	containsDirectory, _ := st.PasswordStore.ContainsDirectory("dir")
	assert.False(t, containsDirectory, "The password store should contain dir")

	testDirectoryPath := filepath.Join(st.PasswordStore.Path, "dir")
	if err := os.Mkdir(testDirectoryPath, 0700); err != nil {
		t.Fatal(err)
	}

	_, err := os.Stat(testDirectoryPath)
	assert.Nil(t, err, "the directory should have been created")

	containsDirectory, _ = st.PasswordStore.ContainsDirectory("dir")
	assert.True(t, containsDirectory, "The password store should contain dir")
}

func TestContainsDirectoryTrailingSlash(t *testing.T) {
	st := gopasstest.NewPasswordStoreTest(t)
	defer st.Close()

	containsDirectory, _ := st.PasswordStore.ContainsDirectory("dir")
	assert.False(t, containsDirectory, "The password store should contain dir")

	testDirectoryPath := filepath.Join(st.PasswordStore.Path, "dir")
	if err := os.Mkdir(testDirectoryPath, 0700); err != nil {
		t.Fatal(err)
	}

	_, err := os.Stat(testDirectoryPath)
	assert.Nil(t, err, "the directory should have been created")

	containsDirectory, _ = st.PasswordStore.ContainsDirectory("/dir///")
	assert.True(t, containsDirectory, "The password store should contain dir")
}
