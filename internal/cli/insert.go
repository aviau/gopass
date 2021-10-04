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

	"golang.org/x/crypto/ssh/terminal"

	gopass_terminal "github.com/aviau/gopass/internal/terminal"
)

// execInsert runs the "insert" command.
func execInsert(cfg CommandConfig, args []string) error {
	var multiline, m bool
	var force, f bool
	var help, h bool

	fs := flag.NewFlagSet("insert", flag.ContinueOnError)
	fs.SetOutput(ioutil.Discard)

	fs.BoolVar(&help, "help", false, "")
	fs.BoolVar(&h, "h", false, "")

	fs.BoolVar(&multiline, "multiline", false, "")
	fs.BoolVar(&m, "m", false, "")

	fs.BoolVar(&force, "force", false, "")
	fs.BoolVar(&f, "f", false, "")

	fs.Usage = func() {
		fmt.Fprintln(cfg.WriterOutput(), `Usage: gopass insert [ --multiline, -m ] [ --force, -f ] pass-name`)
	}

	if err := fs.Parse(args); err != nil {
		return err
	}

	if help || h {
		fs.Usage()
		return nil
	}

	multiline = multiline || m
	force = force || f

	pwname := fs.Arg(0)

	store := cfg.PasswordStore()

	// Check if password already exists
	if containsPassword, _ := store.ContainsPassword(pwname); containsPassword && !force {
		if !gopass_terminal.AskYesNo(cfg.WriterOutput(), fmt.Sprintf("Password \"%s\" already exists. Would you like to overwrite? [y/n] ", pwname)) {
			return nil
		}
	}

	var password string

	if multiline {
		file, _ := ioutil.TempFile(os.TempDir(), "gopass")
		defer os.Remove(file.Name())

		editor := exec.Command(cfg.Editor(), file.Name())
		editor.Stdout = os.Stdout
		editor.Stdin = os.Stdin
		editor.Run()

		pwText, _ := ioutil.ReadFile(file.Name())
		password = string(pwText)
	} else {
		fd := int(os.Stdin.Fd())
		for {
			fmt.Fprintln(cfg.WriterOutput(), "Enter password:")
			try1, _ := terminal.ReadPassword(fd)
			fmt.Fprintln(cfg.WriterOutput())

			fmt.Fprintln(cfg.WriterOutput(), "Enter confirmation:")
			try2, _ := terminal.ReadPassword(fd)
			fmt.Fprintln(cfg.WriterOutput())

			if string(try1) == string(try2) {
				password = string(try1)
				break
			} else {
				fmt.Fprintln(cfg.WriterOutput(), "Passwords did not match, try again...")
			}

		}
	}

	if err := store.InsertPassword(pwname, password); err != nil {
		return err
	}

	fmt.Fprintf(cfg.WriterOutput(), "Password \"%s\" added to the store.\n", pwname)
	return nil
}
