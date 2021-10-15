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

package cli_test

import (
	"strings"
	"testing"

	"github.com/aviau/gopass/internal/cli/clitest"
	"github.com/stretchr/testify/assert"
)

func TestInsertHelp(t *testing.T) {
	cliTest := clitest.NewCliTest(t)
	defer cliTest.Close()

	result, err := cliTest.Run([]string{"insert", "--help"})

	assert.Nil(t, err)
	assert.Equal(t, "", result.Stderr.String())
	assert.True(t, strings.Contains(result.Stdout.String(), "Usage: gopass insert"))
}

func TestInsertDashH(t *testing.T) {
	cliTest := clitest.NewCliTest(t)
	defer cliTest.Close()

	result, err := cliTest.Run([]string{"insert", "--h"})

	assert.Nil(t, err)
	assert.Equal(t, "", result.Stderr.String())
	assert.True(t, strings.Contains(result.Stdout.String(), "Usage: gopass insert"))
}

func TestInsertMultiline(t *testing.T) {
	cliTest := clitest.NewCliTest(t)
	defer cliTest.Close()

	// Setup the EditFunc callback.
	editFuncCalled := false
	editFunc := func(content string) (string, error) {
		editFuncCalled = true
		return "edited password", nil
	}

	// Ensure the password does not already exists.
	containsPassword, _ := cliTest.PasswordStore().ContainsPassword("test.com")
	assert.False(t, containsPassword, "there should be no preexisting password")

	// Run the command
	_, err := cliTest.Run(
		[]string{"insert", "-m", "test.com"},
		clitest.WithEditFunc(editFunc),
	)

	// Asserts
	assert.Nil(t, err)
	assert.True(t, editFuncCalled, "the editor should have been opened")

	containsPassword, _ = cliTest.PasswordStore().ContainsPassword("test.com")
	assert.True(t, containsPassword, "the password should have been created")

	decryptedPassword, err := cliTest.PasswordStore().GetPassword("test.com")
	assert.Nil(t, err)
	assert.Equal(t, "edited password", decryptedPassword)
}
