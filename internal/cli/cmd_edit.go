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
)

// execEdit runs the "edit" command.
func execEdit(cfg CommandConfig, args []string) error {
	var help, h bool

	fs := flag.NewFlagSet("edit", flag.ContinueOnError)
	fs.SetOutput(ioutil.Discard)

	fs.Usage = func() { fmt.Fprintln(cfg.WriterOutput(), "Usage: gopass edit pass-name") }

	fs.BoolVar(&help, "help", false, "")
	fs.BoolVar(&h, "h", false, "")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if help || h {
		fs.Usage()
		return nil
	}

	store := cfg.PasswordStore()

	passName := fs.Arg(0)

	action := "inserted"
	decryptedPassword := ""
	if containsPasword, _ := store.ContainsPassword(passName); containsPasword {
		var err error
		decryptedPassword, err = store.GetPassword(passName)
		if err != nil {
			return err
		}
		action = "edited"
	}

	editedPassword, err := cfg.Edit(decryptedPassword)
	if err != nil {
		return fmt.Errorf("could not edit password: %w", err)
	}

	if err := store.InsertPassword(passName, editedPassword); err != nil {
		return err
	}

	fmt.Fprintf(cfg.WriterOutput(), "Succesfully %s password \"%s\".\n", action, passName)
	return nil
}
