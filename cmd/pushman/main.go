package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime/debug"
	"strings"

	"github.com/pushmanhq/pushman-cli/internal/browser"
	"github.com/pushmanhq/pushman-cli/internal/cli"
	pushclient "github.com/pushmanhq/pushman-cli/internal/client"
	"github.com/pushmanhq/pushman-cli/internal/credential"
	"github.com/pushmanhq/pushman-cli/internal/selfupdate"
)

var (
	version                    = "dev"
	commit                     = "none"
	date                       = "unknown"
	defaultBaseURL             = pushclient.DefaultBaseURL
	credentialNamespace        = ""
	automationTokenEnvironment = "PUSHMAN_TOKEN"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	baseURL := os.Getenv("PUSHMAN_API_URL")
	if baseURL == "" {
		baseURL = defaultBaseURL
	}
	validatedBaseURL, err := pushclient.ValidateBaseURL(baseURL)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(cli.ExitCode(err))
	}
	keyring := credential.NewKeyring(credentialServiceName(validatedBaseURL, credentialNamespace))
	service, err := pushclient.New(validatedBaseURL, keyring, os.Getenv(automationTokenEnvironment), nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(cli.ExitCode(err))
	}

	resolvedVersion, resolvedCommit, resolvedDate := resolveVersionInfo(version, commit, date, readBuildInfo())
	var update func(context.Context) (string, error)
	if credentialNamespace == "" {
		updater := selfupdate.New()
		update = updater.Update
	}
	app := cli.New(cli.Dependencies{
		In:          os.Stdin,
		Out:         os.Stdout,
		ErrOut:      os.Stderr,
		IsTerminal:  cli.IsTerminalFile(os.Stdin),
		OpenBrowser: browser.Open,
		SelfUpdate:  update,
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

func credentialServiceName(baseURL, namespace string) string {
	service := pushclient.CredentialServiceName(baseURL)
	if namespace == "" {
		return service
	}
	return service + "." + namespace
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
