//    Copyright (C) 2021 Alexandre Viau <alexandre@alexandreviau.net>
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
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
)

func editUsingTempfile(cfg CommandConfig, content string) (string, error) {
	file, err := ioutil.TempFile(os.TempDir(), "gopass")
	if err != nil {
		return "", fmt.Errorf("can't create tempfile: %w", err)
	}
	defer file.Close()
	defer os.Remove(file.Name())

	if _, err := file.WriteString(content); err != nil {
		return "", fmt.Errorf("could not write content to tempfile: %w", err)
	}

	editor := exec.Command(cfg.Editor(), file.Name())
	editor.Stdout = cfg.WriterOutput()
	editor.Stderr = cfg.WriterError()
	editor.Stdin = cfg.ReaderInput()
	editor.Run()

	if _, err := file.Seek(0, 0); err != nil {
		return "", fmt.Errorf("could not seek to file start: %w", err)
	}

	editedContentBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("could not read edited file: %w", err)
	}

	editedContent := string(editedContentBytes)

	return editedContent, nil
}
