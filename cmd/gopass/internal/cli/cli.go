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
	"os"
	"path"

	"github.com/aviau/gopass"
)

//commandLine holds options from the main parser
type commandLine struct {
	Path         string    //Path to the password store
	Editor       string    //Text editor to use
	WriterOutput io.Writer //The writer to use for output
}

//Run parses the arguments and executes the gopass CLI
func Run(args []string, writerOutput io.Writer) error {
	c := commandLine{}
	c.WriterOutput = writerOutput

	//Parse the common flags
	var h, help bool

	fs := flag.NewFlagSet("default", flag.ExitOnError)
	fs.StringVar(&c.Path, "PASSWORD_STORE_DIR", "", "Path to the password store")
	fs.StringVar(&c.Editor, "EDITOR", "", "Text editor to use")

	fs.BoolVar(&help, "help", false, "")
	fs.BoolVar(&h, "h", false, "")

	fs.Parse(args)

	if h || help {
		err := execHelp(&c)
		return err
	}

	// Retrieve command name as first argument.
	cmd := fs.Arg(0)

	switch cmd {
	case "show":
		return execShow(&c, args[1:])
	case "edit":
		return execEdit(&c, args[1:])
	case "insert", "add":
		return execInsert(&c, args[1:])
	case "find", "ls", "search", "list":
		return execFind(&c, args[1:])
	case "":
		return execFind(&c, args)
	case "grep":
		return execGrep(&c, args[1:])
	case "cp", "copy":
		return execCp(&c, args[1:])
	case "mv", "rename":
		return execMv(&c, args[1:])
	case "rm", "remove", "delete":
		return execRm(&c, args[1:])
	case "generate":
		return execGenerate(&c, args[1:])
	case "git":
		return execGit(&c, args[1:])
	case "help":
		return execHelp(&c)
	case "init":
		return execInit(&c, args[1:])
	case "version":
		return execVersion(&c)
	default:
		return execShow(&c, args)
	}

}

func getDefaultPasswordStoreDir(c *commandLine) string {
	//Look for the store path in the commandLine,
	// env var, or default to $HOME/.password-store
	storePath := c.Path
	if storePath == "" {
		storePath = os.Getenv("PASSWORD_STORE_DIR")
		if storePath == "" {
			storePath = path.Join(os.Getenv("HOME"), ".password-store")
		}
	}
	return storePath
}

func getEditor(c *commandLine) string {
	// Look for the editor to use in the commandLine,
	// env var, or default to editor.
	editor := c.Editor
	if editor == "" {
		editor = os.Getenv("EDITOR")
		if editor == "" {
			editor = "editor"
		}
	}
	return editor
}

//getStore finds and returns the PasswordStore
func getStore(c *commandLine) *gopass.PasswordStore {
	storePath := getDefaultPasswordStoreDir(c)
	s := gopass.NewPasswordStore(storePath)
	return s
}
