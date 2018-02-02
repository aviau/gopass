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

package rm

import (
	"flag"
	"fmt"

	"github.com/aviau/gopass/cmd/gopass/internal/cli/command"
	gopass_terminal "github.com/aviau/gopass/cmd/gopass/internal/terminal"
)

//ExecRm runs the "rm" command.
func ExecRm(cfg *command.Config, args []string) error {
	var recursive, r bool
	var force, f bool

	fs := flag.NewFlagSet("rm", flag.ContinueOnError)
	fs.Usage = func() {
		fmt.Fprintln(cfg.WriterOutput, "Usage: gopass rm pass-name")
	}

	fs.BoolVar(&recursive, "recursive", false, "")
	fs.BoolVar(&r, "r", false, "")

	fs.BoolVar(&force, "force", false, "")
	fs.BoolVar(&f, "f", false, "")

	err := fs.Parse(args)
	if err != nil {
		return err
	}

	force = force || f
	recursive = recursive || r

	store := cfg.GetStore()

	pwname := fs.Arg(0)
	if pwname == "" {
		fs.Usage()
		return nil
	}

	if containsPassword, _ := store.ContainsPassword(pwname); containsPassword {

		if !force {
			if !gopass_terminal.AskYesNo(cfg.WriterOutput, fmt.Sprintf("Are you sure you would like to delete %s? [y/n] ", pwname)) {
				return nil
			}
		}

		if err := store.RemovePassword(pwname); err != nil {
			return err
		}

	} else if containsDirectory, _ := store.ContainsDirectory(pwname); containsDirectory {

		if !recursive {
			return fmt.Errorf("\"%s\" is a directory, use -r to remove recursively", pwname)
		}

		if !force {
			if !gopass_terminal.AskYesNo(cfg.WriterOutput, fmt.Sprintf("Are you sure you would like to delete \"%s\" recursively? [y/n] ", pwname)) {
				return nil
			}
		}

		if err := store.RemoveDirectory(pwname); err != nil {
			return err
		}
	}

	fmt.Fprintf(cfg.WriterOutput, "Removed password/directory at path \"%s\".\n", fs.Arg(0))
	return nil
}
