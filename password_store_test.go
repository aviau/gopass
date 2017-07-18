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
	"github.com/stretchr/testify/assert"
	"testing"

	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/aviau/gopass"
)

func TestGetPasswordsList(t *testing.T) {
	dir, err := ioutil.TempDir("", "gopass")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	_, err = os.Create(filepath.Join(dir, "test.com.gpg"))
	if err != nil {
		t.Fatal(err)
	}

	_, err = os.Create(filepath.Join(dir, "test2.com.gpg"))
	if err != nil {
		t.Fatal(err)
	}

	_, err = os.Create(filepath.Join(dir, "test3"))
	if err != nil {
		t.Fatal(err)
	}

	s := gopass.NewPasswordStore(dir)
	passwords := s.GetPasswordsList()
	assert.Equal(
		t,
		passwords,
		[]string{"test.com", "test2.com"},
		"Password list should contain test.com and test2.com",
	)
}
