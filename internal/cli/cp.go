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
	"io/ioutil"
	"path/filepath"
	"strings"
)

// execCp runs the "cp" command.
func execCp(cfg CommandConfig, args []string) error {
	var recursive, r bool
	var force, f bool
	var help, h bool

	fs := flag.NewFlagSet("cp", flag.ContinueOnError)
	fs.SetOutput(ioutil.Discard)

	fs.Usage = func() { fmt.Fprintln(cfg.WriterOutput(), "Usage: gopass cp old-path new-path") }

	fs.BoolVar(&help, "help", false, "")
	fs.BoolVar(&h, "h", false, "")

	fs.BoolVar(&recursive, "recursive", false, "")
	fs.BoolVar(&r, "r", false, "")

	fs.BoolVar(&force, "force", false, "")
	fs.BoolVar(&f, "f", false, "")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if help || h {
		fs.Usage()
		return nil
	}

	recursive = recursive || r
	force = force || f

	store := cfg.PasswordStore()

	source := fs.Arg(0)
	dest := fs.Arg(1)

	if source == "" || dest == "" {
		return fmt.Errorf("received empty source or dest argument")
	}

	// If the dest ends with a '/', then it is a directory.
	if strings.HasSuffix(dest, "/") {
		_, sourceFile := filepath.Split(source)
		dest = filepath.Join(dest, sourceFile)
	}

	if sourceIsPassword, _ := store.ContainsPassword(source); sourceIsPassword {

		if destAlreadyExists, _ := store.ContainsPassword(dest); destAlreadyExists {
			if !force {
				return fmt.Errorf("destination \"%s\" already exists. Use -f to override", dest)
			}
		}

		if err := store.CopyPassword(source, dest); err != nil {
			return err
		}

		fmt.Fprintf(cfg.WriterOutput(), "Copied password from \"%s\" to \"%s\".\n", source, dest)
		return nil
	}

	if sourceIsDirectory, _ := store.ContainsDirectory(source); sourceIsDirectory {

		if !recursive {
			return fmt.Errorf("\"%s\" is a directory, use -r to copy recursively", source)
		}

		if err := store.CopyDirectory(source, dest); err != nil {
			return err
		}

		fmt.Fprintf(cfg.WriterOutput(), "Copied directory from \"%s\" to \"%s\".\n", source, dest)
		return nil
	}

	return fmt.Errorf("could not find source \"%s\" to copy", source)
}
