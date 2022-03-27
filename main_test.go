package main

import (
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestGetDotEnvVar(t *testing.T) {
	InitLogger()

	result := getDotEnvVar("TEST_ARG")
	if result != "BrUh" {
		t.Errorf("FAIL: Expected %v, got %v", "BrUh", result)
	}
}

func TestParseDate(t *testing.T) {
	InitLogger()

	expected := strconv.Itoa(time.Now().Add(-time.Hour * (time.Duration(2) * 24)).Day())
	result := ParseDate("Il y a 2 jours")
	if !strings.Contains(result, expected) {
		t.Errorf("FAIL : %v is not present in %v", expected, result)
	}
}
