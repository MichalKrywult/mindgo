package main

import (
	"flag"
	"fmt"

	"github.com/MichalKrywult/mindgo/internal/app"
	"github.com/MichalKrywult/mindgo/internal/cli"
)

func main() {

	parsedFlags, err := cli.ParseFlags()
	if err != nil {
		if err == flag.ErrHelp {
			return
		}

		fmt.Println("Error with parsing flags:", err)
		return
	}

	dataPath, err := cli.BuildPath(parsedFlags)
	if err != nil {
		fmt.Println("Error with path building:", err)
		return
	}

	config := app.Config{DataPath: dataPath}
	err = app.Run(config)
	if err != nil {
		fmt.Println("Error with path building:", err)
		return
	}
}
