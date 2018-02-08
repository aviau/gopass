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

package find

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"

	"github.com/aviau/gopass/cmd/gopass/internal/cli/command"
)

//ExecFind runs the "find" command.
func ExecFind(cfg command.Config, args []string) error {
	var help, h bool

	fs := flag.NewFlagSet("find", flag.ContinueOnError)
	fs.SetOutput(ioutil.Discard)

	fs.Usage = func() { fmt.Fprintln(cfg.WriterOutput(), "Usage: gopass find patterns...") }

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

	terms := fs.Args()
	pattern := "*" + strings.Join(terms, "*|*") + "*"

	find := exec.Command(
		"tree",
		"-C",
		"-l",
		"--noreport",
		"--prune",       // tree>=1.5
		"--matchdirs",   // tree>=1.7
		"--ignore-case", // tree>=1.7
		"-P",
		pattern,
		store.Path)

	find.Stdout = cfg.WriterOutput()
	find.Stderr = cfg.WriterError()
	find.Run()
	return nil
}
