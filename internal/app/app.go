package app

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/MichalKrywult/mindgo/internal/cli"
	"github.com/MichalKrywult/mindgo/internal/storage"
	"github.com/MichalKrywult/mindgo/internal/tracker"
)

type Config struct {
	DataPath string
}

func Run(config Config) error {

	if err := os.MkdirAll(filepath.Dir(config.DataPath), 0755); err != nil {
		fmt.Println("Error creating data directory:", err)
		return err
	}

	fileStorage := storage.FileStorage{Filename: config.DataPath}
	moodTracker, err := tracker.NewMoodTracker(&fileStorage)
	if err != nil {
		fmt.Println("Error initializing tracker:", err)
		return err
	}

	cli := cli.NewCLI(moodTracker, os.Stdin) // the tracker is a pointer, received from NewMoodTracker function
	cli.Show()

	return nil
}
