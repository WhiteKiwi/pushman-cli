package selfupdate

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

const formula = "whitekiwi/tap/pushman"

type commandRunner func(context.Context, string, ...string) (string, error)

type Updater struct {
	executable   func() (string, error)
	lookPath     func(string) (string, error)
	evalSymlinks func(string) (string, error)
	run          commandRunner
	goos         string
}

func New() *Updater {
	return &Updater{
		executable:   os.Executable,
		lookPath:     exec.LookPath,
		evalSymlinks: filepath.EvalSymlinks,
		run: func(ctx context.Context, name string, args ...string) (string, error) {
			output, err := exec.CommandContext(ctx, name, args...).CombinedOutput()
			return strings.TrimSpace(string(output)), err
		},
		goos: runtime.GOOS,
	}
}

func (u *Updater) Update(ctx context.Context) (string, error) {
	if u.goos != "darwin" && u.goos != "linux" {
		return "", fmt.Errorf("self-update supports Homebrew installations on macOS and Linux; use the installation guide for this platform")
	}

	brew, err := u.lookPath("brew")
	if err != nil {
		return "", fmt.Errorf("self-update requires a Homebrew-managed Pushman installation; update using the original installation method")
	}
	prefixOutput, err := u.run(ctx, brew, "--prefix", formula)
	if err != nil {
		return "", commandError("locate the Homebrew Pushman formula", prefixOutput, err)
	}
	prefix := strings.TrimSpace(prefixOutput)
	if !filepath.IsAbs(prefix) {
		return "", fmt.Errorf("Homebrew returned an invalid Pushman prefix")
	}

	executable, err := u.executable()
	if err != nil {
		return "", fmt.Errorf("locate the running Pushman executable: %w", err)
	}
	running, err := u.evalSymlinks(executable)
	if err != nil {
		return "", fmt.Errorf("resolve the running Pushman executable: %w", err)
	}
	managed, err := u.evalSymlinks(filepath.Join(prefix, "bin", "pushman"))
	if err != nil {
		return "", fmt.Errorf("resolve the Homebrew Pushman executable: %w", err)
	}
	if running != managed {
		return "", fmt.Errorf("self-update refused to replace a non-Homebrew Pushman installation; update using the original installation method")
	}

	output, err := u.run(ctx, brew, "upgrade", formula)
	if err != nil {
		return "", commandError("upgrade Pushman with Homebrew", output, err)
	}
	if output == "" {
		return "Pushman is already up to date.", nil
	}
	return output, nil
}

func commandError(action, output string, err error) error {
	if output == "" {
		return fmt.Errorf("%s: %w", action, err)
	}
	return fmt.Errorf("%s: %s: %w", action, output, err)
}
