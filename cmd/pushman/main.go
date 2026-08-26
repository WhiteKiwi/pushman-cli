package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime/debug"
	"strings"

	"github.com/WhiteKiwi/pushman-cli/internal/browser"
	"github.com/WhiteKiwi/pushman-cli/internal/cli"
	pushclient "github.com/WhiteKiwi/pushman-cli/internal/client"
	"github.com/WhiteKiwi/pushman-cli/internal/credential"
	"github.com/WhiteKiwi/pushman-cli/internal/selfupdate"
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

	resolvedVersion, resolvedCommit, resolvedDate := resolveVersionInfo(version, commit, date, readBuildInfo())
	updater := selfupdate.New()
	app := cli.New(cli.Dependencies{
		In:          os.Stdin,
		Out:         os.Stdout,
		ErrOut:      os.Stderr,
		IsTerminal:  cli.IsTerminalFile(os.Stdin),
		OpenBrowser: browser.Open,
		SelfUpdate:  updater.Update,
		Service:     service,
		Version: cli.VersionInfo{
			Version: resolvedVersion,
			Commit:  resolvedCommit,
			Date:    resolvedDate,
		},
	})

	if err := app.ExecuteContext(ctx); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(cli.ExitCode(err))
	}
}

func readBuildInfo() *debug.BuildInfo {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return nil
	}
	return info
}

func resolveVersionInfo(version, commit, date string, info *debug.BuildInfo) (string, string, string) {
	if info == nil {
		return version, commit, date
	}
	if version == "dev" && info.Main.Version != "" && info.Main.Version != "(devel)" {
		version = strings.TrimPrefix(info.Main.Version, "v")
	}
	for _, setting := range info.Settings {
		switch setting.Key {
		case "vcs.revision":
			if commit == "none" && setting.Value != "" {
				commit = setting.Value
			}
		case "vcs.time":
			if date == "unknown" && setting.Value != "" {
				date = setting.Value
			}
		}
	}
	return version, commit, date
}
