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

	"github.com/aviau/gopass/internal/gopasstest"
)

func TestCopyPassword(t *testing.T) {
	st := gopasstest.NewPasswordStoreTest(t)
	defer st.Close()

	testPasswordPath := filepath.Join(st.PasswordStore.Path, "test.com.gpg")
	_, err := os.Create(testPasswordPath)
	if err != nil {
		t.Fatal(err)
	}

	_, err = os.Stat(testPasswordPath)
	assert.Nil(t, err, "test.com.gpg should have been created")

	destPasswordPath := filepath.Join(st.PasswordStore.Path, "test2.com.gpg")
	_, err = os.Stat(destPasswordPath)
	assert.True(t, os.IsNotExist(err), "test2.com.gpg should not have been created yet")

	st.PasswordStore.CopyPassword("test.com", "test2.com")
	_, err = os.Stat(destPasswordPath)
	assert.Nil(t, err, "test.com.gpg shoudl have been copied to test2.com.gpg")
}

func TestCopyPasswordInDirectory(t *testing.T) {
	st := gopasstest.NewPasswordStoreTest(t)
	defer st.Close()

	testPasswordPath := filepath.Join(st.PasswordStore.Path, "test.com.gpg")
	_, err := os.Create(testPasswordPath)
	if err != nil {
		t.Fatal(err)
	}

	_, err = os.Stat(testPasswordPath)
	assert.Nil(t, err, "test.com.gpg should have been created")

	testDirectoryPath := filepath.Join(st.PasswordStore.Path, "dir")
	if err := os.Mkdir(testDirectoryPath, 0700); err != nil {
		t.Fatal(err)
	}

	destPasswordPath := filepath.Join(st.PasswordStore.Path, "dir", "test.com.gpg")
	_, err = os.Stat(destPasswordPath)
	assert.True(t, os.IsNotExist(err), "test2.com.gpg should have been created yet")

	st.PasswordStore.CopyPassword("test.com", "dir/")
	_, err = os.Stat(destPasswordPath)
	assert.Nil(t, err, "test.com.gpg shoudl have been copied to dir/test2.com.gpg")
}
