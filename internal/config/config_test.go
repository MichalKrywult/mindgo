package config_test

import (
	"strings"
	"testing"

	"github.com/MichalKrywult/mindgo/internal/config"
)

func TestGetDefaultConfig(t *testing.T) {

	config, err := config.GetDefaultConfig()
	if err != nil {
		t.Fatal("unexpected error: ", err)
	}

	if len(config.DataFilePath) == 0 {
		t.Fatal("filepath is empty")
	}

	if !strings.HasSuffix(config.DataFilePath, "moods.json") {
		t.Fatalf("filepath '%s' does not end with 'moods.json'", config.DataFilePath)
	}
}
