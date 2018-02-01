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
	"os"
	"os/exec"
)

//execEdit rund the "edit" command.
func execEdit(cmd *commandLine, args []string) error {
	fs := flag.NewFlagSet("edit", flag.ExitOnError)
	fs.Parse(args)

	store := cmd.getStore()

	passname := fs.Arg(0)

	password, err := store.GetPassword(passname)

	if err != nil {
		return err
	}

	file, _ := ioutil.TempFile(os.TempDir(), "gopass")
	defer os.Remove(file.Name())

	ioutil.WriteFile(file.Name(), []byte(password), 0600)

	editor := exec.Command(cmd.getEditor(), file.Name())
	editor.Stdout = cmd.WriterOutput
	editor.Stderr = cmd.WriterError
	editor.Stdin = cmd.ReaderInput
	editor.Run()

	pwText, _ := ioutil.ReadFile(file.Name())
	password = string(pwText)

	if err := store.InsertPassword(passname, password); err != nil {
		return err
	}

	fmt.Fprintf(cmd.WriterOutput, "Succesfully edited password '%s'\n", passname)
	return nil
}
