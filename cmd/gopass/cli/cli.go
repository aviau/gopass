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

package cli

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/aviau/gopass"
	"github.com/aviau/gopass/pwgen"
	"github.com/aviau/gopass/version"
	"github.com/mgutz/ansi"
	"golang.org/x/crypto/ssh/terminal"
)

//commandLine holds options from the main parser
type commandLine struct {
	Path         string    //Path to the password store
	Editor       string    //Text editor to use
	WriterOutput io.Writer //The writer to use for output
}

//Run parses the arguments and executes the gopass CLI
func Run(args []string, writerOutput io.Writer) {
	c := commandLine{}
	c.WriterOutput = writerOutput

	//Parse the common flags
	var h, help bool

	fs := flag.NewFlagSet("default", flag.ExitOnError)
	fs.StringVar(&c.Path, "PASSWORD_STORE_DIR", "", "Path to the password store")
	fs.StringVar(&c.Editor, "EDITOR", "editor", "Text editor to use")

	fs.BoolVar(&help, "help", false, "")
	fs.BoolVar(&h, "h", false, "")

	fs.Parse(args)

	if h || help {
		execHelp(&c)
		return
	}

	// Retrieve command name as first argument.
	cmd := fs.Arg(0)

	switch cmd {
	case "show":
		execShow(&c, args[1:])
	case "edit":
		execEdit(&c, args[1:])
	case "insert", "add":
		execInsert(&c, args[1:])
	case "find", "ls", "search":
		execFind(&c, args[1:])
	case "":
		execFind(&c, args)
	case "grep":
		execGrep(&c, args[1:])
	case "cp":
		execCp(&c, args[1:])
	case "mv":
		execMv(&c, args[1:])
	case "rm":
		execRm(&c, args[1:])
	case "generate":
		execGenerate(&c, args[1:])
	case "git":
		execGit(&c, args[1:])
	case "help":
		execHelp(&c)
	case "init":
		execInit(&c, args[1:])
	case "version":
		execVersion(&c)
	default:
		execShow(&c, args)
	}

}

func execHelp(c *commandLine) {
	fmt.Fprintln(c.WriterOutput, `Usage:
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

func execVersion(c *commandLine) {
	fmt.Fprintf(c.WriterOutput, "gopass v%s\n", version.Version)
}

func execInit(c *commandLine, args []string) {
	var path, p string

	fs := flag.NewFlagSet("init", flag.ContinueOnError)

	fs.StringVar(&path, "path", getDefaultPasswordStoreDir(c), "")
	fs.StringVar(&p, "p", "", "")

	fs.Usage = func() {
		fmt.Fprintln(c.WriterOutput, `Usage: gopass init [ --path=sub-folder, -p sub-folder ] gpg-id`)
	}
	err := fs.Parse(args)
	if err != nil {
		return
	}

	if p != "" {
		path = p
	}

	path, err = filepath.Abs(path)
	if err != nil {
		fmt.Fprintln(c.WriterOutput, err)
		return
	}

	if fs.NArg() != 1 {
		fs.Usage()
		return
	}

	gpgID := fs.Arg(0)

	store := gopass.NewPasswordStore(path)
	err = store.Init(gpgID)
	if err != nil {
		fmt.Fprintln(c.WriterOutput, err)
		return
	}

	fmt.Fprintf(c.WriterOutput, "Sucessfully created Password Store at %s\n", path)
}

//execInsert runs the "insert" command.
func execInsert(c *commandLine, args []string) {
	var multiline, m bool
	var force, f bool

	fs := flag.NewFlagSet("insert", flag.ContinueOnError)

	fs.BoolVar(&multiline, "multiline", false, "")
	fs.BoolVar(&m, "m", false, "")

	fs.BoolVar(&force, "force", false, "")
	fs.BoolVar(&f, "f", false, "")

	fs.Usage = func() {
		fmt.Fprintln(c.WriterOutput, `Usage: gopass insert [ --multiline, -m ] [ --force, -f ] pass-name`)
	}
	err := fs.Parse(args)
	if err != nil {
		return
	}

	multiline = multiline || m
	force = force || f

	pwname := fs.Arg(0)

	store := getStore(c)

	passwordPath := path.Join(store.Path, pwname+".gpg")

	// Check if password already exists
	if _, err := os.Stat(passwordPath); err == nil && !force {
		fmt.Fprintf(c.WriterOutput, "Error: Password already exists at '%s', use -f to force\n", passwordPath)
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
			fmt.Fprintln(c.WriterOutput, "Enter password:")
			try1, _ := terminal.ReadPassword(fd)
			fmt.Fprintln(c.WriterOutput)

			fmt.Fprintln(c.WriterOutput, "Enter confirmation:")
			try2, _ := terminal.ReadPassword(fd)
			fmt.Fprintln(c.WriterOutput)

			if string(try1) == string(try2) {
				password = string(try1)
				break
			} else {
				fmt.Fprintln(c.WriterOutput, "Passwords did not match, try again...")
			}

		}
	}

	err = store.InsertPassword(pwname, password)
	if err != nil {
		fmt.Fprintln(c.WriterOutput, err)
		return
	}

	fmt.Fprintf(c.WriterOutput, "Password %s added to the store\n", pwname)
}

//execEdit rund the "edit" command.
func execEdit(cmd *commandLine, args []string) {
	fs := flag.NewFlagSet("edit", flag.ExitOnError)
	fs.Parse(args)

	store := getStore(cmd)

	passname := fs.Arg(0)

	password, err := store.GetPassword(passname)

	if err != nil {
		fmt.Fprintf(cmd.WriterOutput, "Error: %s\n", err)
		os.Exit(1)
	}

	file, _ := ioutil.TempFile(os.TempDir(), "gopass")
	defer os.Remove(file.Name())

	ioutil.WriteFile(file.Name(), []byte(password), 0600)

	editor := exec.Command(cmd.Editor, file.Name())
	editor.Stdout = os.Stdout
	editor.Stdin = os.Stdin
	editor.Stderr = os.Stderr
	editor.Run()

	pwText, _ := ioutil.ReadFile(file.Name())
	password = string(pwText)

	err = store.InsertPassword(passname, password)
	if err != nil {
		fmt.Fprintln(cmd.WriterOutput, err)
		os.Exit(1)
	}

	fmt.Fprintf(cmd.WriterOutput, "Succesfully edited password '%s'\n", passname)
}

//execGenerate runs the "generate" command.
func execGenerate(cmd *commandLine, args []string) {
	var noSymbols, n bool
	var force, f bool

	fs := flag.NewFlagSet("generate", flag.ContinueOnError)
	fs.Usage = func() {
		fmt.Fprintln(cmd.WriterOutput, `Usage: gopass generate [--no-symbols,-n] [--force,-f] pass-name pass-length`)
	}

	fs.BoolVar(&noSymbols, "no-symbols", false, "")
	fs.BoolVar(&n, "n", false, "")

	fs.BoolVar(&force, "force", false, "")
	fs.BoolVar(&f, "f", false, "")

	err := fs.Parse(args)
	if err != nil {
		return
	}

	noSymbols = noSymbols || n
	force = force || f

	passName := fs.Arg(0)
	passLength, err := strconv.ParseInt(fs.Arg(1), 0, 64)
	if err != nil {
		fmt.Fprintf(cmd.WriterOutput, "Second argument must be an int, got '%s'\n", fs.Arg(1))
		os.Exit(1)
	}

	runes := append(pwgen.Alpha, pwgen.Num...)
	if !noSymbols {
		runes = append(runes, pwgen.Symbols...)
	}

	password := pwgen.RandSeq(int(passLength), runes)

	store := getStore(cmd)

	//TODO: Check if password already exists
	store.InsertPassword(passName, password)
	fmt.Fprintf(cmd.WriterOutput, "Password %s added to the store\n", passName)
}

//execRm runs the "rm" command.
func execRm(c *commandLine, args []string) {
	fs := flag.NewFlagSet("rm", flag.ExitOnError)
	fs.Usage = func() { fmt.Fprintln(c.WriterOutput, "Usage: gopass rm pass-name") }
	fs.Parse(args)

	store := getStore(c)

	err := store.Remove(fs.Arg(0))
	if err != nil {
		fmt.Fprintf(c.WriterOutput, "Error: %s\n", err)
	} else {
		fmt.Fprintln(c.WriterOutput, "Removed password/directory at path", fs.Arg(0))
	}
}

//execMv runs the "mv" comand.
func execMv(c *commandLine, args []string) {
	fs := flag.NewFlagSet("mv", flag.ExitOnError)
	fs.Usage = func() { fmt.Fprintln(c.WriterOutput, "Usage: gopass mv old-path new-path") }
	fs.Parse(args)

	store := getStore(c)

	source := fs.Arg(0)
	dest := fs.Arg(1)

	if source == "" || dest == "" {
		fmt.Fprintln(c.WriterOutput, "Error: Received empty source or dest argument")
		os.Exit(1)
	}

	err := store.Move(source, dest)
	if err != nil {
		fmt.Fprintf(c.WriterOutput, "Error: %s\n", err)
	} else {
		fmt.Fprintf(c.WriterOutput, "Moved password/directory from '%s' to '%s'\n", source, dest)
	}
}

//execCp runs the "cp" command.
func execCp(c *commandLine, args []string) {
	fs := flag.NewFlagSet("cp", flag.ExitOnError)
	fs.Usage = func() { fmt.Fprintln(c.WriterOutput, "Usage: gopass cp old-path new-path") }
	fs.Parse(args)

	store := getStore(c)

	source := fs.Arg(0)
	dest := fs.Arg(1)

	if source == "" || dest == "" {
		fmt.Fprintln(c.WriterOutput, "Error: Received empty source or dest argument")
		os.Exit(1)
	}

	err := store.Copy(source, dest)
	if err != nil {
		fmt.Fprintf(c.WriterOutput, "Error: %s\n", err)
		os.Exit(1)
	} else {
		fmt.Fprintf(c.WriterOutput, "Copied password/directory from '%s' to '%s'\n", source, dest)
	}
}

//execShow runs the "show" command.
func execShow(c *commandLine, args []string) {
	fs := flag.NewFlagSet("show", flag.ExitOnError)
	fs.Usage = func() { fmt.Fprintln(c.WriterOutput, `Usage: gopass show [pass-name]`) }
	fs.Parse(args)
	password := fs.Arg(0)

	store := getStore(c)

	password, err := store.GetPassword(password)

	if err != nil {
		fmt.Fprintf(c.WriterOutput, "Error: %s\n", err)
		os.Exit(1)
	}

	fmt.Fprintln(c.WriterOutput, password)
}

//execFind runs the "find" command.
func execFind(c *commandLine, args []string) {
	fs := flag.NewFlagSet("find", flag.ExitOnError)
	fs.Parse(args)

	store := getStore(c)

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
func execGrep(c *commandLine, args []string) {
	fs := flag.NewFlagSet("grep", flag.ExitOnError)
	fs.Parse(args)

	pattern, _ := regexp.CompilePOSIX(fs.Arg(0))

	store := getStore(c)

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
			fmt.Fprintf(c.WriterOutput, "%s:\n%s", ansi.Color(password, "cyan+b"), output)
		}
	}
}

//execGit runs the "git" command
func execGit(c *commandLine, args []string) {
	store := getStore(c)

	gitArgs := []string{
		"--git-dir=" + store.GitDir,
		"--work-tree=" + store.Path}

	gitArgs = append(gitArgs, args...)

	git := exec.Command(
		"git",
		gitArgs...)

	git.Stdout = os.Stdout
	git.Stderr = os.Stderr
	git.Stdin = os.Stdin
	git.Run()
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

//getStore finds and returns the PasswordStore
func getStore(c *commandLine) *gopass.PasswordStore {
	storePath := getDefaultPasswordStoreDir(c)
	s := gopass.NewPasswordStore(storePath)
	return s
}
