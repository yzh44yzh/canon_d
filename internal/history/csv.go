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
	"strings"

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
		return nil, csvReadError(path, err)
	}

	if !slices.Equal(header, columns) {
		line, _ := reader.FieldPos(0)
		return nil, fmt.Errorf("%s:%d: expected CSV header %s",
			path, line, strings.Join(columns, ","))
	}

	results := make(statistics)
	firstLines := make(map[model.CardIdentity]int)

	for {
		fields, err := reader.Read()
		if errors.Is(err, io.EOF) {
			return results, nil
		}

		if err != nil {
			return nil, csvReadError(path, err)
		}

		line, _ := reader.FieldPos(0)
		id, total, err := parseRow(fields)
		if err != nil {
			var fieldErr *fieldError
			if errors.As(err, &fieldErr) {
				line, _ = reader.FieldPos(fieldErr.index)
			}

			return nil, fmt.Errorf("%s:%d: %w", path, line, err)
		}

		if first, exists := firstLines[id]; exists {
			return nil, fmt.Errorf(
				"%s:%d: duplicate card identity (first occurrence: %s:%d)",
				path, line, path, first)
		}

		results[id] = total
		firstLines[id] = line
	}
}

func csvReadError(path string, err error) error {
	var parseErr *csv.ParseError
	if errors.As(err, &parseErr) {
		return fmt.Errorf("%s:%d: invalid CSV: %w", path, parseErr.Line, err)
	}

	return fmt.Errorf("%s:1: read CSV header or record: %w", path, err)
}

func parseRow(fields []string) (model.CardIdentity, model.CardStatistics, error) {
	id := model.CardIdentity{Deck: fields[0], Header: fields[1]}
	var counts [3]uint64

	if err := validateIdentity(id); err != nil {
		return id, model.CardStatistics{}, err
	}

	for i := range counts {
		value, err := strconv.ParseUint(fields[i+2], 10, 64)
		if err != nil {
			reason := &fieldError{index: i + 2, err: err}

			return id, model.CardStatistics{}, reason
		}

		counts[i] = value
	}

	total := model.CardStatistics{
		Shows: counts[0], CorrectAnswers: counts[1], IncorrectAnswers: counts[2],
	}
	if err := validateTotals(total); err != nil {
		return id, model.CardStatistics{}, err
	}

	return id, total, nil
}

func writeCSV(path string, results statistics) error {
	ids := slices.SortedFunc(maps.Keys(results), compareIdentities)

	for _, id := range ids {
		if err := validateRow(id, results[id]); err != nil {
			return fmt.Errorf("%s, deck %q, card %q: %w", path, id.Deck, id.Header, err)
		}
	}

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
