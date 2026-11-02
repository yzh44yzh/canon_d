package history

import (
	"cmp"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strconv"

	"com.github.yzh44yzh/canon_d/internal/model"
)

var columns = []string{"deck", "card_header", "shows", "correct_answers", "incorrect_answers"}

func readCSV(path string) (statistics, error) {
	file, err := os.Open(path)
	if errors.Is(err, os.ErrNotExist) {
		return make(statistics), nil
	}

	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	header, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("%s: read CSV header: %w", path, err)
	}

	if !slices.Equal(header, columns) {
		return nil, fmt.Errorf("%s: unexpected CSV header", path)
	}

	results := make(statistics)

	for row := 2; ; row++ {
		fields, err := reader.Read()
		if errors.Is(err, io.EOF) {
			return results, nil
		}

		if err != nil {
			return nil, fmt.Errorf("%s: row %d: %w", path, row, err)
		}

		id, total, err := parseRow(fields)
		if err != nil {
			return nil, fmt.Errorf("%s: row %d: %w", path, row, err)
		}

		if _, exists := results[id]; exists {
			return nil, fmt.Errorf("%s: row %d: duplicate card identity", path, row)
		}

		results[id] = total
	}
}

func parseRow(fields []string) (model.CardIdentity, model.CardStatistics, error) {
	id := model.CardIdentity{Deck: fields[0], Header: fields[1]}
	var counts [3]uint64

	for i := range counts {
		value, err := strconv.ParseUint(fields[i+2], 10, 64)
		if err != nil {
			reason := fmt.Errorf("invalid %s: %w", columns[i+2], err)

			return id, model.CardStatistics{}, reason
		}

		counts[i] = value
	}

	return id, model.CardStatistics{
		Shows: counts[0], CorrectAnswers: counts[1], IncorrectAnswers: counts[2],
	}, nil
}

func writeCSV(path string, results statistics) error {
	file, err := os.CreateTemp(filepath.Dir(path), ".canon_d-*.tmp")
	if err != nil {
		return err
	}
	defer os.Remove(file.Name())
	defer file.Close()

	writer := csv.NewWriter(file)
	if err := writer.Write(columns); err != nil {
		return err
	}

	ids := slices.SortedFunc(maps.Keys(results), compareIdentities)

	for _, id := range ids {
		total := results[id]
		fields := []string{
			id.Deck, id.Header,
			strconv.FormatUint(total.Shows, 10),
			strconv.FormatUint(total.CorrectAnswers, 10),
			strconv.FormatUint(total.IncorrectAnswers, 10),
		}
		if err := writer.Write(fields); err != nil {
			return err
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return err
	}

	if err := file.Sync(); err != nil {
		return err
	}

	if err := file.Close(); err != nil {
		return err
	}

	return os.Rename(file.Name(), path)
}

func compareIdentities(a, b model.CardIdentity) int {
	if deck := cmp.Compare(a.Deck, b.Deck); deck != 0 {
		return deck
	}

	return cmp.Compare(a.Header, b.Header)
}
