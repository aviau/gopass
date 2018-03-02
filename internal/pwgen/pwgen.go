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

package pwgen

import (
	"crypto/rand"
	"math/big"
)

// Alpha is a-Z and A-Z
var Alpha = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// Num is 0-9
var Num = []rune("0123456789")

// Symbols is all printable symbols
var Symbols = []rune("!#$%&'()*+,-./:;<=>?@[]^_`{|}~")

// RandSeq returns a random sequence of lenght n
func RandSeq(n int, runes []rune) string {
	b := make([]rune, n)
	for i := range b {
		randomInt, err := rand.Int(rand.Reader, big.NewInt(int64(len(runes))))
		if err != nil {
			panic(err)
		}
		b[i] = runes[randomInt.Int64()]
	}
	return string(b)
}
