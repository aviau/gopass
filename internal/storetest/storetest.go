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

// Package storetest provides utilities for testing the store.
package storetest

import (
	"os"
	"testing"

	_ "embed"

	"github.com/aviau/gopass/internal/gpg"
	"github.com/aviau/gopass/pkg/store"
)

//go:embed testdata/CED3B67C8F1F6CA9.private.key
var testSecretKey []byte
var testSecretKeyID string = "CED3B67C8F1F6CA9"

// PasswordStoreTest allows for testing password stores.
type PasswordStoreTest struct {
	PasswordStore *store.PasswordStore
	storePath     string
}

// NewPasswordStoreTest creates a password store for testing
func NewPasswordStoreTest(t *testing.T) *PasswordStoreTest {
	storePath := t.TempDir()

	// Prepare GPG
	gnupgHome := t.TempDir()
	if err := os.Chmod(gnupgHome, 0700); err != nil {
		t.Fatal(err)
	}
	gpgBackend := gpg.New(
		"",
		[]string{
			"GNUPGHOME=" + gnupgHome,
		},
		true,
	)
	if err := gpgBackend.Import(testSecretKey); err != nil {
		t.Fatal(err)
	}

	// Create the store
	passwordStore := store.NewPasswordStore(storePath)
	passwordStore.UsesGit = false
	passwordStore.GPGBackend = gpgBackend

	// Init the store
	if err := passwordStore.Init([]string{testSecretKeyID}); err != nil {
		t.Fatal(err)
	}

	passwordStoreTest := PasswordStoreTest{
		PasswordStore: passwordStore,
		storePath:     storePath,
	}

	return &passwordStoreTest
}

// Close cleans the test.
func (test *PasswordStoreTest) Close() {
	// Nothing for now.
}
