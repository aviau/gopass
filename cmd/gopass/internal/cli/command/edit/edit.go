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

package edit

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/aviau/gopass/cmd/gopass/internal/cli/command"
)

//ExecEdit runs the "edit" command.
func ExecEdit(cfg command.Config, args []string) error {
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

	passname := fs.Arg(0)

	action := "inserted"
	password := ""
	if containsPasword, _ := store.ContainsPassword(passname); containsPasword {
		var err error
		password, err = store.GetPassword(passname)
		if err != nil {
			return err
		}
		action = "edited"
	}

	file, _ := ioutil.TempFile(os.TempDir(), "gopass")
	defer os.Remove(file.Name())

	ioutil.WriteFile(file.Name(), []byte(password), 0600)

	editor := exec.Command(cfg.Editor(), file.Name())
	editor.Stdout = cfg.WriterOutput()
	editor.Stderr = cfg.WriterError()
	editor.Stdin = cfg.ReaderInput()
	editor.Run()

	pwText, _ := ioutil.ReadFile(file.Name())
	password = string(pwText)

	if err := store.InsertPassword(passname, password); err != nil {
		return err
	}

	fmt.Fprintf(cfg.WriterOutput(), "Succesfully %s password \"%s\".\n", action, passname)
	return nil
}
