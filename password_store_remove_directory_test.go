//   Copyright (C) 2017 Alexandre Viau <alexandre@alexandreviau.net>
//
//   This file is part of gopass.
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
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRemoveDirectory(t *testing.T) {
	st, err := newPasswordStoreTest()
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()

	testDirectoryPath := filepath.Join(st.StorePath, "dir")
	err = os.Mkdir(testDirectoryPath, 0700)
	if err != nil {
		t.Fatal(err)
	}

	_, err = os.Stat(testDirectoryPath)
	assert.Nil(t, err, "dir should have been created")

	err = st.PasswordStore.RemoveDirectory("dir")
	if err != nil {
		t.Fatal(err)
	}

	_, err = os.Stat(testDirectoryPath)
	assert.True(t, os.IsNotExist(err), "dir should have been removed")
}

func TestRemoveDirectoryTrailingSlash(t *testing.T) {
	st, err := newPasswordStoreTest()
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()

	testDirectoryPath := filepath.Join(st.StorePath, "dir")
	err = os.Mkdir(testDirectoryPath, 0700)
	if err != nil {
		t.Fatal(err)
	}

	_, err = os.Stat(testDirectoryPath)
	assert.Nil(t, err, "dir should have been created")

	err = st.PasswordStore.RemoveDirectory("//dir///")
	if err != nil {
		t.Fatal(err)
	}

	_, err = os.Stat(testDirectoryPath)
	assert.True(t, os.IsNotExist(err), "dir should have been removed")
}

func TestRemoveDirectoryRecursive(t *testing.T) {
	st, err := newPasswordStoreTest()
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()

	testDirectoryPath := filepath.Join(st.StorePath, "dir")
	err = os.Mkdir(testDirectoryPath, 0700)
	if err != nil {
		t.Fatal(err)
	}

	testPasswordPath := filepath.Join(testDirectoryPath, "test.com.gpg")
	_, err = os.Create(testPasswordPath)
	if err != nil {
		t.Fatal(err)
	}

	_, err = os.Stat(testPasswordPath)
	assert.Nil(t, err, "test.com.gpg should have been created")

	_, err = os.Stat(testDirectoryPath)
	assert.Nil(t, err, "dir should have been created")

	err = st.PasswordStore.RemoveDirectory("dir")
	if err != nil {
		t.Fatal(err)
	}

	_, err = os.Stat(testDirectoryPath)
	assert.True(t, os.IsNotExist(err), "dir should have been removed")

	_, err = os.Stat(testPasswordPath)
	assert.True(t, os.IsNotExist(err), "test.com.gpg should have been removed")
}
