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

	"github.com/aviau/gopass/internal/gopasstest"
)

func TestGetPasswordsList(t *testing.T) {
	st := gopasstest.NewPasswordStoreTest(t)
	defer st.Close()

	if _, err := os.Create(filepath.Join(st.PasswordStore.Path, "test.com.gpg")); err != nil {
		t.Fatal(err)
	}

	if _, err := os.Create(filepath.Join(st.PasswordStore.Path, "test2.com.gpg")); err != nil {
		t.Fatal(err)
	}

	if _, err := os.Create(filepath.Join(st.PasswordStore.Path, "test3")); err != nil {
		t.Fatal(err)
	}

	dirPath := filepath.Join(st.PasswordStore.Path, "dir")
	if err := os.Mkdir(dirPath, os.ModePerm); err != nil {
		t.Fatal(err)
	}

	if _, err := os.Create(filepath.Join(dirPath, "test3.com.gpg")); err != nil {
		t.Fatal(err)
	}

	passwords := st.PasswordStore.GetPasswordsList()
	assert.Equal(
		t,
		passwords,
		[]string{"dir/test3.com", "test.com", "test2.com"},
		"Password list should contain test.com and test2.com",
	)
}
