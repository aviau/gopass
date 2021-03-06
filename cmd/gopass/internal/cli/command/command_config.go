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

package command

import (
	"io"

	"github.com/aviau/gopass"
)

// Config contains everything that commands need to run.
type Config interface {
	WriterOutput() io.Writer
	WriterError() io.Writer
	ReaderInput() io.Reader
	Editor() string
	PasswordStoreDir() string
	PasswordStore() *gopass.PasswordStore
}
