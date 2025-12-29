package main

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/kubaliski/golog/pkg/logger"
)

func TestCreateDefaultBot_NoToken(t *testing.T) {
	old := os.Getenv("DISCORD_TOKEN")
	defer os.Setenv("DISCORD_TOKEN", old)
	os.Unsetenv("DISCORD_TOKEN")

	// Ensure no .env present in cwd
	dir, err := os.MkdirTemp("", "maingentest")
	if err != nil {
		t.Fatalf("mktemp: %v", err)
	}
	defer os.RemoveAll(dir)

	pwd, _ := os.Getwd()
	defer os.Chdir(pwd)
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	// call createDefaultBot which should return an error due to missing token
	_, err = createDefaultBot()
	if err == nil {
		t.Fatalf("expected error when DISCORD_TOKEN not set")
	}
}

func TestCreateDefaultBot_WithEnv(t *testing.T) {
	dir, err := os.MkdirTemp("", "maingentest2")
	if err != nil {
		t.Fatalf("mktemp: %v", err)
	}
	defer os.RemoveAll(dir)

	pwd, _ := os.Getwd()
	defer os.Chdir(pwd)
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	// write .env and chdir there so godotenv.Load can find it
	if err := os.WriteFile(filepath.Join(dir, ".env"), []byte("DISCORD_TOKEN=tokx\n"), 0644); err != nil {
		t.Fatalf("write .env: %v", err)
	}

	// ensure no DISCORD_TOKEN envvar overrides
	old := os.Getenv("DISCORD_TOKEN")
	defer os.Setenv("DISCORD_TOKEN", old)
	os.Unsetenv("DISCORD_TOKEN")

	if err := os.Chdir(dir); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	base := context.Background()
	base = logger.SetServiceName(base, "test")
	lg := logger.FromContext(base)

	// Call createDefaultBot; it should load .env and create a bot
	b, err := createDefaultBot()
	if err != nil {
		t.Fatalf("unexpected error creating bot: %v", err)
	}
	if b == nil {
		t.Fatalf("expected non-nil bot")
	}
	_ = lg
}
