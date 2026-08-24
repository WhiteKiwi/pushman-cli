//go:build !windows

package cli

import "os"

func IsTerminalFile(file *os.File) func() bool {
	return func() bool {
		info, err := file.Stat()
		return err == nil && info.Mode()&os.ModeCharDevice != 0
	}
}
