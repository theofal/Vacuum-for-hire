package main

import (
	_ "go.uber.org/zap"
	_ "go.uber.org/zap/zapcore"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	Logger = InitLogger()
	// Do stuff BEFORE the tests!
	exitVal := m.Run()
	// Do stuff AFTER the tests!

	os.Exit(exitVal)
}

// TestGetDotEnvVar
func _(t *testing.T) {
	result := getDotEnvVar("TEST_ARG")
	if result != "BrUh" {
		t.Errorf("FAIL: Expected %v, got %v", "BrUh", result)
	}
}

func TestParseDate(t *testing.T) {
	expected := strconv.Itoa(time.Now().Add(-time.Hour * (time.Duration(2) * 24)).Day())
	result := ParseDate("Il y a 2 jours")
	if !strings.Contains(result, expected) {
		t.Errorf("FAIL : %v is not present in %v", expected, result)
	}
}
