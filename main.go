package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/MichalKrywult/mindgo/internal/app"
	"github.com/MichalKrywult/mindgo/internal/cli"
	"github.com/MichalKrywult/mindgo/internal/config"
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

	if err := app.Run(newConf); err != nil {
		return fmt.Errorf("start application: %w", err)
	}

	return nil
}