package config

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	cfg := Load()
	if cfg.ServerPort == "" {
		t.Error("ServerPort should not be empty")
	}
	if cfg.DSN() == "" {
		t.Error("DSN should not be empty")
	}
}

func TestDSN(t *testing.T) {
	cfg := &Config{
		DBUser:     "user",
		DBPassword: "pass",
		DBHost:     "localhost",
		DBPort:     "5432",
		DBName:     "testdb",
	}
	dsn := cfg.DSN()
	if dsn != "postgres://user:pass@localhost:5432/testdb?sslmode=disable" {
		t.Errorf("unexpected DSN: %s", dsn)
	}
}

func TestGetEnv(t *testing.T) {
	os.Setenv("TEST_VAR", "custom")
	defer os.Unsetenv("TEST_VAR")

	if v := getEnv("TEST_VAR", "default"); v != "custom" {
		t.Errorf("expected custom, got %s", v)
	}
	if v := getEnv("MISSING_VAR", "default"); v != "default" {
		t.Errorf("expected default, got %s", v)
	}
}
