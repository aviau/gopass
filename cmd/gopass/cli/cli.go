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
	gopass_terminal "github.com/aviau/gopass/terminal"
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

func execHelp(c *commandLine) error {
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
	return nil
}

func execVersion(c *commandLine) error {
	fmt.Fprintf(c.WriterOutput, "gopass v%s\n", version.Version)
	return nil
}

func execInit(c *commandLine, args []string) error {
	var path, p string

	fs := flag.NewFlagSet("init", flag.ContinueOnError)

	fs.StringVar(&path, "path", getDefaultPasswordStoreDir(c), "")
	fs.StringVar(&p, "p", "", "")

	fs.Usage = func() {
		fmt.Fprintln(c.WriterOutput, `Usage: gopass init [ --path=sub-folder, -p sub-folder ] gpg-id`)
	}
	err := fs.Parse(args)
	if err != nil {
		return err
	}

	if p != "" {
		path = p
	}

	path, err = filepath.Abs(path)
	if err != nil {
		return err
	}

	if fs.NArg() != 1 {
		fs.Usage()
		return nil
	}

	gpgID := fs.Arg(0)

	store := gopass.NewPasswordStore(path)
	err = store.Init(gpgID)
	if err != nil {
		return err
	}

	fmt.Fprintf(c.WriterOutput, "Sucessfully created Password Store at %s\n", path)
	return nil
}

//execInsert runs the "insert" command.
func execInsert(c *commandLine, args []string) error {
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
		return err
	}

	multiline = multiline || m
	force = force || f

	pwname := fs.Arg(0)

	store := getStore(c)

	passwordPath := path.Join(store.Path, pwname+".gpg")

	// Check if password already exists
	if _, err := os.Stat(passwordPath); err == nil && !force {
		fmt.Fprintf(c.WriterOutput, "Error: Password already exists at '%s', use -f to force\n", passwordPath)
		return nil
	}

	var password string

	if multiline {
		file, _ := ioutil.TempFile(os.TempDir(), "gopass")
		defer os.Remove(file.Name())

		editor := exec.Command(getEditor(c), file.Name())
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
		return err
	}

	fmt.Fprintf(c.WriterOutput, "Password %s added to the store\n", pwname)
	return nil
}

//execEdit rund the "edit" command.
func execEdit(cmd *commandLine, args []string) error {
	fs := flag.NewFlagSet("edit", flag.ExitOnError)
	fs.Parse(args)

	store := getStore(cmd)

	passname := fs.Arg(0)

	password, err := store.GetPassword(passname)

	if err != nil {
		return err
	}

	file, _ := ioutil.TempFile(os.TempDir(), "gopass")
	defer os.Remove(file.Name())

	ioutil.WriteFile(file.Name(), []byte(password), 0600)

	editor := exec.Command(getEditor(cmd), file.Name())
	editor.Stdout = os.Stdout
	editor.Stdin = os.Stdin
	editor.Stderr = os.Stderr
	editor.Run()

	pwText, _ := ioutil.ReadFile(file.Name())
	password = string(pwText)

	err = store.InsertPassword(passname, password)
	if err != nil {
		return err
	}

	fmt.Fprintf(cmd.WriterOutput, "Succesfully edited password '%s'\n", passname)
	return nil
}

//execGenerate runs the "generate" command.
func execGenerate(cmd *commandLine, args []string) error {
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
		return err
	}

	noSymbols = noSymbols || n
	force = force || f

	passName := fs.Arg(0)

	store := getStore(cmd)

	if containsPassword, _ := store.ContainsPassword(passName); containsPassword {
		if !force {
			fmt.Fprintf(cmd.WriterOutput, "Error: '%s' already exists. Use -f to override.\n", passName)
			return nil
		}
	}

	passLength, err := strconv.ParseInt(fs.Arg(1), 0, 64)
	if err != nil {
		fmt.Fprintf(cmd.WriterOutput, "Second argument must be an int, got '%s'\n", fs.Arg(1))
		return err
	}

	runes := append(pwgen.Alpha, pwgen.Num...)
	if !noSymbols {
		runes = append(runes, pwgen.Symbols...)
	}

	password := pwgen.RandSeq(int(passLength), runes)

	err = store.InsertPassword(passName, password)
	if err != nil {
		return err
	}

	fmt.Fprintf(cmd.WriterOutput, "Password %s added to the store\n", passName)
	return nil
}

//execRm runs the "rm" command.
func execRm(c *commandLine, args []string) error {
	var recursive, r bool
	var force, f bool

	fs := flag.NewFlagSet("rm", flag.ContinueOnError)
	fs.Usage = func() {
		fmt.Fprintln(c.WriterOutput, "Usage: gopass rm pass-name")
	}

	fs.BoolVar(&recursive, "recursive", false, "")
	fs.BoolVar(&r, "r", false, "")

	fs.BoolVar(&force, "force", false, "")
	fs.BoolVar(&f, "f", false, "")

	err := fs.Parse(args)
	if err != nil {
		return err
	}

	force = force || f
	recursive = recursive || r

	store := getStore(c)

	pwname := fs.Arg(0)
	if pwname == "" {
		fs.Usage()
		return nil
	}

	if containsPassword, _ := store.ContainsPassword(pwname); containsPassword {

		if !force {
			if !gopass_terminal.AskYesNo(c.WriterOutput, fmt.Sprintf("Are you sure you would like to delete %s? [y/n] ", pwname)) {
				return nil
			}
		}

		err = store.RemovePassword(pwname)
		if err != nil {
			return err
		}

	} else if containsDirectory, _ := store.ContainsDirectory(pwname); containsDirectory {

		if !recursive {
			fmt.Fprintf(c.WriterOutput, "Error: %s is a directory, use -r to remove recursively\n", pwname)
			return nil
		}

		if !force {
			if !gopass_terminal.AskYesNo(c.WriterOutput, fmt.Sprintf("Are you sure you would like to delete %s recursively? [y/n] ", pwname)) {
				return nil
			}
		}

		err = store.RemoveDirectory(pwname)
		if err != nil {
			return err
		}
	}

	fmt.Fprintln(c.WriterOutput, "Removed password/directory at path", fs.Arg(0))
	return nil
}

//execMv runs the "mv" comand.
func execMv(c *commandLine, args []string) error {
	var force, f bool

	fs := flag.NewFlagSet("mv", flag.ExitOnError)
	fs.Usage = func() { fmt.Fprintln(c.WriterOutput, "Usage: gopass mv old-path new-path") }

	fs.BoolVar(&force, "force", false, "")
	fs.BoolVar(&f, "f", false, "")

	err := fs.Parse(args)
	if err != nil {
		return err
	}

	force = force || f

	store := getStore(c)

	source := fs.Arg(0)
	dest := fs.Arg(1)

	if source == "" || dest == "" {
		fmt.Fprintln(c.WriterOutput, "Error: Received empty source or dest argument")
		return nil
	}

	// If the dest ends with a '/', then it is a directory.
	if strings.HasSuffix(dest, "/") {
		_, sourceFile := filepath.Split(source)
		dest = filepath.Join(dest, sourceFile)
	}

	if sourceIsPassword, _ := store.ContainsPassword(source); sourceIsPassword {

		if destAlreadyExists, _ := store.ContainsPassword(dest); destAlreadyExists {
			if !force {
				fmt.Fprintf(c.WriterOutput, "Error: destination %s already exists. Use -f to override\n", dest)
				return nil
			}
		}

		err = store.MovePassword(source, dest)
		if err != nil {
			return err
		}

		fmt.Fprintf(c.WriterOutput, "Moved password from '%s' to '%s'\n", source, dest)
		return nil
	}

	if sourceIsDirectory, _ := store.ContainsDirectory(source); sourceIsDirectory {
		err = store.MoveDirectory(source, dest)
		if err != nil {
			return err
		}
		fmt.Fprintf(c.WriterOutput, "Moved directory from '%s' to '%s'\n", source, dest)
		return nil
	}

	fmt.Fprintf(c.WriterOutput, "Error: could not find source '%s' to copy \n", source)
	return nil
}

//execCp runs the "cp" command.
func execCp(c *commandLine, args []string) error {
	var recursive, r bool
	var force, f bool

	fs := flag.NewFlagSet("cp", flag.ExitOnError)
	fs.Usage = func() { fmt.Fprintln(c.WriterOutput, "Usage: gopass cp old-path new-path") }

	fs.BoolVar(&recursive, "recursive", false, "")
	fs.BoolVar(&r, "r", false, "")

	fs.BoolVar(&force, "force", false, "")
	fs.BoolVar(&f, "f", false, "")

	err := fs.Parse(args)
	if err != nil {
		return err
	}

	recursive = recursive || r
	force = force || f

	store := getStore(c)

	source := fs.Arg(0)
	dest := fs.Arg(1)

	if source == "" || dest == "" {
		fmt.Fprintln(c.WriterOutput, "Error: Received empty source or dest argument")
		return nil
	}

	// If the dest ends with a '/', then it is a directory.
	if strings.HasSuffix(dest, "/") {
		_, sourceFile := filepath.Split(source)
		dest = filepath.Join(dest, sourceFile)
	}

	if sourceIsPassword, _ := store.ContainsPassword(source); sourceIsPassword {

		if destAlreadyExists, _ := store.ContainsPassword(dest); destAlreadyExists {
			if !force {
				fmt.Fprintf(c.WriterOutput, "Error: destination %s already exists. Use -f to override\n", dest)
				return nil
			}
		}

		err = store.CopyPassword(source, dest)
		if err != nil {
			return err
		}

		fmt.Fprintf(c.WriterOutput, "Copied password from '%s' to '%s'\n", source, dest)
		return nil
	}

	if sourceIsDirectory, _ := store.ContainsDirectory(source); sourceIsDirectory {

		if !recursive {
			fmt.Fprintf(c.WriterOutput, "Error: %s is a directory, use -r to copy recursively\n", source)
			return nil
		}

		err = store.CopyDirectory(source, dest)
		if err != nil {
			return err
		}

		fmt.Fprintf(c.WriterOutput, "Copied directory from '%s' to '%s'\n", source, dest)
		return nil
	}

	fmt.Fprintf(c.WriterOutput, "Error: could not find source '%s' to copy \n", source)
	return nil
}

//execShow runs the "show" command.
func execShow(c *commandLine, args []string) error {
	fs := flag.NewFlagSet("show", flag.ExitOnError)
	fs.Usage = func() { fmt.Fprintln(c.WriterOutput, `Usage: gopass show [pass-name]`) }
	fs.Parse(args)
	password := fs.Arg(0)

	store := getStore(c)

	password, err := store.GetPassword(password)

	if err != nil {
		return err
	}

	fmt.Fprintln(c.WriterOutput, password)
	return nil
}

//execFind runs the "find" command.
func execFind(c *commandLine, args []string) error {
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
	return nil
}

//execGrep runs the "grep" command
func execGrep(c *commandLine, args []string) error {
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
	return nil
}

//execGit runs the "git" command
func execGit(c *commandLine, args []string) error {
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
	return nil
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
