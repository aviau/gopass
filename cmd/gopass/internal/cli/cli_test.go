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

package cli_test

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/aviau/gopass/cmd/gopass/internal/cli"
	"github.com/aviau/gopass/cmd/gopass/internal/version"
)

func TestVersion(t *testing.T) {
	var writer bytes.Buffer

	cli.Run([]string{"version"}, &writer)

	assert.Equal(t,
		writer.String(),
		fmt.Sprintf("gopass v%s\n", version.Version))
}

func TestHelp(t *testing.T) {
	var writer bytes.Buffer

	cli.Run([]string{"help"}, &writer)
	assert.True(t, strings.Contains(writer.String(), "Usage"))

	writer.Reset()

	cli.Run([]string{"--help"}, &writer)
	assert.True(t, strings.Contains(writer.String(), "Usage"))

	writer.Reset()

	cli.Run([]string{"-h"}, &writer)
	assert.True(t, strings.Contains(writer.String(), "Usage"))
}

func TestInitHelp(t *testing.T) {
	var writer bytes.Buffer

	cli.Run([]string{"init", "--help"}, &writer)
	assert.True(t, strings.Contains(writer.String(), "Usage: gopass init"))

	writer.Reset()

	cli.Run([]string{"init", "-h"}, &writer)
	assert.True(t, strings.Contains(writer.String(), "Usage: gopass init"))
}

func TestInsertHelp(t *testing.T) {
	var writer bytes.Buffer

	cli.Run([]string{"insert", "--help"}, &writer)
	assert.True(t, strings.Contains(writer.String(), "Usage: gopass insert"))

	writer.Reset()

	cli.Run([]string{"insert", "-h"}, &writer)
	assert.True(t, strings.Contains(writer.String(), "Usage: gopass insert"))
}

func TestGenerateHelp(t *testing.T) {
	var writer bytes.Buffer

	cli.Run([]string{"generate", "--help"}, &writer)
	assert.True(t, strings.Contains(writer.String(), "Usage: gopass generate"))

	writer.Reset()

	cli.Run([]string{"generate", "-h"}, &writer)
	assert.True(t, strings.Contains(writer.String(), "Usage: gopass generate"))
}

func TestRmHelp(t *testing.T) {
	var writer bytes.Buffer

	cli.Run([]string{"rm", "--help"}, &writer)
	assert.True(t, strings.Contains(writer.String(), "Usage: gopass rm"))

	writer.Reset()

	cli.Run([]string{"rm", "-h"}, &writer)
	assert.True(t, strings.Contains(writer.String(), "Usage: gopass rm"))
}
