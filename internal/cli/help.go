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
	"fmt"
)

// execHelp runs the "help" command.
func execHelp(cfg CommandConfig) error {
	fmt.Fprint(cfg.WriterOutput(), `Usage:
      init                  Initialize a new password store.
      ls                    List passwords.
      find                  List passwords that match a string.
      show                  Show an encryped password.
      grep                  Search for a string in all passwords.
      insert                Insert a new password.
      edit                  Edit an existing password.
      generate              Generate a new password.
      rm                    Remove a password.
      mv                    Move a password.
      cp                    Copy a password.
      git                   Execute a git command.
      help                  Show this text.
      version               Show version information.
`)
	return nil
}
