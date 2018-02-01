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

	"github.com/aviau/gopass/cmd/gopass/internal/cli/command"
	"github.com/aviau/gopass/cmd/gopass/internal/cli/config"
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

	cfg := config.NewCliConfig(path, editor, writerOutput, writerError, readerInput)

	if h || help {
		err := command.ExecHelp(cfg)
		return err
	}

	// Retrieve command name as first argument.
	cmd := fs.Arg(0)

	switch cmd {
	case "show":
		return command.ExecShow(cfg, args[1:])
	case "edit":
		return command.ExecEdit(cfg, args[1:])
	case "insert", "add":
		return command.ExecInsert(cfg, args[1:])
	case "find", "ls", "search", "list":
		return command.ExecFind(cfg, args[1:])
	case "":
		return command.ExecFind(cfg, args)
	case "grep":
		return command.ExecGrep(cfg, args[1:])
	case "cp", "copy":
		return command.ExecCp(cfg, args[1:])
	case "mv", "rename":
		return command.ExecMv(cfg, args[1:])
	case "rm", "remove", "delete":
		return command.ExecRm(cfg, args[1:])
	case "generate":
		return command.ExecGenerate(cfg, args[1:])
	case "git":
		return command.ExecGit(cfg, args[1:])
	case "help":
		return command.ExecHelp(cfg)
	case "init":
		return command.ExecInit(cfg, args[1:])
	case "version":
		return command.ExecVersion(cfg)
	default:
		return command.ExecShow(cfg, args)
	}

}
