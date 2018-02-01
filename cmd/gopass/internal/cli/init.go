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
	"github.com/aviau/gopass"
	"path/filepath"
)

func execInit(c *commandLine, args []string) error {
	var path, p string

	fs := flag.NewFlagSet("init", flag.ContinueOnError)

	fs.StringVar(&path, "path", getDefaultPasswordStoreDir(c), "")
	fs.StringVar(&p, "p", "", "")

	fs.Usage = func() {
		fmt.Fprintln(c.WriterOutput, `Usage: gopass init [ --path=sub-folder, -p sub-folder ] gpg-id`)
	}

	if err := fs.Parse(args); err != nil {
		return err
	}

	if p != "" {
		path = p
	}

	path, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	if fs.NArg() != 1 {
		fs.Usage()
		return nil
	}

	gpgID := fs.Arg(0)

	store := gopass.NewPasswordStore(path)
	if err := store.Init(gpgID); err != nil {
		return err
	}

	fmt.Fprintf(c.WriterOutput, "Successfully created Password Store at %s\n", path)
	return nil
}
