package main

import (
	"runtime/debug"
	"testing"
)

func TestResolveVersionInfoUsesModuleAndVCSMetadata(t *testing.T) {
	info := &debug.BuildInfo{
		Main: debug.Module{Version: "v0.1.0-beta.2"},
		Settings: []debug.BuildSetting{
			{Key: "vcs.revision", Value: "abc123"},
			{Key: "vcs.time", Value: "2026-08-25T10:00:00Z"},
		},
	}

	version, commit, date := resolveVersionInfo("dev", "none", "unknown", info)
	if version != "0.1.0-beta.2" || commit != "abc123" || date != "2026-08-25T10:00:00Z" {
		t.Fatalf("resolveVersionInfo() = %q, %q, %q", version, commit, date)
	}
}

func TestResolveVersionInfoKeepsLinkerValues(t *testing.T) {
	info := &debug.BuildInfo{
		Main: debug.Module{Version: "v0.1.0-beta.2"},
		Settings: []debug.BuildSetting{
			{Key: "vcs.revision", Value: "build-commit"},
			{Key: "vcs.time", Value: "2026-08-25T10:00:00Z"},
		},
	}

	version, commit, date := resolveVersionInfo("release", "linker-commit", "linker-date", info)
	if version != "release" || commit != "linker-commit" || date != "linker-date" {
		t.Fatalf("resolveVersionInfo() = %q, %q, %q", version, commit, date)
	}
}

func TestResolveVersionInfoIgnoresDevelopmentModuleVersion(t *testing.T) {
	version, commit, date := resolveVersionInfo("dev", "none", "unknown", &debug.BuildInfo{Main: debug.Module{Version: "(devel)"}})
	if version != "dev" || commit != "none" || date != "unknown" {
		t.Fatalf("resolveVersionInfo() = %q, %q, %q", version, commit, date)
	}
}

func TestCredentialServiceNameSeparatesDevelopmentCredentials(t *testing.T) {
	baseURL := "https://api.pushman.example/v1"
	production := credentialServiceName(baseURL, "")
	development := credentialServiceName(baseURL, "dev")

	if production == development {
		t.Fatal("development credential service must not match production")
	}
	if development != production+".dev" {
		t.Fatalf("development credential service = %q", development)
	}
}
