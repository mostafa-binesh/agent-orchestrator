package gitworktree

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestRunCommandCancelsTransportChildren(t *testing.T) {
	script := filepath.Join(t.TempDir(), "fake-git")
	if err := os.WriteFile(script, []byte("#!/bin/sh\nsleep 10 &\nwait\n"), 0o755); err != nil {
		t.Fatalf("write fake git: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	started := time.Now()
	if _, err := runCommand(ctx, script, "-C", "/tmp", "fetch", "origin"); err == nil {
		t.Fatal("runCommand succeeded, want a timeout error")
	}
	if elapsed := time.Since(started); elapsed > 5*time.Second {
		t.Fatalf("runCommand returned after %s; a killed child kept the output pipe open", elapsed)
	}
}
