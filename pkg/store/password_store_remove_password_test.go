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

func TestRemovePassword(t *testing.T) {
	st := gopasstest.NewPasswordStoreTest(t)
	defer st.Close()

	testPasswordPath := filepath.Join(st.PasswordStore.Path, "test.com.gpg")
	_, err := os.Create(testPasswordPath)
	if err != nil {
		t.Fatal(err)
	}

	_, err = os.Stat(testPasswordPath)
	assert.Nil(t, err, "test.com.gpg should have been created")

	st.PasswordStore.RemovePassword("test.com")
	_, err = os.Stat(testPasswordPath)
	assert.True(t, os.IsNotExist(err), "test.com should have been removed")
}

func TestRemovePasswordTrailingSlash(t *testing.T) {
	st := gopasstest.NewPasswordStoreTest(t)
	defer st.Close()

	testPasswordPath := filepath.Join(st.PasswordStore.Path, "test.com.gpg")
	_, err := os.Create(testPasswordPath)
	if err != nil {
		t.Fatal(err)
	}

	_, err = os.Stat(testPasswordPath)
	assert.Nil(t, err, "test.com.gpg should have been created")

	st.PasswordStore.RemovePassword("test.com/")
	_, err = os.Stat(testPasswordPath)
	assert.Nil(t, err, "RemovePassword with a trailing slash should not remove a password")
}

func TestRemovePasswordDirectory(t *testing.T) {
	st := gopasstest.NewPasswordStoreTest(t)
	defer st.Close()

	testDirectoryPath := filepath.Join(st.PasswordStore.Path, "test.com")
	if err := os.Mkdir(testDirectoryPath, 0700); err != nil {
		t.Fatal(err)
	}

	_, err := os.Stat(testDirectoryPath)
	assert.Nil(t, err, "test.com.gpg should have been created")

	st.PasswordStore.RemovePassword("test.com")
	_, err = os.Stat(testDirectoryPath)
	assert.False(
		t,
		os.IsNotExist(err),
		"RemovePassword should not remove directories",
	)
}
