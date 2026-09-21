package main

import (
	"flag"
	"fmt"

	"github.com/MichalKrywult/mindgo/internal/app"
	"github.com/MichalKrywult/mindgo/internal/cli"
	"github.com/MichalKrywult/mindgo/internal/config"
)

func main() {

	configPath, err := config.BuildConfigPath()
	if err != nil {
		fmt.Println("Error with ConfigFile path:", err)
		return
	}

	err = config.EnsureConfigFile(configPath)
	if err != nil {
		fmt.Println("Error with ConfigFile:", err)
		return
	}

	conf, err := config.LoadConfigFromFile(configPath)
	if err != nil {
		fmt.Println("Error with loading from ConfigFile:", err)
		return
	}

	parsedFlags, err := cli.ParseFlags()
	if err != nil {
		if err == flag.ErrHelp {
			return
		}

		fmt.Println("Error with parsing flags:", err)
		return
	}

	dataPath, err := cli.SelectDataPath(parsedFlags, conf)
	if err != nil {
		fmt.Println("Error with path building:", err)
		return
	}

	err = app.Run(config.Config{DataPath: dataPath})
	if err != nil {
		fmt.Println("Error starting apllication:", err)
		return
	}
}
