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
	"github.com/aviau/gopass/clipboard"
	"strings"
)

//execShow runs the "show" command.
func execShow(cmd *commandLine, args []string) error {
	var clip, c bool

	fs := flag.NewFlagSet("show", flag.ExitOnError)
	fs.Usage = func() { fmt.Fprintln(cmd.WriterOutput, `Usage: gopass show [pass-name]`) }

	fs.BoolVar(&clip, "clip", false, "")
	fs.BoolVar(&c, "c", false, "")

	if err := fs.Parse(args); err != nil {
		return err
	}

	clip = clip || c

	password := fs.Arg(0)

	store := cmd.getStore()

	password, err := store.GetPassword(password)
	if err != nil {
		return err
	}

	if clip {
		firstPasswordLine := strings.Split(password, "\n")[0]

		if err := clipboard.CopyToClipboard(firstPasswordLine); err != nil {
			return err
		}

		fmt.Fprintln(cmd.WriterOutput, "the first line of the password was copied to clipboard.")
	} else {
		fmt.Fprintln(cmd.WriterOutput, password)
	}

	return nil
}
