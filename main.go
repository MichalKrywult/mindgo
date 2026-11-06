package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/MichalKrywult/mindgo/internal/app"
	"github.com/MichalKrywult/mindgo/internal/cli"
	"github.com/MichalKrywult/mindgo/internal/config"
	"github.com/MichalKrywult/mindgo/internal/storage"
	"github.com/MichalKrywult/mindgo/internal/tracker"
)

func main() {
	if err := runMain(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func runMain() error {
	configPath, err := config.BuildConfigPath()
	if err != nil {
		return fmt.Errorf("config path: %w", err)
	}

	if err := config.EnsureConfigFile(configPath); err != nil {
		return fmt.Errorf("ensure config file: %w", err)
	}

	oldConf, err := config.LoadConfigFromFile(configPath)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	parsedFlags, err := cli.ParseFlags()
	if err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return nil //show help is not an err
		}
		return fmt.Errorf("parse flags: %w", err)
	}

	dataPath, err := cli.SelectDataPath(parsedFlags, oldConf)
	if err != nil {
		return fmt.Errorf("select data path: %w", err)
	}

	newConf := config.Config{DataPath: dataPath}
	if oldConf != newConf {
		if err := config.SaveConfigToFile(configPath, newConf); err != nil {
			return fmt.Errorf("save config: %w", err)
		}
	}

	cliApp, err := buildCLI(newConf)
	if err != nil {
		return fmt.Errorf("build cli: %w", err)
	}

	if err := app.Run(cliApp); err != nil {
		return fmt.Errorf("start application: %w", err)
	}

	return nil
}

func buildCLI(conf config.Config) (*cli.CLI, error) {

	if err := os.MkdirAll(filepath.Dir(conf.DataPath), 0755); err != nil {
		return nil, err
	}

	fileStorage := storage.FileStorage{Filename: conf.DataPath}
	moodTracker, err := tracker.NewMoodTracker(&fileStorage)
	if err != nil {
		return nil, err
	}

	cliApp := cli.NewCLI(moodTracker, os.Stdin) // the tracker is a pointer, received from NewMoodTracker function

	return cliApp, nil
}
