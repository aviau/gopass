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

func TestHelp(t *testing.T) {
	cliTest := clitest.NewCliTest(t)
	defer cliTest.Close()

	result, err := cliTest.Run([]string{"help"})

	assert.Nil(t, err)
	assert.Equal(t, "", result.Stderr.String())
	assert.True(t, strings.Contains(result.Stdout.String(), "Usage"))
	assert.True(t, strings.Contains(result.Stdout.String(), "init"))
	assert.True(t, strings.Contains(result.Stdout.String(), "insert"))
}

func TestDashDashHelp(t *testing.T) {
	cliTest := clitest.NewCliTest(t)
	defer cliTest.Close()

	result, err := cliTest.Run([]string{"--help"})

	assert.Nil(t, err)
	assert.Equal(t, "", result.Stderr.String())
	assert.True(t, strings.Contains(result.Stdout.String(), "Usage"))
	assert.True(t, strings.Contains(result.Stdout.String(), "init"))
	assert.True(t, strings.Contains(result.Stdout.String(), "insert"))
}

func TestDashH(t *testing.T) {
	cliTest := clitest.NewCliTest(t)
	defer cliTest.Close()

	result, err := cliTest.Run([]string{"-h"})

	assert.Nil(t, err)
	assert.Equal(t, "", result.Stderr.String())
	assert.True(t, strings.Contains(result.Stdout.String(), "Usage"))
	assert.True(t, strings.Contains(result.Stdout.String(), "init"))
	assert.True(t, strings.Contains(result.Stdout.String(), "insert"))
}
