package browser

import (
	"fmt"
	"os/exec"
	"runtime"
)

func Open(rawURL string) error {
	var command *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		command = exec.Command("open", rawURL)
	case "linux":
		command = exec.Command("xdg-open", rawURL)
	case "windows":
		command = exec.Command("rundll32", "url.dll,FileProtocolHandler", rawURL)
	default:
		return fmt.Errorf("browser launch is not supported on %s", runtime.GOOS)
	}
	if err := command.Start(); err != nil {
		return err
	}
	return command.Process.Release()
}
