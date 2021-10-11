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

// Package clipboard implements utilities for interating with the clipboard.
package clipboard

import (
	"os/exec"
)

// CopyToClipboard copies a string to the clipboard using xclip
func CopyToClipboard(s string) error {
	xclip := exec.Command(
		"xclip",
		"-in",
		"-selection",
		"clipboard",
	)

	stdin, err := xclip.StdinPipe()
	if err != nil {
		return err
	}

	if err := xclip.Start(); err != nil {
		return err
	}

	_, err = stdin.Write([]byte(s))
	if err != nil {
		return err
	}

	if err := stdin.Close(); err != nil {
		return err
	}

	err = xclip.Wait()

	return err
}
