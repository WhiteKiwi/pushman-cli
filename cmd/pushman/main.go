package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/WhiteKiwi/pushman-cli/internal/cli"
	pushclient "github.com/WhiteKiwi/pushman-cli/internal/client"
	"github.com/WhiteKiwi/pushman-cli/internal/credential"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	baseURL := os.Getenv("PUSHMAN_API_URL")
	if baseURL == "" {
		baseURL = pushclient.DefaultBaseURL
	}
	validatedBaseURL, err := pushclient.ValidateBaseURL(baseURL)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(cli.ExitCode(err))
	}
	keyring := credential.NewKeyring(pushclient.CredentialServiceName(validatedBaseURL))
	service, err := pushclient.New(validatedBaseURL, keyring, os.Getenv("PUSHMAN_TOKEN"), nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(cli.ExitCode(err))
	}

	app := cli.New(cli.Dependencies{
		In:         os.Stdin,
		Out:        os.Stdout,
		ErrOut:     os.Stderr,
		IsTerminal: cli.IsTerminalFile(os.Stdin),
		Service:    service,
		Version: cli.VersionInfo{
			Version: version,
			Commit:  commit,
			Date:    date,
		},
	})

	if err := app.ExecuteContext(ctx); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(cli.ExitCode(err))
	}
}
