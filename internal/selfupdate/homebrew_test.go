package selfupdate

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"
)

func TestUpdateUsesHomebrewForTheRunningFormula(t *testing.T) {
	var calls [][]string
	updater := &Updater{
		goos:       "darwin",
		executable: func() (string, error) { return "/opt/homebrew/bin/pushman", nil },
		lookPath:   func(string) (string, error) { return "/opt/homebrew/bin/brew", nil },
		evalSymlinks: func(string) (string, error) {
			return "/opt/homebrew/Cellar/pushman/0.1.0/bin/pushman", nil
		},
		run: func(_ context.Context, name string, args ...string) (string, error) {
			calls = append(calls, append([]string{name}, args...))
			if reflect.DeepEqual(args, []string{"--prefix", formula}) {
				return "/opt/homebrew/opt/pushman\n", nil
			}
			return "upgraded pushman", nil
		},
	}

	output, err := updater.Update(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if output != "upgraded pushman" {
		t.Fatalf("output = %q", output)
	}
	want := [][]string{
		{"/opt/homebrew/bin/brew", "--prefix", formula},
		{"/opt/homebrew/bin/brew", "upgrade", formula},
	}
	if !reflect.DeepEqual(calls, want) {
		t.Fatalf("calls = %#v, want %#v", calls, want)
	}
}

func TestUpdateRefusesAnUnmanagedExecutable(t *testing.T) {
	updater := &Updater{
		goos:       "linux",
		executable: func() (string, error) { return "/usr/local/bin/pushman", nil },
		lookPath:   func(string) (string, error) { return "/home/linuxbrew/.linuxbrew/bin/brew", nil },
		evalSymlinks: func(path string) (string, error) {
			if path == "/usr/local/bin/pushman" {
				return path, nil
			}
			return "/home/linuxbrew/.linuxbrew/Cellar/pushman/0.1.0/bin/pushman", nil
		},
		run: func(_ context.Context, _ string, args ...string) (string, error) {
			if reflect.DeepEqual(args, []string{"--prefix", formula}) {
				return "/home/linuxbrew/.linuxbrew/opt/pushman", nil
			}
			t.Fatal("upgrade ran for an unmanaged executable")
			return "", nil
		},
	}

	_, err := updater.Update(context.Background())
	if err == nil || !strings.Contains(err.Error(), "non-Homebrew") {
		t.Fatalf("error = %v", err)
	}
}

func TestUpdateExplainsWhenHomebrewIsUnavailable(t *testing.T) {
	updater := &Updater{
		goos:       "darwin",
		executable: func() (string, error) { return "", nil },
		lookPath:   func(string) (string, error) { return "", errors.New("not found") },
	}

	_, err := updater.Update(context.Background())
	if err == nil || !strings.Contains(err.Error(), "Homebrew-managed") {
		t.Fatalf("error = %v", err)
	}
}
