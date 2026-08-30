package main

import (
	"testing"
	"time"
)

func TestBuildServer(t *testing.T) {
	srv := buildServer("9999")
	if srv.Addr != ":9999" {
		t.Errorf("Addr = %q, want %q", srv.Addr, ":9999")
	}
	if srv.ReadTimeout == 0 {
		t.Error("ReadTimeout must be non-zero")
	}
	if srv.WriteTimeout == 0 {
		t.Error("WriteTimeout must be non-zero")
	}
	if srv.IdleTimeout == 0 {
		t.Error("IdleTimeout must be non-zero")
	}
	if srv.Handler == nil {
		t.Error("Handler must not be nil")
	}
	if srv.ReadTimeout != 10*time.Second {
		t.Errorf("ReadTimeout = %v, want 10s", srv.ReadTimeout)
	}
}

func TestBuildServer_DefaultPortPattern(t *testing.T) {
	srv := buildServer("8080")
	if srv.Addr != ":8080" {
		t.Errorf("Addr = %q, want %q", srv.Addr, ":8080")
	}
}
