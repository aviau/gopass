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

package gopass_test

import (
	"github.com/stretchr/testify/assert"
	"testing"

	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/aviau/gopass"
)

type passwordStoreTest struct {
	PasswordStore *gopass.PasswordStore
	StorePath     string
}

func newPasswordStoreTest() (*passwordStoreTest, error) {
	storePath, err := ioutil.TempDir("", "gopass")
	if err != nil {
		return nil, err
	}

	passwordStore := gopass.NewPasswordStore(storePath)
	passwordStore.UsesGit = false

	err = passwordStore.Init("test")
	if err != nil {
		return nil, err
	}

	passwordStoreTest := passwordStoreTest{
		PasswordStore: passwordStore,
		StorePath:     storePath,
	}

	return &passwordStoreTest, nil
}

func (test *passwordStoreTest) Close() error {
	err := os.RemoveAll(test.StorePath)
	return err
}

func TestRemoveDirectory(t *testing.T) {
	st, err := newPasswordStoreTest()
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()

	testDirectoryPath := filepath.Join(st.StorePath, "dir")
	err = os.Mkdir(testDirectoryPath, 0700)
	if err != nil {
		t.Fatal(err)
	}

	_, err = os.Stat(testDirectoryPath)
	assert.False(
		t,
		os.IsNotExist(err),
		"dir should have been created",
	)

	err = st.PasswordStore.RemoveDirectory("dir")
	if err != nil {
		t.Fatal(err)
	}

	_, err = os.Stat(testDirectoryPath)
	assert.True(
		t,
		os.IsNotExist(err),
		"dir should have been removed",
	)
}

func TestRemoveDirectoryTrailingSlash(t *testing.T) {
	st, err := newPasswordStoreTest()
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()

	testDirectoryPath := filepath.Join(st.StorePath, "dir")
	err = os.Mkdir(testDirectoryPath, 0700)
	if err != nil {
		t.Fatal(err)
	}

	_, err = os.Stat(testDirectoryPath)
	assert.False(
		t,
		os.IsNotExist(err),
		"dir should have been created",
	)

	err = st.PasswordStore.RemoveDirectory("//dir///")
	if err != nil {
		t.Fatal(err)
	}

	_, err = os.Stat(testDirectoryPath)
	assert.True(
		t,
		os.IsNotExist(err),
		"dir should have been removed",
	)
}

func TestRemoveDirectoryRecursive(t *testing.T) {
	st, err := newPasswordStoreTest()
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()

	testDirectoryPath := filepath.Join(st.StorePath, "dir")
	err = os.Mkdir(testDirectoryPath, 0700)
	if err != nil {
		t.Fatal(err)
	}

	testPasswordPath := filepath.Join(testDirectoryPath, "test.com.gpg")
	_, err = os.Create(testPasswordPath)
	if err != nil {
		t.Fatal(err)
	}

	_, err = os.Stat(testPasswordPath)
	assert.False(
		t,
		os.IsNotExist(err),
		"test.com.gpg should have been created",
	)

	_, err = os.Stat(testDirectoryPath)
	assert.False(
		t,
		os.IsNotExist(err),
		"dir should have been created",
	)

	err = st.PasswordStore.RemoveDirectory("dir")
	if err != nil {
		t.Fatal(err)
	}

	_, err = os.Stat(testDirectoryPath)
	assert.True(
		t,
		os.IsNotExist(err),
		"dir should have been removed",
	)

	_, err = os.Stat(testPasswordPath)
	assert.True(
		t,
		os.IsNotExist(err),
		"test.com.gpg should have been removed",
	)
}

func TestRemovePassword(t *testing.T) {
	st, err := newPasswordStoreTest()
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()

	testPasswordPath := filepath.Join(st.StorePath, "test.com.gpg")
	_, err = os.Create(testPasswordPath)
	if err != nil {
		t.Fatal(err)
	}

	_, err = os.Stat(testPasswordPath)
	assert.False(
		t,
		os.IsNotExist(err),
		"test.com.gpg should have been created",
	)

	st.PasswordStore.RemovePassword("test.com")
	_, err = os.Stat(testPasswordPath)
	assert.True(
		t,
		os.IsNotExist(err),
		"test.com should have been removed",
	)
}

func TestRemovePasswordTrailingSlash(t *testing.T) {
	st, err := newPasswordStoreTest()
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()

	testPasswordPath := filepath.Join(st.StorePath, "test.com.gpg")
	_, err = os.Create(testPasswordPath)
	if err != nil {
		t.Fatal(err)
	}

	_, err = os.Stat(testPasswordPath)
	assert.False(
		t,
		os.IsNotExist(err),
		"test.com.gpg should have been created",
	)

	st.PasswordStore.RemovePassword("test.com/")
	_, err = os.Stat(testPasswordPath)
	assert.False(
		t,
		os.IsNotExist(err),
		"RemovePassword with a trailing slash should not remove a password",
	)
}

func TestRemovePasswordDirectory(t *testing.T) {
	st, err := newPasswordStoreTest()
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()

	testDirectoryPath := filepath.Join(st.StorePath, "test.com")
	err = os.Mkdir(testDirectoryPath, 0700)
	if err != nil {
		t.Fatal(err)
	}

	_, err = os.Stat(testDirectoryPath)
	assert.False(
		t,
		os.IsNotExist(err),
		"test.com.gpg should have been created",
	)

	st.PasswordStore.RemovePassword("test.com")
	_, err = os.Stat(testDirectoryPath)
	assert.False(
		t,
		os.IsNotExist(err),
		"RemovePassword should not remove directories",
	)
}

func TestCopyPassword(t *testing.T) {
	st, err := newPasswordStoreTest()
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()

	testPasswordPath := filepath.Join(st.StorePath, "test.com.gpg")
	_, err = os.Create(testPasswordPath)
	if err != nil {
		t.Fatal(err)
	}

	_, err = os.Stat(testPasswordPath)
	assert.False(
		t,
		os.IsNotExist(err),
		"test.com.gpg should have been created",
	)

	destPasswordPath := filepath.Join(st.StorePath, "test2.com.gpg")
	_, err = os.Stat(destPasswordPath)
	assert.True(
		t,
		os.IsNotExist(err),
		"test2.com.gpg should not have been created yet",
	)

	st.PasswordStore.CopyPassword("test.com", "test2.com")
	_, err = os.Stat(destPasswordPath)
	assert.False(
		t,
		os.IsNotExist(err),
		"test.com.gpg shoudl have been copied to test2.com.gpg",
	)
}

func TestCopyPasswordInDirectory(t *testing.T) {
	st, err := newPasswordStoreTest()
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()

	testPasswordPath := filepath.Join(st.StorePath, "test.com.gpg")
	_, err = os.Create(testPasswordPath)
	if err != nil {
		t.Fatal(err)
	}

	_, err = os.Stat(testPasswordPath)
	assert.False(
		t,
		os.IsNotExist(err),
		"test.com.gpg should have been created",
	)

	testDirectoryPath := filepath.Join(st.StorePath, "dir")
	err = os.Mkdir(testDirectoryPath, 0700)
	if err != nil {
		t.Fatal(err)
	}

	destPasswordPath := filepath.Join(st.StorePath, "dir", "test.com.gpg")
	_, err = os.Stat(destPasswordPath)
	assert.True(
		t,
		os.IsNotExist(err),
		"test2.com.gpg should have been created yet",
	)

	st.PasswordStore.CopyPassword("test.com", "dir/")
	_, err = os.Stat(destPasswordPath)
	assert.False(
		t,
		os.IsNotExist(err),
		"test.com.gpg shoudl have been copied to dir/test2.com.gpg",
	)
}

func TestContainsPassword(t *testing.T) {
	st, err := newPasswordStoreTest()
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()

	containsPassword, _ := st.PasswordStore.ContainsPassword("test.com")
	assert.False(t, containsPassword, "The password store should not contain test.com")

	testPasswordPath := filepath.Join(st.StorePath, "test.com.gpg")
	_, err = os.Create(testPasswordPath)
	if err != nil {
		t.Fatal(err)
	}

	_, err = os.Stat(testPasswordPath)
	assert.False(
		t,
		os.IsNotExist(err),
		"test.com.gpg should have been created",
	)

	containsPassword, _ = st.PasswordStore.ContainsPassword("test.com")
	assert.True(t, containsPassword, "The password store should contain test.com")
}

func TestContainsPasswordDirectory(t *testing.T) {
	st, err := newPasswordStoreTest()
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()

	containsPassword, _ := st.PasswordStore.ContainsPassword("test.com")
	assert.False(t, containsPassword, "The password store should not contain test.com")

	testDirectoryPath := filepath.Join(st.StorePath, "test.com")
	err = os.Mkdir(testDirectoryPath, 0700)
	if err != nil {
		t.Fatal(err)
	}

	_, err = os.Stat(testDirectoryPath)
	assert.False(
		t,
		os.IsNotExist(err),
		"test.com directory should have been created",
	)

	containsPassword, _ = st.PasswordStore.ContainsPassword("test.com")
	assert.False(t, containsPassword, "The password store should not a password named test.com")
}

func TestContainsDirectory(t *testing.T) {
	st, err := newPasswordStoreTest()
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()

	containsDirectory, _ := st.PasswordStore.ContainsDirectory("dir")
	assert.False(t, containsDirectory, "The password store should contain dir")

	testDirectoryPath := filepath.Join(st.StorePath, "dir")
	err = os.Mkdir(testDirectoryPath, 0700)
	if err != nil {
		t.Fatal(err)
	}

	_, err = os.Stat(testDirectoryPath)
	assert.False(
		t,
		os.IsNotExist(err),
		"the directory should have been created",
	)

	containsDirectory, _ = st.PasswordStore.ContainsDirectory("dir")
	assert.True(t, containsDirectory, "The password store should contain dir")
}

func TestContainsDirectoryTrailingSlash(t *testing.T) {
	st, err := newPasswordStoreTest()
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()

	containsDirectory, _ := st.PasswordStore.ContainsDirectory("dir")
	assert.False(t, containsDirectory, "The password store should contain dir")

	testDirectoryPath := filepath.Join(st.StorePath, "dir")
	err = os.Mkdir(testDirectoryPath, 0700)
	if err != nil {
		t.Fatal(err)
	}

	_, err = os.Stat(testDirectoryPath)
	assert.False(
		t,
		os.IsNotExist(err),
		"the directory should have been created",
	)

	containsDirectory, _ = st.PasswordStore.ContainsDirectory("/dir///")
	assert.True(t, containsDirectory, "The password store should contain dir")
}

func TestGetPasswordsList(t *testing.T) {
	st, err := newPasswordStoreTest()
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()

	_, err = os.Create(filepath.Join(st.StorePath, "test.com.gpg"))
	if err != nil {
		t.Fatal(err)
	}

	_, err = os.Create(filepath.Join(st.StorePath, "test2.com.gpg"))
	if err != nil {
		t.Fatal(err)
	}

	_, err = os.Create(filepath.Join(st.StorePath, "test3"))
	if err != nil {
		t.Fatal(err)
	}

	passwords := st.PasswordStore.GetPasswordsList()
	assert.Equal(
		t,
		passwords,
		[]string{"test.com", "test2.com"},
		"Password list should contain test.com and test2.com",
	)
}
