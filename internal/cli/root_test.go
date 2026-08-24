package cli

import (
	"bytes"
	"strings"
	"testing"
)

func TestVersionUsesInjectedOutput(t *testing.T) {
	t.Parallel()
	out := new(bytes.Buffer)
	root := New(Dependencies{Out: out, Version: VersionInfo{Version: "1.2.3", Commit: "abc", Date: "2026-08-25"}})
	root.SetArgs([]string{"version"})
	if err := root.Execute(); err != nil {
		t.Fatal(err)
	}
	if got := out.String(); !strings.Contains(got, "1.2.3") || !strings.Contains(got, "abc") {
		t.Fatalf("version output = %q", got)
	}
}

func TestInvalidArgumentIsUsageError(t *testing.T) {
	t.Parallel()
	root := New(Dependencies{})
	root.SetArgs([]string{"status", "extra"})
	err := root.Execute()
	if err == nil {
		t.Fatal("expected error")
	}
	if got := ExitCode(err); got != 2 {
		t.Fatalf("exit code = %d, want 2", got)
	}
}
