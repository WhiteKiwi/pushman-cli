package cli

import (
	"io"
	"os"

	"github.com/spf13/cobra"
)

type Dependencies struct {
	In         io.Reader
	Out        io.Writer
	ErrOut     io.Writer
	IsTerminal func() bool
	Hostname   func() (string, error)
	Service    Service
	Version    VersionInfo
}

func New(deps Dependencies) *cobra.Command {
	deps = withDefaults(deps)
	root := &cobra.Command{
		Use:           "pushman",
		Short:         "Send push notifications to your phone",
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Help()
		},
	}
	root.SetIn(deps.In)
	root.SetOut(deps.Out)
	root.SetErr(deps.ErrOut)
	root.CompletionOptions.DisableDefaultCmd = true
	root.SetFlagErrorFunc(func(_ *cobra.Command, err error) error {
		return &UsageError{Err: err}
	})
	root.AddCommand(
		newPairCommand(deps),
		newStatusCommand(deps),
		newRenameCommand(deps),
		newLogoutCommand(deps),
		newPushCommand(deps),
		newDevicesCommand(deps),
		newHistoryCommand(deps),
		newUsageCommand(deps),
		newDoctorCommand(deps),
		newVersionCommand(deps),
	)
	root.InitDefaultHelpCmd()
	return root
}

func withDefaults(deps Dependencies) Dependencies {
	if deps.In == nil {
		deps.In = os.Stdin
	}
	if deps.Out == nil {
		deps.Out = io.Discard
	}
	if deps.ErrOut == nil {
		deps.ErrOut = io.Discard
	}
	if deps.IsTerminal == nil {
		deps.IsTerminal = func() bool { return true }
	}
	if deps.Hostname == nil {
		deps.Hostname = os.Hostname
	}
	if deps.Service == nil {
		deps.Service = UnconfiguredService{}
	}
	return deps
}

func noArgs(args []string) error {
	if len(args) != 0 {
		return usagef("expected no arguments")
	}
	return nil
}

func exactlyOne(args []string, name string) error {
	if len(args) != 1 {
		return usagef("expected exactly one %s", name)
	}
	return nil
}
