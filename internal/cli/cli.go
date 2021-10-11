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

// Package cli implements the gopass CLI.
package cli

import "context"

// Run parses the arguments and executes the gopass CLI
func Run(ctx context.Context, cfg CommandConfig, cmdAndArgs []string) error {
	cmd := ""
	if len(cmdAndArgs) > 0 {
		cmd = cmdAndArgs[0]
	}

	switch cmd {
	case "show":
		return execShow(cfg, cmdAndArgs[1:])
	case "edit":
		return execEdit(cfg, cmdAndArgs[1:])
	case "insert", "add":
		return execInsert(cfg, cmdAndArgs[1:])
	case "find", "ls", "search", "list":
		return execFind(cfg, cmdAndArgs[1:])
	case "":
		return execFind(cfg, cmdAndArgs)
	case "grep":
		return execGrep(cfg, cmdAndArgs[1:])
	case "cp", "copy":
		return execCp(cfg, cmdAndArgs[1:])
	case "mv", "rename":
		return execMv(cfg, cmdAndArgs[1:])
	case "rm", "remove", "delete":
		return execRm(cfg, cmdAndArgs[1:])
	case "generate":
		return execGenerate(cfg, cmdAndArgs[1:])
	case "git":
		return execGit(cfg, cmdAndArgs[1:])
	case "help", "-h", "--help":
		return execHelp(cfg)
	case "init":
		return execInit(cfg, cmdAndArgs[1:])
	case "version":
		return execVersion(cfg)
	default:
		return execShow(cfg, cmdAndArgs)
	}

}
