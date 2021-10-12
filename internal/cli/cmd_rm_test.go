//    Copyright (C) 2017-2018 Alexandre Viau <alexandre@alexandreviau.net>
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

package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRmDashDashHelp(t *testing.T) {
	cliTest := newCliTest(t)
	defer cliTest.Close()

	result, err := cliTest.Run([]string{"rm", "--help"})

	assert.Nil(t, err)
	assert.Equal(t, "", result.Stderr.String())
	assert.True(t, strings.Contains(result.Stdout.String(), "Usage: gopass rm"))
}

func TestRmDashH(t *testing.T) {
	cliTest := newCliTest(t)
	defer cliTest.Close()

	result, err := cliTest.Run([]string{"rm", "-h"})

	assert.Nil(t, err)
	assert.Equal(t, "", result.Stderr.String())
	assert.True(t, strings.Contains(result.Stdout.String(), "Usage: gopass rm"))
}

func TestRmDirectoryWithoutRecursive(t *testing.T) {
	cliTest := newCliTest(t)
	defer cliTest.Close()

	testDirectoryPath := filepath.Join(cliTest.PasswordStore().Path, "dir")
	if err := os.Mkdir(testDirectoryPath, 0700); err != nil {
		t.Fatal(err)
	}

	_, err := cliTest.Run([]string{"rm", "dir"})

	assert.EqualError(t, err, "\"dir\" is a directory, use -r to remove recursively")
}

func TestRmUnexistingPassword(t *testing.T) {
	cliTest := newCliTest(t)
	defer cliTest.Close()

	_, err := cliTest.Run([]string{"rm", "dir"})

	assert.EqualError(t, err, "could not find password or directory to remove")
}
