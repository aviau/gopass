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

	cmd_cp "github.com/aviau/gopass/cmd/gopass/internal/cli/command/cp"
	cmd_edit "github.com/aviau/gopass/cmd/gopass/internal/cli/command/edit"
	cmd_find "github.com/aviau/gopass/cmd/gopass/internal/cli/command/find"
	cmd_generate "github.com/aviau/gopass/cmd/gopass/internal/cli/command/generate"
	cmd_git "github.com/aviau/gopass/cmd/gopass/internal/cli/command/git"
	cmd_grep "github.com/aviau/gopass/cmd/gopass/internal/cli/command/grep"
	cmd_help "github.com/aviau/gopass/cmd/gopass/internal/cli/command/help"
	cmd_init "github.com/aviau/gopass/cmd/gopass/internal/cli/command/init"
	cmd_insert "github.com/aviau/gopass/cmd/gopass/internal/cli/command/insert"
	cmd_mv "github.com/aviau/gopass/cmd/gopass/internal/cli/command/mv"
	cmd_rm "github.com/aviau/gopass/cmd/gopass/internal/cli/command/rm"
	cmd_show "github.com/aviau/gopass/cmd/gopass/internal/cli/command/show"
	cmd_version "github.com/aviau/gopass/cmd/gopass/internal/cli/command/version"

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
		err := cmd_help.ExecHelp(cfg)
		return err
	}

	// Retrieve command name as first argument.
	cmd := fs.Arg(0)

	switch cmd {
	case "show":
		return cmd_show.ExecShow(cfg, args[1:])
	case "edit":
		return cmd_edit.ExecEdit(cfg, args[1:])
	case "insert", "add":
		return cmd_insert.ExecInsert(cfg, args[1:])
	case "find", "ls", "search", "list":
		return cmd_find.ExecFind(cfg, args[1:])
	case "":
		return cmd_find.ExecFind(cfg, args)
	case "grep":
		return cmd_grep.ExecGrep(cfg, args[1:])
	case "cp", "copy":
		return cmd_cp.ExecCp(cfg, args[1:])
	case "mv", "rename":
		return cmd_mv.ExecMv(cfg, args[1:])
	case "rm", "remove", "delete":
		return cmd_rm.ExecRm(cfg, args[1:])
	case "generate":
		return cmd_generate.ExecGenerate(cfg, args[1:])
	case "git":
		return cmd_git.ExecGit(cfg, args[1:])
	case "help":
		return cmd_help.ExecHelp(cfg)
	case "init":
		return cmd_init.ExecInit(cfg, args[1:])
	case "version":
		return cmd_version.ExecVersion(cfg)
	default:
		return cmd_show.ExecShow(cfg, args)
	}

}
