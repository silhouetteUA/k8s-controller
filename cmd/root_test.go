package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/rs/zerolog"
)

func TestParseLogLevel(t *testing.T) {
	tests := map[string]zerolog.Level{
		"trace": zerolog.TraceLevel,
		"debug": zerolog.DebugLevel,
		"info":  zerolog.InfoLevel,
		"warn":  zerolog.WarnLevel,
		"error": zerolog.ErrorLevel,
		"":      zerolog.InfoLevel, // default
		"FOO":   zerolog.InfoLevel, // invalid → default
	}

	for input, expected := range tests {
		result := parseLogLevel(input)
		if result != expected {
			t.Errorf("Expected %v for input %q, got %v", expected, input, result)
		}
	}

}

func TestRootCmd_HelpShownOnNoArgs(t *testing.T) {
	// Capture help output
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{}) // simulate: `kctl`

	err := rootCmd.Execute()
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Usage:") {
		t.Errorf("Expected usage/help text, got: %q", output)
	}
}
