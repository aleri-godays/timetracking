package db

import (
	"github.com/aleri-godays/timetracking"
)

type StormEntry struct {
	ID                int `storm:"id,increment"`
	DateTS            int64
	Project           int    `storm:"index"`
	User              string `storm:"index"`
	Comment           string
	DurationInSeconds int64
}

func entryToStormEntry(e *timetracking.Entry) *StormEntry {
	return &StormEntry{
		ID:                e.ID,
		DateTS:            e.DateTS,
		Project:           e.Project,
		User:              e.User,
		Comment:           e.Comment,
		DurationInSeconds: e.Duration,
	}
}

func stormEntryToEntry(e *StormEntry) *timetracking.Entry {
	return &timetracking.Entry{
		ID:       e.ID,
		DateTS:   e.DateTS,
		Project:  e.Project,
		User:     e.User,
		Comment:  e.Comment,
		Duration: e.DurationInSeconds,
	}
}
