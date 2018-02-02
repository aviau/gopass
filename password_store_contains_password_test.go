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
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContainsPassword(t *testing.T) {
	st := newPasswordStoreTest(t)
	defer st.Close()

	containsPassword, _ := st.PasswordStore.ContainsPassword("test.com")
	assert.False(t, containsPassword, "The password store should not contain test.com")

	testPasswordPath := filepath.Join(st.StorePath, "test.com.gpg")
	_, err := os.Create(testPasswordPath)
	if err != nil {
		t.Fatal(err)
	}

	_, err = os.Stat(testPasswordPath)
	assert.Nil(t, err, "test.com.gpg should have been created")

	containsPassword, _ = st.PasswordStore.ContainsPassword("test.com")
	assert.True(t, containsPassword, "The password store should contain test.com")
}

func TestContainsPasswordDirectory(t *testing.T) {
	st := newPasswordStoreTest(t)
	defer st.Close()

	containsPassword, _ := st.PasswordStore.ContainsPassword("test.com")
	assert.False(t, containsPassword, "The password store should not contain test.com")

	testDirectoryPath := filepath.Join(st.StorePath, "test.com")
	if err := os.Mkdir(testDirectoryPath, 0700); err != nil {
		t.Fatal(err)
	}

	_, err := os.Stat(testDirectoryPath)
	assert.Nil(t, err, "test.com directory should have been created")

	containsPassword, _ = st.PasswordStore.ContainsPassword("test.com")
	assert.False(t, containsPassword, "The password store should not a password named test.com")
}
