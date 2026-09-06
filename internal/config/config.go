package config

import (
	"flag"
	"os"
	"path/filepath"
)

type Config struct {
	DataFilePath string
}

func GetConfig() (Config, error) {

	file := flag.String("file", "moods.json", "moods saving file name")
	flag.Parse() // super important, otherwise flags won't work 

	configDir, err := os.UserConfigDir()
	if err != nil {
		return Config{}, err
	}

	appDir := filepath.Join(configDir, "mindgo")
	fullPath := filepath.Join(appDir, *file)

	err = os.MkdirAll(appDir, 0755)
	if err != nil {
		return Config{}, err
	}

	return Config{DataFilePath: fullPath}, nil
}
