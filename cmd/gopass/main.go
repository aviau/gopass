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

package main

import (
	"fmt"
	"os"

	"github.com/aviau/gopass/cmd/gopass/internal/cli"
)

func main() {
	// Retrieve args and Shift binary name off argument list.
	args := os.Args[1:]

	if err := cli.Run(args, os.Stdout, os.Stderr, os.Stdin); err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
}
