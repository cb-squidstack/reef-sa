package logger

import (
	"testing"
)

func TestSetLevel(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"debug level", "debug", "debug"},
		{"info level", "info", "info"},
		{"warn level", "warn", "warn"},
		{"error level", "error", "error"},
		{"default to info", "invalid", "info"},
		{"empty defaults to info", "", "info"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetLevel(tt.input)
			got := GetLevel()
			if got != tt.expected {
				t.Errorf("SetLevel(%q): expected %q, got %q", tt.input, tt.expected, got)
			}
		})
	}
}

func TestInit(t *testing.T) {
	// Test that Init doesn't panic
	Init("info")
	if GetLevel() != "info" {
		t.Errorf("expected level 'info' after Init, got %q", GetLevel())
	}

	Init("debug")
	if GetLevel() != "debug" {
		t.Errorf("expected level 'debug' after Init, got %q", GetLevel())
	}
}

func TestLevelPrecedence(t *testing.T) {
	// Set to error (highest threshold)
	SetLevel("error")
	if GetLevel() != "error" {
		t.Errorf("expected 'error', got %q", GetLevel())
	}

	// Set to debug (lowest threshold)
	SetLevel("debug")
	if GetLevel() != "debug" {
		t.Errorf("expected 'debug', got %q", GetLevel())
	}

	// Set to warn (middle threshold)
	SetLevel("warn")
	if GetLevel() != "warn" {
		t.Errorf("expected 'warn', got %q", GetLevel())
	}
}

func TestLogFunctions(t *testing.T) {
	// Set level to debug to enable all log functions
	SetLevel("debug")

	// Test that log functions don't panic
	t.Run("Debugf", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Debugf panicked: %v", r)
			}
		}()
		Debugf("test debug message")
	})

	t.Run("Infof", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Infof panicked: %v", r)
			}
		}()
		Infof("test info message")
	})

	t.Run("Warnf", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Warnf panicked: %v", r)
			}
		}()
		Warnf("test warn message")
	})

	t.Run("Errorf", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Errorf panicked: %v", r)
			}
		}()
		Errorf("test error message")
	})
}

func TestLogLevelFiltering(t *testing.T) {
	// When level is set to error, only error logs should be visible
	// We can't easily test output, but we can verify it doesn't panic
	SetLevel("error")

	Debugf("this should not be logged")
	Infof("this should not be logged")
	Warnf("this should not be logged")
	Errorf("this should be logged")

	// No panic means test passes
}

func TestConcurrentLevelChange(t *testing.T) {
	// Test that concurrent level changes don't cause data races
	done := make(chan bool)

	for i := 0; i < 10; i++ {
		go func() {
			SetLevel("debug")
			_ = GetLevel()
			SetLevel("info")
			_ = GetLevel()
			done <- true
		}()
	}

	for i := 0; i < 10; i++ {
		<-done
	}
}
