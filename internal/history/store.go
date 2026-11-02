// Package history persists cumulative card statistics in a deck folder.
package history

import (
	"errors"
	"fmt"
	"maps"
	"math"
	"os"
	"path/filepath"

	"com.github.yzh44yzh/canon_d/internal/model"
)

type statistics = map[model.CardIdentity]model.CardStatistics

// Store holds an exclusive folder lock and the totals loaded before practice.
type Store struct {
	path  string
	lock  *os.File
	base  statistics
	saved statistics
}

func Open(folder string) (*Store, error) {
	path := filepath.Join(folder, ".canon_d.csv")
	lock, err := os.OpenFile(path+".lock", os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return nil, fmt.Errorf("lock statistics: %w", err)
	}

	if err := lockFile(lock); err != nil {
		lock.Close()

		return nil, fmt.Errorf("lock %s: another run may be active: %w", path, err)
	}

	base, err := readCSV(path)
	if err != nil {
		lock.Close()

		return nil, fmt.Errorf("load statistics: %w", err)
	}

	return &Store{path: path, lock: lock, base: base}, nil
}

// Close releases the lock. Its file remains in place to keep one shared lock inode.
func (s *Store) Close() error { return s.lock.Close() }

// InitialStatistics returns a copy of the saved totals loaded when the store opened.
// It excludes current-run results, even after checkpoints have been saved.
func (s *Store) InitialStatistics() statistics { return maps.Clone(s.base) }

// Save merges cumulative results from this run with the original stored totals.
// Repeated saves never add the same submissions or shows twice.
func (s *Store) Save(results statistics) error {
	if maps.Equal(results, s.saved) {
		return nil
	}

	merged := maps.Clone(s.base)

	for id, result := range results {
		total, err := addStatistics(merged[id], result)
		if err != nil {
			return fmt.Errorf("save %s, card %q: %w", s.path, id.Header, err)
		}

		merged[id] = total
	}

	if err := writeCSV(s.path, merged); err != nil {
		return fmt.Errorf("save statistics: %w", err)
	}

	s.saved = maps.Clone(results)

	return nil
}

func addStatistics(a, b model.CardStatistics) (model.CardStatistics, error) {
	if b.Shows > math.MaxUint64-a.Shows ||
		b.CorrectAnswers > math.MaxUint64-a.CorrectAnswers ||
		b.IncorrectAnswers > math.MaxUint64-a.IncorrectAnswers {
		return model.CardStatistics{}, errors.New("statistics counter overflow")
	}

	return model.CardStatistics{
		Shows:            a.Shows + b.Shows,
		CorrectAnswers:   a.CorrectAnswers + b.CorrectAnswers,
		IncorrectAnswers: a.IncorrectAnswers + b.IncorrectAnswers,
	}, nil
}
