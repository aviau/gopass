//    Copyright (C) 2017 Alexandre Viau <alexandre@alexandreviau.net>
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

// Package store implements the Password Store.
package store

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/aviau/gopass/internal/gpg"
	gopassio "github.com/aviau/gopass/internal/io"
)

// PasswordStore represents a password store.
type PasswordStore struct {
	Path       string     // path of the store
	GitDir     string     // The path of the git directory
	GPGIDs     []string   // The GPG IDs used for the store
	GPGBackend GPGBackend // The store's GPG backend.
	UsesGit    bool       // Whether or not the store uses git
}

// GPGBackend the PasswordStore's GPG backend.
type GPGBackend interface {
	Encrypt(content []byte, recipients []string) ([]byte, error)
	Decrypt(content []byte) ([]byte, error)
}

// Returns the GPG ids for a given directory
func loadGPGIDs(directory string) ([]string, error) {
	gpgIDPath := path.Join(directory, ".gpg-id")

	file, err := os.Open(gpgIDPath)
	if err != nil {
		return nil, err
	}

	var gpgIDs []string

	fscanner := bufio.NewScanner(file)
	for fscanner.Scan() {
		gpgID := strings.TrimSpace(fscanner.Text())
		gpgIDs = append(gpgIDs, gpgID)
	}

	return gpgIDs, nil
}

// NewPasswordStore returns a new password store.
func NewPasswordStore(storePath string) *PasswordStore {
	s := PasswordStore{}
	s.Path = storePath
	s.UsesGit = true
	s.GitDir = path.Join(s.Path, ".git")
	s.GPGBackend = gpg.New("", nil, false)

	//Read the .gpg-id file
	gpgIDs, _ := loadGPGIDs(s.Path)
	s.GPGIDs = gpgIDs

	return &s
}

// Init creates a Password Store at the Path
func (store *PasswordStore) Init(gpgIDs []string) error {
	// Check if the password path already exists
	if fi, err := os.Stat(store.Path); os.IsNotExist(err) {
		// Path does not exist, create it
		if err := os.Mkdir(store.Path, 0700); err != nil {
			return err
		}
	} else if err != nil {
		return err
	} else {
		// Path exists, but is it a directory?
		if !fi.Mode().IsDir() {
			return fmt.Errorf(
				"could not create password store. Path \"%s\" already exists and it is not a directory",
				store.Path)
		}
	}

	// Check if the .gpg-id file already exists.
	gpgIDFilePath := path.Join(store.Path, ".gpg-id")
	if _, err := os.Stat(gpgIDFilePath); err == nil {
		return fmt.Errorf("there is already a .gpg-id file at \"%s\". Stopping init", gpgIDFilePath)
	} else if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("could not look for an existing .gpg-id: %w", err)
	}

	gpgIDFile, err := os.Create(path.Join(store.Path, ".gpg-id"))
	if err != nil {
		return err
	}
	defer gpgIDFile.Close()
	for _, gpgID := range gpgIDs {
		gpgIDFile.WriteString(gpgID + "\n")
	}
	store.GPGIDs = gpgIDs

	if err := store.git("init"); err != nil {
		return err
	}

	return store.AddAndCommit("initial commit", ".gpg-id")
}

// SetGPGIDs will set the store's GPG ids
func (store *PasswordStore) SetGPGIDs(gpgIDs []string) error {
	gpgIDFile, err := os.OpenFile(
		path.Join(store.Path, ".gpg-id"),
		os.O_WRONLY|os.O_CREATE|os.O_TRUNC,
		0644,
	)
	if err != nil {
		return err
	}
	defer gpgIDFile.Close()

	for _, gpgID := range gpgIDs {
		gpgIDFile.WriteString(gpgID + "\n")
	}
	store.GPGIDs = gpgIDs

	return store.AddAndCommit(
		fmt.Sprintf("Set GPG id to %s", strings.Join(gpgIDs, ", ")),
		".gpg-id",
	)
}

// ReencryptPassword will reencrypt a password to the current GPG ids
func (store *PasswordStore) ReencryptPassword(pwname string) error {
	containsPassword, passwordPath := store.ContainsPassword(pwname)

	// Error if the password does not exist
	if !containsPassword {
		return fmt.Errorf("could not find password \"%s\" at path \"%s\"", pwname, passwordPath)
	}

	encryptedPassword, err := ioutil.ReadFile(passwordPath)
	if err != nil {
		return fmt.Errorf("could not read encrypted password: %w", err)
	}

	decryptedPassword, err := store.GPGBackend.Decrypt(encryptedPassword)
	if err != nil {
		return fmt.Errorf("could not decrypt the password: %w", err)
	}

	reEncryptedPassword, err := store.GPGBackend.Encrypt(decryptedPassword, store.GPGIDs)
	if err != nil {
		return fmt.Errorf("could not re-encrypt the password: %w", err)
	}

	if err := ioutil.WriteFile(passwordPath, reEncryptedPassword, 0600); err != nil {
		return fmt.Errorf("could not write the newly encrypted password: %w", err)
	}

	return nil
}

// InsertPassword inserts a new password or overwrites an existing one
func (store *PasswordStore) InsertPassword(pwname, pwtext string) error {
	containsPassword, passwordPath := store.ContainsPassword(pwname)

	// Check if password already exists
	var gitAction string
	if containsPassword {
		gitAction = "edited"
	} else {
		gitAction = "added"
	}

	encryptedPassword, err := store.GPGBackend.Encrypt([]byte(pwtext), store.GPGIDs)
	if err != nil {
		return fmt.Errorf("could not encrypt the password: %w", err)
	}

	if err := ioutil.WriteFile(passwordPath, encryptedPassword, 0600); err != nil {
		return fmt.Errorf("could not write the newly encrypted password: %w", err)
	}

	store.AddAndCommit(
		fmt.Sprintf("%s password \"%s\"", gitAction, pwname),
		passwordPath)

	return nil
}

// RemoveDirectory removes the directory at the given path
func (store *PasswordStore) RemoveDirectory(dirname string) error {
	containsDirectory, directoryPath := store.ContainsDirectory(dirname)

	if !containsDirectory {
		return fmt.Errorf("could not find directory at path \"%s\"", directoryPath)
	}

	if err := os.RemoveAll(directoryPath); err != nil {
		return err
	}

	store.AddAndCommit(
		fmt.Sprintf("removed directory \"%s\" from the store", dirname),
		directoryPath)

	return nil
}

// RemovePassword removes the password at the given path
func (store *PasswordStore) RemovePassword(pwname string) error {
	containsPassword, passwordPath := store.ContainsPassword(pwname)

	if !containsPassword {
		return fmt.Errorf("could not find password at path \"%s\"", passwordPath)
	}

	os.Remove(passwordPath)

	store.AddAndCommit(
		fmt.Sprintf("removed password \"%s\" from the store", pwname),
		passwordPath)

	return nil
}

// MoveDirectory moves a directory from source to dest
func (store *PasswordStore) MoveDirectory(source, dest string) error {
	containsDirectory, sourceDirectoryPath := store.ContainsDirectory(source)
	if !containsDirectory {
		return fmt.Errorf("could not find directory at path \"%s\"", sourceDirectoryPath)
	}

	destDirectoryPath := path.Join(store.Path, dest)

	if err := os.Rename(sourceDirectoryPath, destDirectoryPath); err != nil {
		return err
	}

	store.AddAndCommit(
		fmt.Sprintf("moved directory \"%s\" to \"%s\"", source, dest),
		sourceDirectoryPath,
		destDirectoryPath)

	return nil
}

// MovePassword moves a passsword or directory from source to dest.
func (store *PasswordStore) MovePassword(source, dest string) error {
	containsPassword, sourcePasswordPath := store.ContainsPassword(source)

	if !containsPassword {
		return fmt.Errorf("could not find password path \"%s\"", sourcePasswordPath)
	}

	// If the dest ends with a '/', then it is a directory.
	var destPasswordPath string
	if strings.HasSuffix(dest, "/") {
		_, file := filepath.Split(sourcePasswordPath)
		destPasswordPath = path.Join(store.Path, dest, file)
	} else {
		destPasswordPath = path.Join(store.Path, dest+".gpg")
	}

	if err := os.Rename(sourcePasswordPath, destPasswordPath); err != nil {
		return err
	}

	store.AddAndCommit(
		fmt.Sprintf("moved Password \"%s\" to \"%s\"", source, dest),
		sourcePasswordPath,
		destPasswordPath)

	return nil
}

// CopyPassword copies a password from source to dest
func (store *PasswordStore) CopyPassword(source, dest string) error {
	containsPassword, sourcePasswordPath := store.ContainsPassword(source)

	if !containsPassword {
		return fmt.Errorf("could not find password or at path \"%s\"", sourcePasswordPath)
	}

	// If the dest ends with a '/', then it is a directory.
	var destPasswordPath string
	if strings.HasSuffix(dest, "/") {
		_, file := filepath.Split(sourcePasswordPath)
		destPasswordPath = path.Join(store.Path, dest, file)
	} else {
		destPasswordPath = path.Join(store.Path, dest+".gpg")
	}

	if err := gopassio.CopyFileContents(sourcePasswordPath, destPasswordPath); err != nil {
		return err
	}

	store.AddAndCommit(
		fmt.Sprintf("copied Password \"%s\" to \"%s\"", source, dest),
		destPasswordPath)

	return nil
}

// CopyDirectory copies a directory from source to dest
func (store *PasswordStore) CopyDirectory(source, dest string) error {
	containsDirectory, sourceDirectoryPath := store.ContainsDirectory(source)

	if !containsDirectory {
		return fmt.Errorf("could not find directory at path \"%s\"", path.Join(store.Path, source))
	}

	destDirectoryPath := path.Join(store.Path, dest)
	if err := exec.Command("cp", "-r", sourceDirectoryPath, destDirectoryPath).Run(); err != nil {
		return err
	}

	store.AddAndCommit(
		fmt.Sprintf("copied directory \"%s\" to \"%s\"", source, dest),
		destDirectoryPath)

	return nil
}

// GetPassword returns a decrypted password
func (store *PasswordStore) GetPassword(pwname string) (string, error) {
	containsPassword, passwordPath := store.ContainsPassword(pwname)

	// Error if the password does not exist
	if !containsPassword {
		return "", fmt.Errorf("could not find password \"%s\" at path \"%s\"", pwname, passwordPath)
	}

	encryptedPassword, err := ioutil.ReadFile(passwordPath)
	if err != nil {
		return "", fmt.Errorf("could not read the encrypted password: %w", err)
	}

	decryptedPassword, err := store.GPGBackend.Decrypt(encryptedPassword)
	if err != nil {
		return "", fmt.Errorf("could not decrypt the password: %w", err)
	}

	return strings.TrimSpace(string(decryptedPassword)), nil
}

// ContainsPassword returns whether or not the store contains a password with this name.
// it also conveniently returns the password path that was checked
func (store *PasswordStore) ContainsPassword(pwname string) (bool, string) {
	passwordPath := path.Join(store.Path, pwname+".gpg")

	if _, err := os.Stat(passwordPath); os.IsNotExist(err) {
		return false, passwordPath
	}

	return true, passwordPath
}

// ContainsDirectory returns whether or not the store contains a directory with this name.
// it also conveniently returns the directory path that was checked
func (store *PasswordStore) ContainsDirectory(dirname string) (bool, string) {
	directoryPath := path.Join(store.Path, dirname)

	if _, err := os.Stat(directoryPath); os.IsNotExist(err) {
		return false, directoryPath
	}

	return true, directoryPath
}

// GetPasswordsList returns a list of all the passwords
func (store *PasswordStore) GetPasswordsList() []string {
	var list []string

	var scan = func(path string, fileInfo os.FileInfo, inpErr error) (err error) {
		if strings.HasSuffix(path, ".gpg") {
			password := strings.TrimSuffix(
				strings.TrimPrefix(path, store.Path+"/"),
				".gpg",
			)
			list = append(list, password)
		}
		return
	}

	filepath.Walk(store.Path, scan)

	return list
}

// AddAndCommit adds paths to the index and creates a commit
func (store *PasswordStore) AddAndCommit(message string, paths ...string) error {
	store.git("reset")

	for _, path := range paths {
		store.git("add", path)
	}

	store.git("commit", "-m", message)

	return nil
}

// git executes a git command
func (store *PasswordStore) git(args ...string) error {
	if !store.UsesGit {
		return nil
	}

	gitArgs := []string{
		"--git-dir=" + store.GitDir,
		"--work-tree=" + store.Path}

	gitArgs = append(gitArgs, args...)

	git := exec.Command("git", gitArgs...)

	// Should we do that?
	git.Stdout = os.Stdout
	git.Stderr = os.Stderr
	git.Stdin = os.Stdin

	if err := git.Run(); err != nil {
		fmt.Println(err.Error())
		return fmt.Errorf("git error: \"%s\"", err.Error())
	}

	return nil
}
