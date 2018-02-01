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

package mv

import (
	"flag"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/aviau/gopass/cmd/gopass/internal/cli/config"
)

//ExecMv runs the "mv" comand.
func ExecMv(cfg *config.CliConfig, args []string) error {
	var force, f bool

	fs := flag.NewFlagSet("mv", flag.ExitOnError)
	fs.Usage = func() { fmt.Fprintln(cfg.WriterOutput, "Usage: gopass mv old-path new-path") }

	fs.BoolVar(&force, "force", false, "")
	fs.BoolVar(&f, "f", false, "")

	if err := fs.Parse(args); err != nil {
		return err
	}

	force = force || f

	store := cfg.GetStore()

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

		if err := store.MovePassword(source, dest); err != nil {
			return err
		}

		fmt.Fprintf(cfg.WriterOutput, "Moved password from \"%s\" to \"%s\".\n", source, dest)
		return nil
	}

	if sourceIsDirectory, _ := store.ContainsDirectory(source); sourceIsDirectory {
		if err := store.MoveDirectory(source, dest); err != nil {
			return err
		}
		fmt.Fprintf(cfg.WriterOutput, "Moved directory from \"%s\" to \"%s\".\n", source, dest)
		return nil
	}

	return fmt.Errorf("could not find source \"%s\" to copy", source)
}
