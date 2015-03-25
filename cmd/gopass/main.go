//   This file is part of gopass.
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

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strconv"
	"strings"

	"github.com/ReAzem/gopass"
	"github.com/ReAzem/gopass/pwgen"
	"github.com/mgutz/ansi"
	"golang.org/x/crypto/ssh/terminal"
)

//CommandLine holds options from the main parser
type CommandLine struct {
	Path   string //Path to the password store
	Editor string //Text editor to use
}

func main() {

	// Shift binary name off argument list.
	args := os.Args[1:]

	// Retrieve command name as first argument.
	var cmd string
	if len(args) > 0 && !strings.HasPrefix(args[0], "-") {
		cmd = args[0]
	}

	c := CommandLine{}

	//Parse the common flags
	fs := flag.NewFlagSet("default", flag.ExitOnError)
	fs.StringVar(&c.Path, "PASSWORD_STORE_DIR", "", "Path to the password store")
	fs.StringVar(&c.Editor, "EDITOR", "editor", "Text editor to use")
	fs.Usage = func() { execHelp() }
	fs.Parse(args)

	switch cmd {
	case "show":
		execShow(&c, args[1:])
	case "edit":
		execEdit(&c, args[1:])
	case "insert":
		execInsert(&c, args[1:])
	case "find", "ls":
		execFind(&c, args[1:])
	case "":
		execFind(&c, args)
	case "grep":
		execGrep(&c, args[1:])
	case "cp":
		fmt.Println("Executing", cmd)
	case "mv":
		execMv(&c, args[1:])
	case "rm":
		execRm(&c, args[1:])
	case "generate":
		execGenerate(&c, args[1:])
	case "git":
		execGit(&c, args[1:])
	case "help":
		execHelp()
	case "init":
		fmt.Println("Executing", cmd)
	default:
		execShow(&c, args)
	}

}

func execHelp() {
	fmt.Println(`Usage:
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
}

//execInsert runs the "insert" command.
func execInsert(c *CommandLine, args []string) {
	var multiline, m bool
	var force, f bool

	fs := flag.NewFlagSet("insert", flag.ExitOnError)

	fs.BoolVar(&multiline, "multiline", false, "")
	fs.BoolVar(&m, "m", false, "")

	fs.BoolVar(&force, "force", false, "")
	fs.BoolVar(&f, "f", false, "")

	fs.Usage = func() { fmt.Println(`Usage: gopass insert [ --multiline, -m ] [ --force, -f ] pass-name`) }
	fs.Parse(args)

	multiline = multiline || m
	force = force || f

	pwname := fs.Arg(0)

	store := GetStore(c)

	passwordPath := path.Join(store.Path, pwname+".gpg")

	// Check if password already exists
	if _, err := os.Stat(passwordPath); err == nil && !force {
		fmt.Printf("Error: Password already exists at '%s', use -f to force\n", passwordPath)
		return
	}

	var password string

	if multiline {
		file, _ := ioutil.TempFile(os.TempDir(), "gopass")
		defer os.Remove(file.Name())

		editor := exec.Command(c.Editor, file.Name())
		editor.Stdout = os.Stdout
		editor.Stdin = os.Stdin
		editor.Run()

		pwText, _ := ioutil.ReadFile(file.Name())
		password = string(pwText)
	} else {
		fd := int(os.Stdin.Fd())
		for {
			fmt.Println("Enter password:")
			try1, _ := terminal.ReadPassword(fd)
			fmt.Println()

			fmt.Println("Enter confirmation:")
			try2, _ := terminal.ReadPassword(fd)
			fmt.Println()

			if string(try1) == string(try2) {
				password = string(try1)
				break
			} else {
				fmt.Println("Passwords did not match, try again...")
			}

		}
	}

	err := store.InsertPassword(pwname, password)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Password %s added to the store\n", pwname)
}

//execEdit rund the "edit" command.
func execEdit(cmd *CommandLine, args []string) {
	fs := flag.NewFlagSet("edit", flag.ExitOnError)
	fs.Parse(args)

	store := GetStore(cmd)

	passname := fs.Arg(0)

	password, err := store.GetPassword(passname)

	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}

	file, _ := ioutil.TempFile(os.TempDir(), "gopass")
	defer os.Remove(file.Name())

	ioutil.WriteFile(file.Name(), []byte(password), 0600)

	editor := exec.Command(cmd.Editor, file.Name())
	editor.Stdout = os.Stdout
	editor.Stdin = os.Stdin
	editor.Run()

	pwText, _ := ioutil.ReadFile(file.Name())
	password = string(pwText)

	err = store.InsertPassword(passname, password)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Succesfully edited password '%s'\n", passname)
}

//execGenerate runs the "generate" command.
func execGenerate(cmd *CommandLine, args []string) {
	var noSymbols, n bool
	var clip, c bool
	var inPlace, i bool
	var force, f bool

	fs := flag.NewFlagSet("generate", flag.ExitOnError)
	fs.Usage = func() {
		fmt.Println(`Usage: pass generate [--no-symbols,-n] [--clip,-c] [--in-place,-i | --force,-f] pass-name pass-length`)
	}

	fs.BoolVar(&noSymbols, "no-symbols", false, "")
	fs.BoolVar(&n, "n", false, "")

	fs.BoolVar(&clip, "clip", false, "")
	fs.BoolVar(&c, "c", false, "")

	fs.BoolVar(&inPlace, "in-place", false, "")
	fs.BoolVar(&i, "i", false, "")

	fs.BoolVar(&force, "force", false, "")
	fs.BoolVar(&f, "f", false, "")

	fs.Parse(args)

	noSymbols = noSymbols || n
	clip = clip || c
	inPlace = inPlace || i
	force = force || f

	passName := fs.Arg(0)
	passLength, err := strconv.ParseInt(fs.Arg(1), 0, 64)
	if err != nil {
		fmt.Printf("Second argument must be an int, got '%s'\n", fs.Arg(1))
		os.Exit(1)
	}

	runes := append(pwgen.Alpha, pwgen.Num...)
	if !noSymbols {
		runes = append(runes, pwgen.Symbols...)
	}

	password := pwgen.RandSeq(int(passLength), runes)

	store := GetStore(cmd)

	store.InsertPassword(passName, password)
	fmt.Printf("Password %s added to the store\n", passName)
}

//execRm runs the "rm" command.
func execRm(c *CommandLine, args []string) {
	fs := flag.NewFlagSet("rm", flag.ExitOnError)
	fs.Parse(args)

	store := GetStore(c)

	err := store.Remove(fs.Arg(0))
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	} else {
		fmt.Println("Removed password/directory at path", fs.Arg(0))
	}
}

//execMv runs the "mv" comand.
func execMv(c *CommandLine, args []string) {
	fs := flag.NewFlagSet("mv", flag.ExitOnError)
	fs.Parse(args)

	store := GetStore(c)

	source := fs.Arg(0)
	dest := fs.Arg(1)

	if source == "" || dest == "" {
		fmt.Println("Error: Received empty source or dest argument")
		os.Exit(1)
	}

	err := store.Move(source, dest)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	} else {
		fmt.Printf("Moved password/directory from '%s' to '%s'\n", source, dest)
	}
}

//execShow runs the "show" command.
func execShow(c *CommandLine, args []string) {
	fs := flag.NewFlagSet("show", flag.ExitOnError)
	fs.Usage = func() { fmt.Println(`Usage: gopass show [--clip,-c] [pass-name]`) }
	fs.Parse(args)
	password := fs.Arg(0)

	store := GetStore(c)

	password, err := store.GetPassword(password)

	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}

	fmt.Println(password)
}

//execFind runs the "find" command.
func execFind(c *CommandLine, args []string) {
	fs := flag.NewFlagSet("find", flag.ExitOnError)
	fs.Parse(args)

	store := GetStore(c)

	terms := fs.Args()
	pattern := "*" + strings.Join(terms, "*|*") + "*"

	find := exec.Command(
		"tree",
		"-C",
		"-l",
		"--noreport",
		"--prune",       // tree>=1.5
		"--matchdirs",   // tree>=1.7
		"--ignore-case", // tree>=1.7
		"-P",
		pattern,
		store.Path)

	find.Stdout = os.Stdout
	find.Stderr = os.Stderr
	find.Run()
}

//execGrep runs the "grep" command
func execGrep(c *CommandLine, args []string) {
	fs := flag.NewFlagSet("grep", flag.ExitOnError)
	fs.Parse(args)

	pattern, _ := regexp.CompilePOSIX(fs.Arg(0))

	store := GetStore(c)

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
			fmt.Printf("%s:\n%s", ansi.Color(password, "cyan+b"), output)
		}
	}
}

//execGit runs the "git" command
func execGit(c *CommandLine, args []string) {
	//TODO: Add work-dir arg.
	git := exec.Command("git", args...)
	git.Stdout = os.Stdout
	git.Stderr = os.Stderr
	git.Run()
}

//GetStore finds and returns the PasswordStore
func GetStore(c *CommandLine) *gopass.PasswordStore {
	//Look for the store path in the CommandLine,
	// env var, or default to $HOME/.password-store
	storePath := c.Path
	if storePath == "" {
		storePath = os.Getenv("PASSWORD_STORE_DIR")
		if storePath == "" {
			storePath = path.Join(os.Getenv("HOME"), ".password-store")
		}
	}

	s := gopass.NewPasswordStore(storePath)
	return s
}
