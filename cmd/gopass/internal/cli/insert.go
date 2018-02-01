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
	gopass_terminal "github.com/aviau/gopass/cmd/gopass/internal/terminal"
	"golang.org/x/crypto/ssh/terminal"
	"io/ioutil"
	"os"
	"os/exec"
)

//execInsert runs the "insert" command.
func execInsert(cmd *commandLine, args []string) error {
	var multiline, m bool
	var force, f bool

	fs := flag.NewFlagSet("insert", flag.ContinueOnError)

	fs.BoolVar(&multiline, "multiline", false, "")
	fs.BoolVar(&m, "m", false, "")

	fs.BoolVar(&force, "force", false, "")
	fs.BoolVar(&f, "f", false, "")

	fs.Usage = func() {
		fmt.Fprintln(cmd.WriterOutput, `Usage: gopass insert [ --multiline, -m ] [ --force, -f ] pass-name`)
	}
	err := fs.Parse(args)
	if err != nil {
		return err
	}

	multiline = multiline || m
	force = force || f

	pwname := fs.Arg(0)

	store := cmd.getStore()

	// Check if password already exists
	if containsPassword, _ := store.ContainsPassword(pwname); containsPassword && !force {
		if !gopass_terminal.AskYesNo(cmd.WriterOutput, fmt.Sprintf("Password '%s' already exists. Would you like to overwrite? [y/n] ", pwname)) {
			return nil
		}
	}

	var password string

	if multiline {
		file, _ := ioutil.TempFile(os.TempDir(), "gopass")
		defer os.Remove(file.Name())

		editor := exec.Command(cmd.getEditor(), file.Name())
		editor.Stdout = os.Stdout
		editor.Stdin = os.Stdin
		editor.Run()

		pwText, _ := ioutil.ReadFile(file.Name())
		password = string(pwText)
	} else {
		fd := int(os.Stdin.Fd())
		for {
			fmt.Fprintln(cmd.WriterOutput, "Enter password:")
			try1, _ := terminal.ReadPassword(fd)
			fmt.Fprintln(cmd.WriterOutput)

			fmt.Fprintln(cmd.WriterOutput, "Enter confirmation:")
			try2, _ := terminal.ReadPassword(fd)
			fmt.Fprintln(cmd.WriterOutput)

			if string(try1) == string(try2) {
				password = string(try1)
				break
			} else {
				fmt.Fprintln(cmd.WriterOutput, "Passwords did not match, try again...")
			}

		}
	}

	if err := store.InsertPassword(pwname, password); err != nil {
		return err
	}

	fmt.Fprintf(cmd.WriterOutput, "Password %s added to the store\n", pwname)
	return nil
}
