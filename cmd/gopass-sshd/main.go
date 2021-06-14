package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/gliderlabs/ssh"
)

var welcomeMessage = `
Welcome to the gopass-sshd server.

Please enter a command.
`

func main() {
	addr := "localhost:2222"

	fmt.Printf("Starting gopass-sshd on port %s...\n", addr)

	server := ssh.Server{
		Addr: addr,
		PublicKeyHandler: ssh.PublicKeyHandler(func(ctx ssh.Context, key ssh.PublicKey) bool {
			// TODO: What keys should we accept?
			return true
		}),
		Handler: ssh.Handler(func(s ssh.Session) {
			cmd := s.Command()

			if len(cmd) == 0 {
				fmt.Fprintln(s, "No command specified.")
				return
			}

			fmt.Fprintf(s, "%s\n", strings.Join(cmd, ""))
		}),
	}

	log.Fatal(server.ListenAndServe())
}
