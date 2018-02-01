package clipboard

import (
	"os/exec"
)

//CopyToClipboard copies a string to the clipboard using xclip
func CopyToClipboard(s string) error {
	xclip := exec.Command(
		"xclip",
		"-in",
		"-selection",
		"clipboard",
	)

	stdin, err := xclip.StdinPipe()
	if err != nil {
		return err
	}

	if err := xclip.Start(); err != nil {
		return err
	}

	_, err = stdin.Write([]byte(s))
	if err != nil {
		return err
	}

	if err := stdin.Close(); err != nil {
		return err
	}

	err = xclip.Wait()

	return err
}
