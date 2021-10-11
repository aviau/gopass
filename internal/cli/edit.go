//    Copyright (C) 2018 Alexandre Viau <alexandre@alexandreviau.net>
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
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
)

// execEdit runs the "edit" command.
func execEdit(cfg CommandConfig, args []string) error {
	var help, h bool

	fs := flag.NewFlagSet("edit", flag.ContinueOnError)
	fs.SetOutput(ioutil.Discard)

	fs.Usage = func() { fmt.Fprintln(cfg.WriterOutput(), "Usage: gopass edit pass-name") }

	fs.BoolVar(&help, "help", false, "")
	fs.BoolVar(&h, "h", false, "")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if help || h {
		fs.Usage()
		return nil
	}

	store := cfg.PasswordStore()

	passName := fs.Arg(0)

	action := "inserted"
	decryptedPassword := ""
	if containsPasword, _ := store.ContainsPassword(passName); containsPasword {
		var err error
		decryptedPassword, err = store.GetPassword(passName)
		if err != nil {
			return err
		}
		action = "edited"
	}

	editedPassword, err := editUsingTempfile(cfg, decryptedPassword)
	if err != nil {
		return fmt.Errorf("could not edit password: %w", err)
	}

	if err := store.InsertPassword(passName, editedPassword); err != nil {
		return err
	}

	fmt.Fprintf(cfg.WriterOutput(), "Succesfully %s password \"%s\".\n", action, passName)
	return nil
}

func editUsingTempfile(cfg CommandConfig, pass string) (string, error) {
	file, err := ioutil.TempFile(os.TempDir(), "gopass")
	if err != nil {
		return "", fmt.Errorf("can't create tempfile: %w", err)
	}
	defer file.Close()
	defer os.Remove(file.Name())

	if _, err := file.WriteString(pass); err != nil {
		return "", fmt.Errorf("could not write password to tempfile: %w", err)
	}

	editor := exec.Command(cfg.Editor(), file.Name())
	editor.Stdout = cfg.WriterOutput()
	editor.Stderr = cfg.WriterError()
	editor.Stdin = cfg.ReaderInput()
	editor.Run()

	if _, err := file.Seek(0, 0); err != nil {
		return "", fmt.Errorf("could not seek to password start: %w", err)
	}

	editedPasswordBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("could not read edited file: %w", err)
	}

	editedPassword := string(editedPasswordBytes)

	return editedPassword, nil
}
