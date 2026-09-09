package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/MichalKrywult/mindgo/internal/cli"
	"github.com/MichalKrywult/mindgo/internal/storage"
	"github.com/MichalKrywult/mindgo/internal/tracker"
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

	path, err := cli.BuildPath(parsedFlags)
	if err != nil {
		fmt.Println("Error with path building:", err)
		return
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		fmt.Println("Error creating data directory:", err)
		return
	}

	fileStorage := storage.FileStorage{Filename: path}
	moodTracker, err := tracker.NewMoodTracker(&fileStorage)
	if err != nil {
		fmt.Println("Error initializing tracker:", err)
		return
	}

	cli := cli.NewCLI(moodTracker, os.Stdin) // the tracker is a pointer, received from NewMoodTracker function
	cli.Show()
}
