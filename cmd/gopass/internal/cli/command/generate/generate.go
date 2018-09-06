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

package generate

import (
	"flag"
	"fmt"
	"github.com/aviau/gopass/internal/clipboard"
	"io/ioutil"
	"strconv"

	"github.com/aviau/gopass/cmd/gopass/internal/cli/command"
	gopass_terminal "github.com/aviau/gopass/cmd/gopass/internal/terminal"
	"github.com/aviau/gopass/internal/pwgen"
)

// ExecGenerate runs the "generate" command.
func ExecGenerate(cfg command.Config, args []string) error {
	var noSymbols, n bool
	var force, f bool
	var help, h bool
	var write, w bool

	fs := flag.NewFlagSet("generate", flag.ContinueOnError)
	fs.SetOutput(ioutil.Discard)

	fs.Usage = func() {
		fmt.Fprintln(cfg.WriterOutput(), `Usage: gopass generate [--no-symbols,-n] [--force,-f] [--write,-w][pass-name] pass-length`)
	}

	fs.BoolVar(&help, "help", false, "")
	fs.BoolVar(&h, "h", false, "")

	fs.BoolVar(&noSymbols, "no-symbols", false, "")
	fs.BoolVar(&n, "n", false, "")

	fs.BoolVar(&force, "force", false, "")
	fs.BoolVar(&f, "f", false, "")

	fs.BoolVar(&write, "write", false, "")
	fs.BoolVar(&w, "w", false, "")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if help || h {
		fs.Usage()
		return nil
	}

	noSymbols = noSymbols || n
	force = force || f
	write = write || w

	var passLength int
	if write {
		l, err := strconv.ParseInt(fs.Arg(1), 0, 64)
		if err != nil {
			return fmt.Errorf("password length argument must be an int, got \"%s\"", fs.Arg(1))
		}
		passLength = int(l)
	} else {
		l, err := strconv.ParseInt(fs.Arg(0), 0, 64)
		if err != nil {
			return fmt.Errorf("password length argument must be an int, got \"%s\"", fs.Arg(0))
		}
		passLength = int(l)
	}

	runes := append(pwgen.Alpha, pwgen.Num...)
	if !noSymbols {
		runes = append(runes, pwgen.Symbols...)
	}

	password := pwgen.RandSeq(int(passLength), runes)

	if (write) {
		store := cfg.PasswordStore()
		passName := fs.Arg(0)
		if passName == "" {
			return fmt.Errorf("password name should be empty")
		}

		if containsPassword, _ := store.ContainsPassword(passName); containsPassword && !force {
			if !gopass_terminal.AskYesNo(cfg.WriterOutput(), fmt.Sprintf("Password \"%s\" already exists. Would you like to overwrite? [y/n] ", passName)) {
				return nil
			}
		}

		if err := store.InsertPassword(passName, password); err != nil {
			return err
		}
		fmt.Fprintf(cfg.WriterOutput(), "Password \"%s\" added to the store.\n", passName)

	} else { //copy
		if err := clipboard.CopyToClipboard(password); err != nil {
			return err
		}
		fmt.Fprintln(cfg.WriterOutput(), "Password was copied to clipboard.")
	}

	return nil
}
