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
	"github.com/aviau/gopass/cmd/gopass/internal/pwgen"
	gopass_terminal "github.com/aviau/gopass/cmd/gopass/internal/terminal"
	"strconv"
)

//execGenerate runs the "generate" command.
func execGenerate(cmd *commandLine, args []string) error {
	var noSymbols, n bool
	var force, f bool

	fs := flag.NewFlagSet("generate", flag.ContinueOnError)
	fs.Usage = func() {
		fmt.Fprintln(cmd.WriterOutput, `Usage: gopass generate [--no-symbols,-n] [--force,-f] pass-name pass-length`)
	}

	fs.BoolVar(&noSymbols, "no-symbols", false, "")
	fs.BoolVar(&n, "n", false, "")

	fs.BoolVar(&force, "force", false, "")
	fs.BoolVar(&f, "f", false, "")

	err := fs.Parse(args)
	if err != nil {
		return err
	}

	noSymbols = noSymbols || n
	force = force || f

	passName := fs.Arg(0)

	store := getStore(cmd)

	if containsPassword, _ := store.ContainsPassword(passName); containsPassword && !force {
		if !gopass_terminal.AskYesNo(cmd.WriterOutput, fmt.Sprintf("Password '%s' already exists. Would you like to overwrite? [y/n] ", passName)) {
			return nil
		}
	}

	passLength, err := strconv.ParseInt(fs.Arg(1), 0, 64)
	if err != nil {
		fmt.Fprintf(cmd.WriterOutput, "Second argument must be an int, got '%s'\n", fs.Arg(1))
		return err
	}

	runes := append(pwgen.Alpha, pwgen.Num...)
	if !noSymbols {
		runes = append(runes, pwgen.Symbols...)
	}

	password := pwgen.RandSeq(int(passLength), runes)

	if err := store.InsertPassword(passName, password); err != nil {
		return err
	}

	fmt.Fprintf(cmd.WriterOutput, "Password %s added to the store\n", passName)
	return nil
}
