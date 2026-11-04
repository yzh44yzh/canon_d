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
	"unicode"
	"unicode/utf8"

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

type fieldError struct {
	index int
	err   error
}

func (e *fieldError) Error() string {
	return fmt.Sprintf("invalid %s: %v", columns[e.index], e.err)
}

func parseRow(fields []string) (model.CardIdentity, model.CardStatistics, error) {
	id := model.CardIdentity{Deck: fields[0], Header: fields[1]}
	var counts [3]uint64

	for i := 0; i < 2; i++ {
		if err := validateIdentityField(fields[i], i == 0); err != nil {
			return id, model.CardStatistics{}, &fieldError{index: i, err: err}
		}
	}

	for i := range counts {
		value, err := strconv.ParseUint(fields[i+2], 10, 64)
		if err != nil {
			reason := &fieldError{index: i + 2, err: err}

			return id, model.CardStatistics{}, reason
		}

		counts[i] = value
	}

	if counts[0] == 0 {
		for i := 1; i < len(counts); i++ {
			if counts[i] > 0 {
				return id, model.CardStatistics{}, &fieldError{
					index: i + 2,
					err:   errors.New("positive answers require shows > 0"),
				}
			}
		}
	}

	return id, model.CardStatistics{
		Shows: counts[0], CorrectAnswers: counts[1], IncorrectAnswers: counts[2],
	}, nil
}

func validateIdentityField(value string, deck bool) error {
	if !utf8.ValidString(value) {
		return errors.New("must contain valid UTF-8")
	}

	if strings.TrimSpace(value) == "" {
		return errors.New("must contain non-whitespace text")
	}

	if strings.ContainsFunc(value, unicode.IsControl) {
		return errors.New("control characters are not allowed")
	}

	if deck && (value == "." || value == ".." || strings.ContainsAny(value, "/\\")) {
		return errors.New("must be a filename, without directory separators")
	}

	return nil
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
