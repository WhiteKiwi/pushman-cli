package cli

import (
	"fmt"
	"runtime"
	"strings"
	"unicode/utf8"

	"github.com/spf13/cobra"
	"golang.org/x/text/unicode/norm"
)

func newPairCommand(deps Dependencies) *cobra.Command {
	cmd := &cobra.Command{Use: "pair", Short: "Pair this CLI with Pushman", Args: func(_ *cobra.Command, args []string) error { return noArgs(args) }, RunE: func(cmd *cobra.Command, _ []string) error {
		hostname, err := deps.Hostname()
		if err != nil {
			return fmt.Errorf("read hostname for pairing: %w", err)
		}
		suggestedName, err := normalizeSuggestedName(hostname)
		if err != nil {
			return err
		}
		result, err := deps.Service.Pair(cmd.Context(), PairRequest{
			Platform: runtime.GOOS, SuggestedName: suggestedName,
			OnChallenge: func(challenge PairChallenge) error {
				_, err := fmt.Fprintf(cmd.OutOrStdout(), "Pairing code: %s\nVerify at: %s\n", challenge.UserCode, challenge.VerificationURL)
				return err
			},
		})
		if err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Paired as %s\n", result.Nickname)
		return nil
	}}
	return cmd
}

func normalizeSuggestedName(value string) (string, error) {
	name := norm.NFC.String(strings.TrimSpace(value))
	if !utf8.ValidString(name) || name == "" {
		return "", usagef("hostname must contain valid non-whitespace UTF-8")
	}
	runes := []rune(name)
	if len(runes) > 64 {
		name = string(runes[:64])
	}
	return name, nil
}

func newStatusCommand(deps Dependencies) *cobra.Command {
	return &cobra.Command{Use: "status", Short: "Show pairing status", Args: func(_ *cobra.Command, args []string) error { return noArgs(args) }, RunE: func(cmd *cobra.Command, _ []string) error {
		result, err := deps.Service.Status(cmd.Context())
		if err != nil {
			return err
		}
		if !result.Paired {
			fmt.Fprintln(cmd.OutOrStdout(), "Not paired")
			return nil
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Paired as %s\n", result.Nickname)
		return nil
	}}
}

func newRenameCommand(deps Dependencies) *cobra.Command {
	return &cobra.Command{Use: "rename <nickname>", Short: "Rename this paired CLI", Args: func(_ *cobra.Command, args []string) error { return exactlyOne(args, "nickname") }, RunE: func(cmd *cobra.Command, args []string) error {
		name := norm.NFC.String(strings.TrimSpace(args[0]))
		if !utf8.ValidString(name) || countScalars(name) < 1 || countScalars(name) > 64 {
			return usagef("nickname must contain 1 through 64 Unicode scalar values")
		}
		if err := deps.Service.Rename(cmd.Context(), name); err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Renamed to %s\n", name)
		return nil
	}}
}

func newLogoutCommand(deps Dependencies) *cobra.Command {
	return &cobra.Command{Use: "logout", Short: "Revoke this CLI credential", Args: func(_ *cobra.Command, args []string) error { return noArgs(args) }, RunE: func(cmd *cobra.Command, _ []string) error {
		if err := deps.Service.Logout(cmd.Context()); err != nil {
			return err
		}
		fmt.Fprintln(cmd.OutOrStdout(), "Logged out")
		return nil
	}}
}

func newDevicesCommand(deps Dependencies) *cobra.Command {
	return &cobra.Command{Use: "devices", Short: "List receiving devices", Args: func(_ *cobra.Command, args []string) error { return noArgs(args) }, RunE: func(cmd *cobra.Command, _ []string) error {
		devices, err := deps.Service.Devices(cmd.Context())
		if err != nil {
			return err
		}
		for _, device := range devices {
			fmt.Fprintf(cmd.OutOrStdout(), "%s\t%s\n", device.Nickname, device.Status)
		}
		return nil
	}}
}

func newHistoryCommand(deps Dependencies) *cobra.Command {
	cmd := &cobra.Command{Use: "history", Short: "List message history", Args: func(_ *cobra.Command, args []string) error { return noArgs(args) }, RunE: func(cmd *cobra.Command, _ []string) error {
		items, err := deps.Service.History(cmd.Context())
		if err != nil {
			return err
		}
		for _, item := range items {
			fmt.Fprintf(cmd.OutOrStdout(), "%s\t%s\t%s\n", item.ID, item.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"), item.Title)
		}
		return nil
	}}
	cmd.AddCommand(&cobra.Command{Use: "show <message-id>", Short: "Show a message and its revisions", Args: func(_ *cobra.Command, args []string) error { return exactlyOne(args, "message ID") }, RunE: func(cmd *cobra.Command, args []string) error {
		detail, err := deps.Service.HistoryShow(cmd.Context(), args[0])
		if err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "%s\n%s\n\n%s\n", detail.ID, detail.Title, detail.Body)
		return nil
	}})
	return cmd
}

func newUsageCommand(deps Dependencies) *cobra.Command {
	return &cobra.Command{Use: "usage", Short: "Show monthly usage", Args: func(_ *cobra.Command, args []string) error { return noArgs(args) }, RunE: func(cmd *cobra.Command, _ []string) error {
		usage, err := deps.Service.Usage(cmd.Context())
		if err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "%d of %d messages used; resets %s\n", usage.Used, usage.Limit, usage.ResetsAt.Format("2006-01-02T15:04:05Z07:00"))
		return nil
	}}
}

func newDoctorCommand(deps Dependencies) *cobra.Command {
	return &cobra.Command{Use: "doctor", Short: "Diagnose local and service connectivity", Args: func(_ *cobra.Command, args []string) error { return noArgs(args) }, RunE: func(cmd *cobra.Command, _ []string) error {
		checks, err := deps.Service.Doctor(cmd.Context())
		if err != nil {
			return err
		}
		for _, check := range checks {
			state := "ok"
			if !check.OK {
				state = "failed"
			}
			fmt.Fprintf(cmd.OutOrStdout(), "%s\t%s\t%s\n", state, check.Name, check.Message)
		}
		return nil
	}}
}

func newVersionCommand(deps Dependencies) *cobra.Command {
	return &cobra.Command{Use: "version", Short: "Print version information", Args: func(_ *cobra.Command, args []string) error { return noArgs(args) }, RunE: func(cmd *cobra.Command, _ []string) error {
		fmt.Fprintf(cmd.OutOrStdout(), "pushman %s (commit %s, built %s)\n", deps.Version.Version, deps.Version.Commit, deps.Version.Date)
		return nil
	}}
}
