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

func TestGetPasswordsList(t *testing.T) {
	st, err := newPasswordStoreTest()
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()

	_, err = os.Create(filepath.Join(st.StorePath, "test.com.gpg"))
	if err != nil {
		t.Fatal(err)
	}

	_, err = os.Create(filepath.Join(st.StorePath, "test2.com.gpg"))
	if err != nil {
		t.Fatal(err)
	}

	_, err = os.Create(filepath.Join(st.StorePath, "test3"))
	if err != nil {
		t.Fatal(err)
	}

	passwords := st.PasswordStore.GetPasswordsList()
	assert.Equal(
		t,
		passwords,
		[]string{"test.com", "test2.com"},
		"Password list should contain test.com and test2.com",
	)
}
