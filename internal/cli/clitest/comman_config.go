// Package clitest provides utilities for testing the CLI.
package clitest

import (
	"io"
	"time"

	"github.com/aviau/gopass/pkg/store"
)

// testCommandConfig is a CommandConfig for testing.
type testCommandConfig struct {
	passwordStore *store.PasswordStore
	runOptions    *runOptions
	writerOutput  io.Writer
	writerError   io.Writer
	readerInput   io.Reader
}

func (cfg *testCommandConfig) PasswordStore() *store.PasswordStore {
	return cfg.passwordStore
}

func (cfg *testCommandConfig) PasswordStoreDir() string {
	return cfg.passwordStore.Path
}

func (cfg *testCommandConfig) Edit(content string) (string, error) {
	if cfg.runOptions.editFunc != nil {
		return cfg.runOptions.editFunc(content)
	}
	return content, nil
}

func (cfg *testCommandConfig) WriterOutput() io.Writer {
	return cfg.writerOutput
}

func (cfg *testCommandConfig) WriterError() io.Writer {
	return cfg.writerError
}

func (cfg *testCommandConfig) ReaderInput() io.Reader {
	return cfg.readerInput
}

func (cfg *testCommandConfig) Now() time.Time {
	if cfg.runOptions.nowFunc != nil {
		return cfg.runOptions.nowFunc()
	}
	return time.Now()
}
