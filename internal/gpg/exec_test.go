package gpg_test

import (
	"testing"

	"github.com/aviau/gopass/internal/gpg"
	"github.com/aviau/gopass/pkg/store"
)

func TestImplementsGPG(t *testing.T) {
	_ = func() store.GPGBackend {
		return gpg.New("", nil, false)
	}
}
