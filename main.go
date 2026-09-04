package main

import (
	"fmt"
	"os"

	"github.com/MichalKrywult/mindgo/internal/cli"
	"github.com/MichalKrywult/mindgo/internal/config"
	"github.com/MichalKrywult/mindgo/internal/storage"
	"github.com/MichalKrywult/mindgo/internal/tracker"
)

func main() {

	cfg, err := config.GetDefaultConfig()
	if err != nil {
		fmt.Printf("Failed to load configuration: %v", err)
	}

	storage := storage.FileStorage{Filename: cfg.DataFilePath}
	tracker, err := tracker.NewMoodTracker(&storage)
	if err != nil {
		fmt.Println("Error initializing tracker:", err)
		return
	}

	cli := cli.NewCLI(tracker, os.Stdin) // the tracker is a pointer, received from NewMoodTracker function
	cli.Show()
}
