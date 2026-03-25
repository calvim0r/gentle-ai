package kiro

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/gentleman-programming/gentle-ai/internal/model"
	"github.com/gentleman-programming/gentle-ai/internal/system"
)

func TestDetect(t *testing.T) {
	tests := []struct {
		name            string
		lookPathPath    string
		lookPathErr     error
		stat            statResult
		wantInstalled   bool
		wantBinaryPath  string
		wantConfigPath  string
		wantConfigFound bool
		wantErr         bool
	}{
		{
			name:            "binary and config directory found",
			lookPathPath:    "/usr/local/bin/kiro",
			stat:            statResult{isDir: true},
			wantInstalled:   true,
			wantBinaryPath:  "/usr/local/bin/kiro",
			wantConfigPath:  filepath.Join("/tmp/home", ".kiro"),
			wantConfigFound: true,
		},
		{
			name:            "config dir exists but no binary — still installed (desktop app)",
			lookPathErr:     errors.New("missing"),
			stat:            statResult{isDir: true},
			wantInstalled:   true,
			wantBinaryPath:  "",
			wantConfigPath:  filepath.Join("/tmp/home", ".kiro"),
			wantConfigFound: true,
		},
		{
			name:            "binary missing and config missing",
			lookPathErr:     errors.New("missing"),
			stat:            statResult{err: os.ErrNotExist},
			wantInstalled:   false,
			wantBinaryPath:  "",
			wantConfigPath:  filepath.Join("/tmp/home", ".kiro"),
			wantConfigFound: false,
		},
		{
			name:           "stat error bubbles up",
			lookPathPath:   "/usr/local/bin/kiro",
			stat:           statResult{err: errors.New("permission denied")},
			wantConfigPath: "",
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Adapter{
				lookPath: func(string) (string, error) {
					return tt.lookPathPath, tt.lookPathErr
				},
				statPath: func(string) statResult {
					return tt.stat
				},
			}

			installed, binaryPath, configPath, configFound, err := a.Detect(context.Background(), "/tmp/home")
			if (err != nil) != tt.wantErr {
				t.Fatalf("Detect() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr {
				return
			}

			if installed != tt.wantInstalled {
				t.Fatalf("Detect() installed = %v, want %v", installed, tt.wantInstalled)
			}

			if binaryPath != tt.wantBinaryPath {
				t.Fatalf("Detect() binaryPath = %q, want %q", binaryPath, tt.wantBinaryPath)
			}

			if configPath != tt.wantConfigPath {
				t.Fatalf("Detect() configPath = %q, want %q", configPath, tt.wantConfigPath)
			}

			if configFound != tt.wantConfigFound {
				t.Fatalf("Detect() configFound = %v, want %v", configFound, tt.wantConfigFound)
			}
		})
	}
}

func TestConfigPathsCrossPlatform(t *testing.T) {
	a := NewAdapter()
	home := "/tmp/home"

	if got := a.GlobalConfigDir(home); got != filepath.Join(home, ".kiro") {
		t.Fatalf("GlobalConfigDir() = %q, want %q", got, filepath.Join(home, ".kiro"))
	}

	if got := a.SkillsDir(home); got != filepath.Join(home, ".kiro", "skills") {
		t.Fatalf("SkillsDir() = %q, want %q", got, filepath.Join(home, ".kiro", "skills"))
	}

	if got := a.SystemPromptDir(home); got != filepath.Join(home, ".kiro", "steering") {
		t.Fatalf("SystemPromptDir() = %q, want %q", got, filepath.Join(home, ".kiro", "steering"))
	}

	if got := a.SystemPromptFile(home); got != filepath.Join(home, ".kiro", "steering", "gentle-ai.md") {
		t.Fatalf("SystemPromptFile() = %q, want %q", got, filepath.Join(home, ".kiro", "steering", "gentle-ai.md"))
	}

	// MCP config path — always ~/.kiro/settings/mcp.json regardless of server name.
	want := filepath.Join(home, ".kiro", "settings", "mcp.json")
	if got := a.MCPConfigPath(home, "context7"); got != want {
		t.Fatalf("MCPConfigPath() = %q, want %q", got, want)
	}
	if got := a.MCPConfigPath(home, "engram"); got != want {
		t.Fatalf("MCPConfigPath(engram) = %q, want %q (server name should be ignored)", got, want)
	}
}

func TestSettingsPathUsesKiroUserProfile(t *testing.T) {
	a := NewAdapter()
	home := "/tmp/home"

	switch runtime.GOOS {
	case "darwin":
		path := a.SettingsPath(home)
		want := filepath.Join(home, "Library", "Application Support", "Kiro", "User", "settings.json")
		if path != want {
			t.Fatalf("SettingsPath() = %q, want %q", path, want)
		}
	case "windows":
		appData := filepath.Join(home, "AppData", "Roaming")
		t.Setenv("APPDATA", appData)
		path := a.SettingsPath(home)
		want := filepath.Join(appData, "Kiro", "User", "settings.json")
		if path != want {
			t.Fatalf("SettingsPath() = %q, want %q", path, want)
		}
	default:
		t.Setenv("XDG_CONFIG_HOME", filepath.Join(home, "xdg"))
		path := a.SettingsPath(home)
		want := filepath.Join(home, "xdg", "Kiro", "User", "settings.json")
		if path != want {
			t.Fatalf("SettingsPath() = %q, want %q", path, want)
		}
	}
}

func TestStrategies(t *testing.T) {
	a := NewAdapter()

	if got := a.SystemPromptStrategy(); got != model.StrategyFileReplace {
		t.Fatalf("SystemPromptStrategy() = %v, want %v", got, model.StrategyFileReplace)
	}

	if got := a.MCPStrategy(); got != model.StrategyMCPConfigFile {
		t.Fatalf("MCPStrategy() = %v, want %v", got, model.StrategyMCPConfigFile)
	}
}

func TestCapabilities(t *testing.T) {
	a := NewAdapter()

	if got := a.Agent(); got != model.AgentKiro {
		t.Fatalf("Agent() = %q, want %q", got, model.AgentKiro)
	}

	if got := a.Tier(); got != model.TierFull {
		t.Fatalf("Tier() = %q, want %q", got, model.TierFull)
	}

	if got := a.SupportsSkills(); !got {
		t.Fatal("SupportsSkills() = false, want true")
	}

	if got := a.SupportsSystemPrompt(); !got {
		t.Fatal("SupportsSystemPrompt() = false, want true")
	}

	if got := a.SupportsMCP(); !got {
		t.Fatal("SupportsMCP() = false, want true")
	}

	if got := a.SupportsOutputStyles(); got {
		t.Fatal("SupportsOutputStyles() = true, want false")
	}

	if got := a.SupportsSlashCommands(); got {
		t.Fatal("SupportsSlashCommands() = true, want false")
	}
}

func TestDesktopAppNotAutoInstallable(t *testing.T) {
	a := NewAdapter()

	if a.SupportsAutoInstall() {
		t.Fatal("Kiro should not support auto-install (desktop IDE)")
	}

	_, err := a.InstallCommand(system.PlatformProfile{})
	if err == nil {
		t.Fatal("InstallCommand() should return error for desktop IDE")
	}
}
