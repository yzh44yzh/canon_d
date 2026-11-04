package history

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"

	"com.github.yzh44yzh/canon_d/internal/model"
)

type fieldError struct {
	index int
	err   error
}

func (e *fieldError) Error() string {
	return fmt.Sprintf("invalid %s: %v", columns[e.index], e.err)
}

// validateRow checks typed persisted data; CSV syntax belongs to the reader.
func validateRow(id model.CardIdentity, total model.CardStatistics) error {
	if err := validateIdentity(id); err != nil {
		return err
	}

	return validateTotals(total)
}

func validateIdentity(id model.CardIdentity) error {
	for i, value := range []string{id.Deck, id.Header} {
		if err := validateIdentityField(value, i == 0); err != nil {
			return &fieldError{index: i, err: err}
		}
	}

	return nil
}

func validateTotals(total model.CardStatistics) error {
	if total.Shows > 0 {
		return nil
	}

	for i, count := range []uint64{total.CorrectAnswers, total.IncorrectAnswers} {
		if count > 0 {
			return &fieldError{
				index: i + 3,
				err:   errors.New("positive answers require shows > 0"),
			}
		}
	}

	return nil
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
