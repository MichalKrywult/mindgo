package tracker

import (
	"fmt"

	"github.com/MichalKrywult/mindgo/internal/domain"
	"github.com/MichalKrywult/mindgo/internal/storage"
)

func (tracker *MoodTracker) save() error {
	if err := tracker.Save(tracker.entries); err != nil {
		return fmt.Errorf("unexpected error occurred when saving to file: %w", err)
	}
	return nil
}

type MoodTracker struct {
	entries []domain.MoodEntry
	entryID int
	storage.Storage
}

type backupMoodTracker struct {
	entries []domain.MoodEntry
	entryID int
}

func (tracker *MoodTracker) createTrackerBackup() backupMoodTracker {
	backupEntries := make([]domain.MoodEntry, len(tracker.entries))
	copy(backupEntries, tracker.entries)

	backupEntryID := tracker.entryID
	return backupMoodTracker{entries: backupEntries, entryID: backupEntryID}
}

func (tracker *MoodTracker) restoreTrackerFromBackup(backup backupMoodTracker) {
	tracker.entries = backup.entries
	tracker.entryID = backup.entryID
}

func NewMoodTracker(storage storage.Storage) (*MoodTracker, error) {
	// read data from the file
	entries, err := storage.Load()
	if err != nil {
		return nil, err
	}

	// search for max id in data that is being loaded
	maxID := 0
	for _, e := range entries {
		if e.ID > maxID {
			maxID = e.ID
		}
	}

	// return pointer to the created struct
	return &MoodTracker{
		entries: entries,
		entryID: maxID,
		Storage: storage,
	}, nil
}

func (tracker *MoodTracker) validateIndex(index int) error {
	if index < 0 || index >= len(tracker.entries) {
		return fmt.Errorf("invalid entry number")
	}
	return nil
}

func (tracker *MoodTracker) GetEntries() []domain.MoodEntry {
	entries := make([]domain.MoodEntry, len(tracker.entries))
	//we have to create space for that copy with make()
	// copy() doesn't make a slice bigger
	copy(entries, tracker.entries)
	//now we can copy tracker.entries into entries

	return entries
}

func (tracker *MoodTracker) AddEntry(entry domain.MoodEntry) error {

	backup := tracker.createTrackerBackup()

	tracker.entryID++
	entry.ID = tracker.entryID

	tracker.entries = append(tracker.entries, entry)
	err := tracker.save()
	if err != nil {
		tracker.restoreTrackerFromBackup(backup)
		return err
	}

	return nil
}

func (tracker *MoodTracker) EditEntryByIndex(index int, newEntry domain.MoodEntry) error {
	if err := tracker.validateIndex(index); err != nil {
		return err
	}

	backup := tracker.createTrackerBackup()

	oldEntry := tracker.entries[index]
	newEntry.ID = oldEntry.ID
	tracker.entries[index] = newEntry

	if err := tracker.save(); err != nil {
		tracker.restoreTrackerFromBackup(backup)
		return err
	}

	return nil
}

func (tracker *MoodTracker) RemoveEntryByIndex(index int) error {
	if err := tracker.validateIndex(index); err != nil {
		return err
	}

	backup := tracker.createTrackerBackup()

	tracker.entries = append(
		tracker.entries[:index],
		tracker.entries[index+1:]...,
	)

	if err := tracker.save(); err != nil {
		tracker.restoreTrackerFromBackup(backup)
		return err
	}

	return nil
}
