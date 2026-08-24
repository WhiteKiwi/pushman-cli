package credential

import (
	"errors"
	"testing"

	keyring "github.com/zalando/go-keyring"
)

func TestKeyringLifecycle(t *testing.T) {
	keyring.MockInit()
	store := NewKeyring("com.pushman.cli.test")
	if _, err := store.Get(); !errors.Is(err, ErrNotFound) {
		t.Fatalf("initial get error = %v", err)
	}
	if err := store.Set("pm_cli_secret"); err != nil {
		t.Fatal(err)
	}
	if value, err := store.Get(); err != nil || value != "pm_cli_secret" {
		t.Fatalf("get = %q, %v", value, err)
	}
	if err := store.Delete(); err != nil {
		t.Fatal(err)
	}
	if _, err := store.Get(); !errors.Is(err, ErrNotFound) {
		t.Fatalf("deleted get error = %v", err)
	}
}

func TestKeyringRejectsEmptyCredential(t *testing.T) {
	keyring.MockInit()
	if err := NewKeyring("com.pushman.cli.test").Set(""); err == nil {
		t.Fatal("expected empty credential to fail")
	}
}
