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

	"github.com/stretchr/testify/assert"
)

func TestShowDashDashHelp(t *testing.T) {
	cliTest := newCliTest(t)
	defer cliTest.Close()

	cliTest.Run([]string{"show", "--help"})

	assert.True(t, strings.Contains(cliTest.OutputWriter.String(), "Usage: gopass show"))
}

func TestShowDashH(t *testing.T) {
	cliTest := newCliTest(t)
	defer cliTest.Close()

	cliTest.Run([]string{"show", "-h"})

	assert.True(t, strings.Contains(cliTest.OutputWriter.String(), "Usage: gopass show"))
}

func TestShowMissingPassword(t *testing.T) {
	cliTest := newCliTest(t)
	defer cliTest.Close()

	err := cliTest.Run([]string{"show"})

	assert.EqualError(t, err, "missing password name")
	assert.Equal(t, cliTest.OutputWriter.String(), "")
	assert.Equal(t, cliTest.ErrorWriter.String(), "")
}
