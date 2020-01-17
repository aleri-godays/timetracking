package db

import (
	"context"
	"fmt"
	"github.com/aleri-godays/timetracking"
	log "github.com/sirupsen/logrus"
)
import "github.com/asdine/storm/v3"

type stormDB struct {
	db *storm.DB
}

func NewStormDB(dbPath string) *storm.DB {
	p := fmt.Sprintf("%s/timetracking.db", dbPath)
	db, err := storm.Open(p)
	if err != nil {
		log.WithFields(log.Fields{
			"db_path": p,
			"error":   err,
		}).Fatal("could not open storm db")
	}
	return db
}

func NewStormRepository(db *storm.DB) timetracking.Repository {
	if err := db.Init(&StormEntry{}); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("could not initialize StormEntry bucket")
	}

	sdb := &stormDB{
		db: db,
	}

	return sdb
}

func (s *stormDB) Add(ctx context.Context, e *timetracking.Entry) (*timetracking.Entry, error) {
	sp := entryToStormEntry(e)
	if err := s.db.Save(sp); err != nil {
		return nil, fmt.Errorf("could not save entry: %w", err)
	}
	return stormEntryToEntry(sp), nil
}

func (s *stormDB) Get(ctx context.Context, id int) (*timetracking.Entry, error) {
	var sp StormEntry
	if err := s.db.One("ID", id, &sp); err != nil {
		if err == storm.ErrNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("could not fetch entry '%d': %w", id, err)
	}
	return stormEntryToEntry(&sp), nil
}

func (s *stormDB) Delete(ctx context.Context, id int) error {
	sp := StormEntry{ID: id}
	if err := s.db.DeleteStruct(&sp); err != nil {
		return fmt.Errorf("could not delete entry '%d': %w", id, err)
	}
	return nil
}

func (s *stormDB) Update(ctx context.Context, e *timetracking.Entry) error {
	sp := entryToStormEntry(e)
	if err := s.db.Update(sp); err != nil {
		return fmt.Errorf("could not update entry '%d': %w", e.ID, err)
	}
	return nil
}

func (s *stormDB) All(ctx context.Context) ([]*timetracking.Entry, error) {
	var sps []StormEntry
	if err := s.db.All(&sps); err != nil {
		return nil, fmt.Errorf("could not fetch all entries: %w", err)
	}
	if len(sps) == 0 {
		return nil, nil
	}
	ps := make([]*timetracking.Entry, 0, len(sps))
	for _, sp := range sps {
		p := stormEntryToEntry(&sp)
		ps = append(ps, p)
	}

	return ps, nil
}
