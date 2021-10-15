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
	"time"

	"github.com/aviau/gopass/internal/cli/clitest"
	"github.com/stretchr/testify/assert"
)

func TestShowDashDashHelp(t *testing.T) {
	cliTest := clitest.NewCliTest(t)
	defer cliTest.Close()

	result, err := cliTest.Run([]string{"show", "--help"})

	assert.Nil(t, err)
	assert.Equal(t, "", result.Stderr.String())
	assert.True(t, strings.Contains(result.Stdout.String(), "Usage: gopass show"))
}

func TestShowDashH(t *testing.T) {
	cliTest := clitest.NewCliTest(t)
	defer cliTest.Close()

	result, err := cliTest.Run([]string{"show", "-h"})

	assert.Nil(t, err)
	assert.Equal(t, "", result.Stderr.String())
	assert.True(t, strings.Contains(result.Stdout.String(), "Usage: gopass show"))
}

func TestShowMissingPassword(t *testing.T) {
	cliTest := clitest.NewCliTest(t)
	defer cliTest.Close()

	result, err := cliTest.Run([]string{"show"})

	assert.EqualError(t, err, "missing password name")
	assert.Equal(t, result.Stdout.String(), "")
	assert.Equal(t, result.Stderr.String(), "")
}

func TestShow(t *testing.T) {
	cliTest := clitest.NewCliTest(t)
	defer cliTest.Close()

	if err := cliTest.PasswordStore().InsertPassword("test.com", "hello world"); err != nil {
		t.Fatal(err)
	}

	result, err := cliTest.Run([]string{"show", "test.com"})

	assert.Nil(t, err)
	assert.Equal(t, result.Stderr.String(), "")
	assert.Equal(t, result.Stdout.String(), "hello world\n")
}

func TestShowTwoFactor(t *testing.T) {
	cliTest := clitest.NewCliTest(t)
	defer cliTest.Close()

	passwordWithTwoFactor := `pass123
---
2fa: JBSWY3DPEHPK3PXP
`

	if err := cliTest.PasswordStore().InsertPassword("test.com", passwordWithTwoFactor); err != nil {
		t.Fatal(err)
	}

	cliTest.NowFunc = func() time.Time {
		return time.Date(2020, 1, 2, 15, 0, 0, 0, time.UTC)
	}

	result, err := cliTest.Run([]string{"show", "--2fa", "test.com"})

	assert.Nil(t, err)
	assert.Equal(t, result.Stderr.String(), "")
	assert.Equal(t, result.Stdout.String(), "891690\n")

}

func TestShowTwoFactorOtpauthURI(t *testing.T) {
	cliTest := clitest.NewCliTest(t)
	defer cliTest.Close()

	passwordWithTwoFactor := `pass123
--
2fa: otpauth://totp/Example:alice@google.com?secret=JBSWY3DPEHPK3PXP&issuer=Example
`

	if err := cliTest.PasswordStore().InsertPassword("test.com", passwordWithTwoFactor); err != nil {
		t.Fatal(err)
	}

	cliTest.NowFunc = func() time.Time {
		return time.Date(2020, 1, 2, 15, 0, 0, 0, time.UTC)
	}

	result, err := cliTest.Run([]string{"show", "--2fa", "test.com"})

	assert.Nil(t, err)
	assert.Equal(t, result.Stderr.String(), "")
	assert.Equal(t, result.Stdout.String(), "891690\n")

}
