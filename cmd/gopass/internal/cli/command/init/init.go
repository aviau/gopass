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

package init

import (
	"flag"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/aviau/gopass"
	"github.com/aviau/gopass/cmd/gopass/internal/cli/command"
)

// ExecInit runs the "init" command.
func ExecInit(cfg command.Config, args []string) error {
	var path, p string
	var help, h bool

	fs := flag.NewFlagSet("init", flag.ContinueOnError)
	fs.SetOutput(ioutil.Discard)

	fs.BoolVar(&help, "help", false, "")
	fs.BoolVar(&h, "h", false, "")

	fs.StringVar(&path, "path", cfg.PasswordStoreDir(), "")
	fs.StringVar(&p, "p", "", "")

	fs.Usage = func() {
		fmt.Fprintln(cfg.WriterOutput(), `Usage: gopass init [--path=subfolder,-p subfolder] gpg-id...`)
	}

	if err := fs.Parse(args); err != nil {
		return err
	}

	if help || h {
		fs.Usage()
		return nil
	}

	if p != "" {
		path = p
	}

	path, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	if fs.NArg() < 1 {
		fs.Usage()
		return nil
	}

	gpgIDs := fs.Args()

	store := gopass.NewPasswordStore(path)
	if err := store.Init(gpgIDs); err != nil {
		return err
	}

	fmt.Fprintf(cfg.WriterOutput(), "Successfully created Password Store at \"%s\".\n", path)
	return nil
}
