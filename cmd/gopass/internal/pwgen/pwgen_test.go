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

package pwgen_test

import (
	"github.com/aviau/gopass/cmd/gopass/internal/pwgen"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestLength(t *testing.T) {
	runes := pwgen.Alpha
	password := pwgen.RandSeq(10, runes)
	assert.Equal(t, len(password), 10)

	password = pwgen.RandSeq(30, runes)
	assert.Equal(t, len(password), 30)
}

func TestContainsNumsOnly(t *testing.T) {
	runes := pwgen.Num
	password := pwgen.RandSeq(50, runes)
	numsString := string(pwgen.Num)
	for _, character := range password {
		assert.True(t, strings.ContainsRune(numsString, character))
	}
}

func TestContainsAlphaOnly(t *testing.T) {
	runes := pwgen.Alpha
	password := pwgen.RandSeq(50, runes)
	numsString := string(pwgen.Alpha)
	for _, character := range password {
		assert.True(t, strings.ContainsRune(numsString, character))
	}
}
