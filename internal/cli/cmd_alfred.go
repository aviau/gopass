//    Copyright (C) 2022 Alexandre Viau <alexandre@alexandreviau.net>
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
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/aviau/gopass/internal/alfred"
)

func matchesTerms(password string, terms []string) bool {
	searchIn := strings.ToLower(password)

	for _, term := range terms {
		index := strings.Index(searchIn, term)
		if index == -1 {
			return false
		}
		searchIn = searchIn[index+len(term):]
	}
	return true
}

// execAlfred runs the "alfred" command.
func execAlfred(cfg CommandConfig, args []string) error {
	var help, h bool

	fs := flag.NewFlagSet("alfred", flag.ContinueOnError)
	fs.SetOutput(ioutil.Discard)

	fs.Usage = func() { fmt.Fprintln(cfg.WriterOutput(), "Usage: gopass alfred terms...") }

	fs.BoolVar(&help, "help", false, "")
	fs.BoolVar(&h, "h", false, "")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if help || h {
		fs.Usage()
		return nil
	}

	var terms []string
	for _, query := range fs.Args() {
		terms = append(terms, strings.Split(query, " ")...)
	}

	allPasswords := cfg.PasswordStore().GetPasswordsList()

	var matchedPasswords []string
	if len(terms) == 0 {
		matchedPasswords = allPasswords
	} else {
		for _, password := range allPasswords {
			if matchesTerms(password, terms) {
				matchedPasswords = append(matchedPasswords, password)
			}
		}
	}

	var alfredItems = make([]*alfred.Item, 0)
	for _, password := range matchedPasswords {
		alfredItems = append(
			alfredItems,
			&alfred.Item{
				UID:   password,
				Title: password,
				Arg:   password,
			},
		)
	}

	marshaledOutput, err := json.Marshal(
		&alfred.Output{
			Items: alfredItems,
		},
	)
	if err != nil {
		return err
	}

	fmt.Fprintln(cfg.WriterOutput(), string(marshaledOutput))

	return nil
}
