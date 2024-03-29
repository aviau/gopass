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

// Package terminal implements utilities for interacting with the terminal.
package terminal

import (
	"fmt"
	"io"
)

// AskYesNo asks a Yes/No question to the user and returns the result
func AskYesNo(writer io.Writer, question string) bool {
	fmt.Fprint(writer, question)

	var response string
	fmt.Scanln(&response)

	switch response {
	case "y", "Y", "yes", "YES":
		return true
	default:
		return false
	}

}
