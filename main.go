package main

import (
	"fmt"
	"os"

	"github.com/MichalKrywult/mindgo/internal/cli"
	"github.com/MichalKrywult/mindgo/internal/storage"
	"github.com/MichalKrywult/mindgo/internal/tracker"
)

func main() {

	storage := storage.FileStorage{Filename: "moods.json"}
	tracker, err := tracker.NewMoodTracker(&storage)
	if err != nil {
		fmt.Println("Error initializing tracker:", err)
		return
	}

	cli := cli.NewCLI(tracker, os.Stdin) // the tracker is a pointer, received from NewMoodTracker function
	cli.Show()
}
