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

package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/aviau/gopass/pkg/cli"
)

func main() {
	// Create a context for the program. Cancel it on SIGINT.
	ctx := func() context.Context {
		ctx, cancel := context.WithCancel(context.Background())

		signalCh := make(chan os.Signal, 1)
		signal.Notify(signalCh, os.Interrupt)

		go func() {
			defer signal.Stop(signalCh)

			select {
			case <-signalCh:
				cancel()
			case <-ctx.Done():
			}
		}()

		return ctx
	}()

	// Create a command configuration
	commandConfig := cli.NewCommandConfig(
		os.Stdout,
		os.Stderr,
		os.Stdin,
	)

	// Retrieve args and Shift binary name off argument list.
	args := os.Args[1:]

	if err := cli.Run(ctx, commandConfig, args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s.\n", err)
		os.Exit(1)
	}
}
