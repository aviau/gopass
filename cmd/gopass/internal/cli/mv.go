//    Copyright (C) 2018 Alexandre Viau <alexandre@alexandreviau.net>
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

package cli

import (
	"flag"
	"fmt"
	"path/filepath"
	"strings"
)

//execMv runs the "mv" comand.
func execMv(c *commandLine, args []string) error {
	var force, f bool

	fs := flag.NewFlagSet("mv", flag.ExitOnError)
	fs.Usage = func() { fmt.Fprintln(c.WriterOutput, "Usage: gopass mv old-path new-path") }

	fs.BoolVar(&force, "force", false, "")
	fs.BoolVar(&f, "f", false, "")

	if err := fs.Parse(args); err != nil {
		return err
	}

	force = force || f

	store := getStore(c)

	source := fs.Arg(0)
	dest := fs.Arg(1)

	if source == "" || dest == "" {
		fmt.Fprintln(c.WriterOutput, "Error: Received empty source or dest argument")
		return nil
	}

	// If the dest ends with a '/', then it is a directory.
	if strings.HasSuffix(dest, "/") {
		_, sourceFile := filepath.Split(source)
		dest = filepath.Join(dest, sourceFile)
	}

	if sourceIsPassword, _ := store.ContainsPassword(source); sourceIsPassword {

		if destAlreadyExists, _ := store.ContainsPassword(dest); destAlreadyExists {
			if !force {
				fmt.Fprintf(c.WriterOutput, "Error: destination %s already exists. Use -f to override\n", dest)
				return nil
			}
		}

		if err := store.MovePassword(source, dest); err != nil {
			return err
		}

		fmt.Fprintf(c.WriterOutput, "Moved password from '%s' to '%s'\n", source, dest)
		return nil
	}

	if sourceIsDirectory, _ := store.ContainsDirectory(source); sourceIsDirectory {
		if err := store.MoveDirectory(source, dest); err != nil {
			return err
		}
		fmt.Fprintf(c.WriterOutput, "Moved directory from '%s' to '%s'\n", source, dest)
		return nil
	}

	fmt.Fprintf(c.WriterOutput, "Error: could not find source '%s' to copy \n", source)
	return nil
}
