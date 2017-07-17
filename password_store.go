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

//Init creates a Password Store at the Path
func (store *PasswordStore) Init(gpgID string) error {
	//Check if the password path already exists
	fi, err := os.Stat(store.Path)
	if err == nil {
		//Path exists, but is it a directory?
		if fi.Mode().IsDir() == false {
			return fmt.Errorf(
				"Could not create password store. Path `%s` already exists and it is not a directory.",
				store.Path)
		}
	} else {
		//Error during os.Stat
		if os.IsNotExist(err) {
			//Path does not exist, create it
			err = os.Mkdir(store.Path, 0700)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	//Check if the .gpg-id file already exists.
	gpgIDFilePath := path.Join(store.Path, ".gpg-id")
	fi, err = os.Stat(gpgIDFilePath)
	if err == nil {
		//.gpg-id already exists
		return fmt.Errorf("There is already a .gpg-id file at %s. Stopping init.", gpgIDFilePath)
	}

	gpgIDFile, err := os.Create(path.Join(store.Path, ".gpg-id"))
	if err != nil {
		return err
	}
	defer gpgIDFile.Close()
	gpgIDFile.WriteString(gpgID + "\n")
	store.GPGID = gpgID

	err = store.git("init")
	if err != nil {
		return err
	}

	return nil
}

//InsertPassword inserts a new password or overwrites an existing one
func (store *PasswordStore) InsertPassword(pwname, pwtext string) error {
	containsPassword, passwordPath := store.ContainsPassword(pwname)

	//Check if password already exists
	var gitAction string
	if containsPassword {
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

//RemoveDirectory removes the directory at the given path
func (store *PasswordStore) RemoveDirectory(dirname string) error {
	containsDirectory, directoryPath := store.ContainsDirectory(dirname)
	if containsDirectory {
		os.RemoveAll(directoryPath)

		store.AddAndCommit(
			fmt.Sprintf("Removed directory '%s' from the store", dirname),
			directoryPath)

		return nil
	}
	return fmt.Errorf("Could not find directory at path %s", directoryPath)
}

//RemovePassword removes the password at the given path
func (store *PasswordStore) RemovePassword(pwname string) error {
	containsPassword, passwordPath := store.ContainsPassword(pwname)
	if containsPassword {
		os.Remove(passwordPath)

		store.AddAndCommit(
			fmt.Sprintf("Removed password '%s' from the store", pwname),
			passwordPath)

		return nil
	}
	return fmt.Errorf("Could not find password at path %s", passwordPath)
}

//Move moves a passsword or directory from source to dest.
func (store *PasswordStore) Move(source, dest string) error {

	//Check if the path is a dir
	containsDirectory, sourceDirectoryPath := store.ContainsDirectory(source)
	if containsDirectory {
		destDirectoryPath := path.Join(store.Path, dest)
		os.Rename(sourceDirectoryPath, destDirectoryPath)

		store.AddAndCommit(
			fmt.Sprintf("Moved directory '%s' to '%s'", source, dest),
			sourceDirectoryPath,
			destDirectoryPath)

		return nil
	}

	//Check if the path is a password
	containsPassword, sourcePasswordPath := store.ContainsPassword(source)
	if containsPassword {
		destPasswordPath := path.Join(store.Path, dest+".gpg")
		os.Rename(sourcePasswordPath, destPasswordPath)

		store.AddAndCommit(
			fmt.Sprintf("Moved Password '%s' to '%s'", source, dest),
			sourcePasswordPath,
			destPasswordPath)

		return nil
	}

	return fmt.Errorf("Could not find password or directory at path %s", path.Join(store.Path, source))
}

//Copy copies a password or directory from source to dest.
func (store *PasswordStore) Copy(source, dest string) error {

	//Check if the path is a dir
	containsDirectory, sourceDirectoryPath := store.ContainsDirectory(source)
	if containsDirectory {
		err := exec.Command("cp", "-r", sourceDirectoryPath, path.Join(store.Path, dest)).Run()
		return err
	}

	//Check if the path is a password
	containsPassword, sourcePasswordPath := store.ContainsPassword(source)
	if containsPassword {
		destPasswordPath := path.Join(store.Path, dest+".gpg")
		gopassio.CopyFileContents(sourcePasswordPath, destPasswordPath)

		store.AddAndCommit(
			fmt.Sprintf("Copied Password '%s' to '%s'", source, dest),
			destPasswordPath)

		return nil
	}

	return fmt.Errorf("Could not find password or directory at path %s", path.Join(store.Path, source))
}

//GetPassword returns a decrypted password
func (store *PasswordStore) GetPassword(pwname string) (string, error) {
	containsPassword, passwordPath := store.ContainsPassword(pwname)

	//Error if the password does not exist
	if containsPassword == false {
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

//ContainsPassword returns whether or not the store contains a password with this name.
//it also conveniently returns the password path that was checked
func (store *PasswordStore) ContainsPassword(pwname string) (bool, string) {
	passwordPath := path.Join(store.Path, pwname+".gpg")

	if _, err := os.Stat(passwordPath); os.IsNotExist(err) {
		return false, passwordPath
	}

	return true, passwordPath
}

//ContainsDirectory returns whether or not the store contains a directory with this name.
//it also conveniently returns the directory path that was checked
func (store *PasswordStore) ContainsDirectory(dirname string) (bool, string) {
	directoryPath := path.Join(store.Path, dirname)

	if _, err := os.Stat(directoryPath); os.IsNotExist(err) {
		return false, directoryPath
	}

	return true, directoryPath
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

	store.git("commit", "-m", message)

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
