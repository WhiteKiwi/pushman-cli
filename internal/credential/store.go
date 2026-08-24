package credential

import (
	"errors"
	"fmt"

	keyring "github.com/zalando/go-keyring"
)

var ErrNotFound = errors.New("paired credential not found")

type Store interface {
	Get() (string, error)
	Set(string) error
	Delete() error
}

type Keyring struct {
	service string
	user    string
}

func NewKeyring(service string) *Keyring {
	return &Keyring{service: service, user: "paired-credential"}
}

func (k *Keyring) Get() (string, error) {
	token, err := keyring.Get(k.service, k.user)
	if errors.Is(err, keyring.ErrNotFound) {
		return "", ErrNotFound
	}
	if err != nil {
		return "", fmt.Errorf("read paired credential from operating-system keyring: %w", err)
	}
	if token == "" {
		return "", ErrNotFound
	}
	return token, nil
}

func (k *Keyring) Set(token string) error {
	if token == "" {
		return fmt.Errorf("refuse to store an empty credential")
	}
	if err := keyring.Set(k.service, k.user, token); err != nil {
		return fmt.Errorf("store paired credential in operating-system keyring: %w", err)
	}
	return nil
}

func (k *Keyring) Delete() error {
	err := keyring.Delete(k.service, k.user)
	if errors.Is(err, keyring.ErrNotFound) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("delete paired credential from operating-system keyring: %w", err)
	}
	return nil
}
