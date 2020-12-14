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

package show

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/aviau/gopass/cmd/gopass/internal/cli/command"
	"github.com/aviau/gopass/internal/clipboard"
)

var usernameRegex = regexp.MustCompile(`(username|user|email):\s*(?P<username>\S*)`)

// ExecShow runs the "show" command.
func ExecShow(cfg command.Config, args []string) error {
	var clip, c bool
	var username, u bool
	var help, h bool

	fs := flag.NewFlagSet("show", flag.ContinueOnError)
	fs.SetOutput(ioutil.Discard)

	fs.Usage = func() {
		fmt.Fprintln(cfg.WriterOutput(), "Usage: gopass show [pass-name]")
	}

	fs.BoolVar(&help, "help", false, "")
	fs.BoolVar(&h, "h", false, "")

	fs.BoolVar(&clip, "clip", false, "")
	fs.BoolVar(&c, "c", false, "")

	fs.BoolVar(&username, "username", false, "")
	fs.BoolVar(&u, "u", false, "")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if help || h {
		fs.Usage()
		return nil
	}

	clip = clip || c

	username = username || u

	password := fs.Arg(0)

	if password == "" {
		return errors.New("missing password name")
	}

	store := cfg.PasswordStore()

	// Decrypt the password
	password, err := store.GetPassword(password)
	if err != nil {
		return err
	}

	// Prepare the password to display or copy.
	outputPassword := password
	if username {
		if matches := usernameRegex.FindStringSubmatch(password); matches != nil {
			outputPassword = matches[2]
		} else {
			return fmt.Errorf("Could not find username in the password")
		}
	} else if clip {
		outputPassword = strings.Split(password, "\n")[0]
	}

	// Eithier display the password or copy it to the clipboard.
	if clip {
		if err := clipboard.CopyToClipboard(outputPassword); err != nil {
			return err
		}
		fmt.Fprintln(cfg.WriterOutput(), "the first line of the password was copied to clipboard.")
	} else {
		fmt.Fprintln(cfg.WriterOutput(), outputPassword)
	}

	return nil
}
