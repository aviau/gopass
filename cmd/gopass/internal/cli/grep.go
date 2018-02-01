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
	"github.com/mgutz/ansi"
	"regexp"
	"strings"
)

//execGrep runs the "grep" command
func execGrep(cmd *commandLine, args []string) error {
	fs := flag.NewFlagSet("grep", flag.ExitOnError)
	fs.Parse(args)

	pattern, _ := regexp.CompilePOSIX(fs.Arg(0))

	store := cmd.getStore()

	passwords := store.GetPasswordsList()

	for _, password := range passwords {
		decryptedPassword, _ := store.GetPassword(password)
		lines := strings.Split(decryptedPassword, "\n")
		output := ""
		for _, line := range lines {
			result := pattern.FindAllString(line, -1)
			if len(result) > 0 {
				output += strings.Replace(line+"\n", result[0], ansi.Color(result[0], "red+b"), -1)
			}
		}
		if output != "" {
			fmt.Fprintf(cmd.WriterOutput, "%s:\n%s", ansi.Color(password, "cyan+b"), output)
		}
	}
	return nil
}
