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

package gopass

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	gopassio "github.com/aviau/gopass/io"
)

//PasswordStore represents a password store.
type PasswordStore struct {
	Path   string //path of the store
	GitDir string //The path of the git directory
	GPGBin string //The GPG binary to use
	GPGID  string //The GPG ID used to encrypt the passwords
}

// NewPasswordStore returns a new password store.
func NewPasswordStore(storePath string) *PasswordStore {
	s := PasswordStore{}
	s.Path = storePath
	s.GitDir = path.Join(s.Path, ".git")

	// Find the GPG bin
	which := exec.Command("which", "gpg2")
	err := which.Run()
	if err == nil {
		s.GPGBin = "gpg2"
	} else {
		s.GPGBin = "gpg"
	}

	//Read the .gpg-id file
	gpgIDPath := path.Join(s.Path, ".gpg-id")
	content, _ := ioutil.ReadFile(gpgIDPath)
	s.GPGID = strings.TrimSpace(string(content))

	return &s
}

//InsertPassword inserts a new password or overwrites an existing one
func (store *PasswordStore) InsertPassword(pwname, pwtext string) error {
	passwordPath := path.Join(store.Path, pwname+".gpg")

	//Check if password already exists
	var gitAction string
	if _, err := os.Stat(passwordPath); err == nil {
		gitAction = "Edited"
	} else {
		gitAction = "Added"
	}

	gpg := exec.Command(
		store.GPGBin,
		"-e",
		"-r", store.GPGID,
		"--batch",
		"--use-agent",
		"--no-tty",
		"--yes",
		"-o", passwordPath)

	stdin, _ := gpg.StdinPipe()
	io.WriteString(stdin, pwtext)
	stdin.Close()
	output, err := gpg.CombinedOutput()

	if err != nil {
		return fmt.Errorf("Error: %s" + string(output))
	}

	store.AddAndCommit(
		fmt.Sprintf("%s password '%s'", gitAction, pwname),
		passwordPath)

	return nil
}

//Remove removes a password or a directory of the store
func (store *PasswordStore) Remove(pwname string) error {
	passwordPath := path.Join(store.Path, pwname)

	//Check if the path is a dir
	if _, err := os.Stat(passwordPath); err == nil {
		os.RemoveAll(passwordPath)

		store.AddAndCommit(
			fmt.Sprintf("Removed directory '%s' from the store", pwname),
			passwordPath)

		return nil
	}

	//Check if the path is a password
	passwordPath += ".gpg"
	if _, err := os.Stat(passwordPath); err == nil {
		os.Remove(passwordPath)

		store.AddAndCommit(
			fmt.Sprintf("Removed password '%s' from the store", pwname),
			passwordPath)

		return nil
	}

	return fmt.Errorf("Could not find password or directory at path %s", path.Join(store.Path, pwname))
}

//Move moves a passsword or directory from source to dest.
func (store *PasswordStore) Move(source, dest string) error {
	passwordPath := path.Join(store.Path, source)

	//Check if the path is a dir
	if _, err := os.Stat(passwordPath); err == nil {
		os.Rename(passwordPath, path.Join(store.Path, dest))

		store.AddAndCommit(
			fmt.Sprintf("Moved directory '%s' to '%s'", source, dest),
			passwordPath)

		return nil
	}

	//Check if the path is a password
	passwordPath += ".gpg"
	if _, err := os.Stat(passwordPath); err == nil {
		os.Rename(passwordPath, path.Join(store.Path, dest+".gpg"))

		store.AddAndCommit(
			fmt.Sprintf("Moved Password '%s' to '%s'", source, dest),
			passwordPath)

		return nil
	}

	return fmt.Errorf("Could not find password or directory at path %s", path.Join(store.Path, source))
}

//Copy copies a password or directory from source to dest.
func (store *PasswordStore) Copy(source, dest string) error {
	passwordPath := path.Join(store.Path, source)

	//Check if the path is a dir
	if _, err := os.Stat(passwordPath); err == nil {
		err := exec.Command("cp", "-r", passwordPath, path.Join(store.Path, dest)).Run()
		return err
	}

	//Check if the path is a password
	passwordPath += ".gpg"
	if _, err := os.Stat(passwordPath); err == nil {
		gopassio.CopyFileContents(passwordPath, path.Join(store.Path, dest+".gpg"))

		store.AddAndCommit(
			fmt.Sprintf("Copied Password '%s' to '%s'", source, dest),
			passwordPath)

		return nil
	}

	return fmt.Errorf("Could not find password or directory at path %s", path.Join(store.Path, source))
}

//GetPassword returns a decrypted password
func (store *PasswordStore) GetPassword(pwname string) (string, error) {
	passwordPath := path.Join(store.Path, pwname+".gpg")

	// Check if the passwiord exists
	if _, err := os.Stat(passwordPath); os.IsNotExist(err) {
		return "", fmt.Errorf("Could not find password '%s' at path '%s'", pwname, passwordPath)
	}

	//TODO: Use GPG lib instead
	show := exec.Command(store.GPGBin, "--quiet", "--batch", "--use-agent", "-d", passwordPath)
	output, err := show.CombinedOutput()

	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(output)), nil
}

//GetPasswordsList returns a list of all the passwords
func (store *PasswordStore) GetPasswordsList() []string {
	var list []string

	var scan = func(path string, fileInfo os.FileInfo, inpErr error) (err error) {
		if strings.HasSuffix(path, ".gpg") {
			_, file := filepath.Split(path)
			password := strings.TrimSuffix(file, ".gpg")
			list = append(list, password)
		}
		return
	}

	filepath.Walk(store.Path, scan)

	return list
}

//AddAndCommit adds paths to the index and creates a commit
func (store *PasswordStore) AddAndCommit(message string, paths ...string) error {
	store.git("reset")

	for _, path := range paths {
		store.git("add", path)
	}

	store.git("commit", "-m", fmt.Sprintf("'%s'", message))

	return nil
}

//git executes a git command
func (store *PasswordStore) git(args ...string) error {
	gitArgs := []string{
		"--git-dir=" + store.GitDir,
		"--work-tree=" + store.Path}

	gitArgs = append(gitArgs, args...)

	git := exec.Command("git", gitArgs...)

	//Should we do that?
	git.Stdout = os.Stdout
	git.Stderr = os.Stderr
	git.Stdin = os.Stdin

	err := git.Run()

	if err != nil {
		fmt.Println(err.Error())
		return fmt.Errorf("Git error: %s", err.Error())
	}

	return nil
}
