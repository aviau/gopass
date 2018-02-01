//    Copyright (C) 2017-2018 Alexandre Viau <alexandre@alexandreviau.net>
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
	"io"
)

//Run parses the arguments and executes the gopass CLI
func Run(args []string, writerOutput io.Writer, writerError io.Writer, readerInput io.Reader) error {

	//Parse the common flags
	var h, help bool
	var path string
	var editor string

	fs := flag.NewFlagSet("default", flag.ExitOnError)
	fs.StringVar(&path, "PASSWORD_STORE_DIR", "", "Path to the password store")
	fs.StringVar(&editor, "EDITOR", "", "Text editor to use")

	fs.BoolVar(&help, "help", false, "")
	fs.BoolVar(&h, "h", false, "")

	fs.Parse(args)

	c := newCommandline(path, editor, writerOutput, writerError, readerInput)

	if h || help {
		err := execHelp(c)
		return err
	}

	// Retrieve command name as first argument.
	cmd := fs.Arg(0)

	switch cmd {
	case "show":
		return execShow(c, args[1:])
	case "edit":
		return execEdit(c, args[1:])
	case "insert", "add":
		return execInsert(c, args[1:])
	case "find", "ls", "search", "list":
		return execFind(c, args[1:])
	case "":
		return execFind(c, args)
	case "grep":
		return execGrep(c, args[1:])
	case "cp", "copy":
		return execCp(c, args[1:])
	case "mv", "rename":
		return execMv(c, args[1:])
	case "rm", "remove", "delete":
		return execRm(c, args[1:])
	case "generate":
		return execGenerate(c, args[1:])
	case "git":
		return execGit(c, args[1:])
	case "help":
		return execHelp(c)
	case "init":
		return execInit(c, args[1:])
	case "version":
		return execVersion(c)
	default:
		return execShow(c, args)
	}

}
