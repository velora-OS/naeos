package errors

import (
	"errors"
	"testing"
)

func TestNew(t *testing.T) {
	err := New(ErrParse, "bad input")
	if err.Code != ErrParse {
		t.Fatalf("expected code %s, got %s", ErrParse, err.Code)
	}
	if err.Message != "bad input" {
		t.Fatalf("expected message %q, got %q", "bad input", err.Message)
	}
	if err.Inner != nil {
		t.Fatal("expected nil Inner")
	}
}

func TestErrorFormat(t *testing.T) {
	err := New(ErrValidation, "missing field")
	got := err.Error()
	want := "VALIDATION_ERROR: missing field"
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestWrapAndUnwrap(t *testing.T) {
	inner := errors.New("underlying problem")
	err := Wrap(ErrPlugin, "load failed", inner)

	if err.Inner != inner {
		t.Fatal("Inner does not match")
	}

	unwrapped := err.Unwrap()
	if unwrapped != inner {
		t.Fatalf("Unwrap returned %v, want %v", unwrapped, inner)
	}

	if !errors.Is(err, inner) {
		t.Fatal("errors.Is should find inner error")
	}
}

func TestIsCode(t *testing.T) {
	err := New(ErrCloud, "timeout")

	if !Is(err, ErrCloud) {
		t.Fatal("Is should match same code")
	}
	if Is(err, ErrParse) {
		t.Fatal("Is should not match different code")
	}
}

func TestIsChained(t *testing.T) {
	inner := New(ErrNetwork, "connection refused")
	outer := Wrap(ErrPlugin, "plugin init failed", inner)

	if !Is(outer, ErrPlugin) {
		t.Fatal("Is should match outer code")
	}
	if !Is(outer, ErrNetwork) {
		t.Fatal("Is should match inner code via chain")
	}
}

func TestIsNonNaeosError(t *testing.T) {
	err := errors.New("plain error")
	if Is(err, ErrInternal) {
		t.Fatal("Is should return false for non-NaeosError")
	}
}

func TestSentinels(t *testing.T) {
	tests := []struct {
		sentinel *NaeosError
		code     ErrorCode
	}{
		{ErrNotConnected, ErrNetwork},
		{ErrInvalidSpec, ErrValidation},
		{ErrPluginNotFound, ErrPlugin},
		{ErrDeployFailed, ErrCloud},
	}

	for _, tt := range tests {
		if tt.sentinel.Code != tt.code {
			t.Fatalf("%s: expected code %s, got %s", tt.sentinel.Message, tt.code, tt.sentinel.Code)
		}
		if !Is(tt.sentinel, tt.code) {
			t.Fatalf("Is should match sentinel %s", tt.sentinel.Message)
		}
	}
}
