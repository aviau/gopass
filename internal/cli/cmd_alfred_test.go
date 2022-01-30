//    Copyright (C) 2022 Alexandre Viau <alexandre@alexandreviau.net>
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
	"encoding/json"
	"strings"
	"testing"

	"github.com/aviau/gopass/internal/alfred"
	"github.com/aviau/gopass/internal/cli/clitest"
	"github.com/stretchr/testify/assert"
)

func TestAlfredHelp(t *testing.T) {
	cliTest := clitest.NewCliTest(t)
	defer cliTest.Close()

	result, err := cliTest.Run([]string{"alfred", "--help"})

	assert.Nil(t, err)
	assert.Equal(t, "", result.Stderr.String())
	assert.True(t, strings.Contains(result.Stdout.String(), "Usage: gopass alfred"))
}

func TestAlfred(t *testing.T) {
	cliTest := clitest.NewCliTest(t)
	defer cliTest.Close()

	if err := cliTest.PasswordStore().InsertPassword("aaaa", ""); err != nil {
		t.Fatal(err)
	}

	if err := cliTest.PasswordStore().InsertPassword("aabbcc", ""); err != nil {
		t.Fatal(err)
	}

	type testCase struct {
		query        []string
		exoectedUIDs []string
	}

	testCases := []*testCase{
		{
			query:        []string{"aa"},
			exoectedUIDs: []string{"aaaa", "aabbcc"},
		},
		{
			query:        []string{"aa", "cc"},
			exoectedUIDs: []string{"aabbcc"},
		},
		{
			query:        []string{"cc", "aa"},
			exoectedUIDs: []string{},
		},
		{
			query:        []string{"bb", "bb"},
			exoectedUIDs: []string{},
		},
	}

	for _, testCase := range testCases {
		result, err := cliTest.Run(
			append([]string{"alfred"}, testCase.query...),
		)

		assert.Nil(t, err)
		assert.Equal(t, "", result.Stderr.String())

		parsedOutput := alfred.Output{}
		if err := json.Unmarshal(result.Stdout.Bytes(), &parsedOutput); err != nil {
			t.Fatal(err)
		}

		var UIDs = make([]string, 0)
		for _, item := range parsedOutput.Items {
			UIDs = append(UIDs, item.UID)
		}

		assert.Equal(t, testCase.exoectedUIDs, UIDs)
	}

}
