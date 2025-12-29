package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad_Success(t *testing.T) {
	// backup and restore DISCORD_TOKEN
	old := os.Getenv("DISCORD_TOKEN")
	defer os.Setenv("DISCORD_TOKEN", old)
	os.Unsetenv("DISCORD_TOKEN")

	dir, err := os.MkdirTemp("", "cfgtest")
	if err != nil {
		t.Fatalf("mktemp: %v", err)
	}
	defer os.RemoveAll(dir)

	pwd, _ := os.Getwd()
	defer os.Chdir(pwd)
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	// write .env with DISCORD_TOKEN
	envPath := filepath.Join(dir, ".env")
	if err := os.WriteFile(envPath, []byte("DISCORD_TOKEN=tok123\n"), 0644); err != nil {
		t.Fatalf("write .env: %v", err)
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if cfg.DiscordToken != "tok123" {
		t.Fatalf("expected token tok123, got %q", cfg.DiscordToken)
	}
}

func TestLoad_NoDotEnv_ReturnsError(t *testing.T) {
	old := os.Getenv("DISCORD_TOKEN")
	defer os.Setenv("DISCORD_TOKEN", old)
	os.Unsetenv("DISCORD_TOKEN")

	dir, err := os.MkdirTemp("", "cfgtest2")
	if err != nil {
		t.Fatalf("mktemp: %v", err)
	}
	defer os.RemoveAll(dir)

	pwd, _ := os.Getwd()
	defer os.Chdir(pwd)
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	if _, err := Load(); err == nil {
		t.Fatalf("expected error when .env missing and DISCORD_TOKEN unset")
	}
}

func TestLoad_NoTokenInDotEnv_ReturnsError(t *testing.T) {
	old := os.Getenv("DISCORD_TOKEN")
	defer os.Setenv("DISCORD_TOKEN", old)
	os.Unsetenv("DISCORD_TOKEN")

	dir, err := os.MkdirTemp("", "cfgtest3")
	if err != nil {
		t.Fatalf("mktemp: %v", err)
	}
	defer os.RemoveAll(dir)

	pwd, _ := os.Getwd()
	defer os.Chdir(pwd)
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	if err := os.WriteFile(filepath.Join(dir, ".env"), []byte("FOO=bar\n"), 0644); err != nil {
		t.Fatalf("write .env: %v", err)
	}

	if _, err := Load(); err == nil {
		t.Fatalf("expected error when DISCORD_TOKEN not set in .env")
	}
}
