package history

import (
	"bytes"
	"encoding/csv"
	"math"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func encodeHistory(t *testing.T, rows ...[]string) string {
	t.Helper()

	var output bytes.Buffer
	writer := csv.NewWriter(&output)
	if err := writer.WriteAll(append([][]string{columns}, rows...)); err != nil {
		t.Fatal(err)
	}

	return output.String()
}

func TestHistoryFieldValidationPreservesFileAndReleasesLock(t *testing.T) {
	tests := []struct {
		name, value, reason string
		field               int
	}{
		{name: "empty deck", field: 0, reason: "non-whitespace"},
		{name: "blank header", field: 1, value: " \u2003", reason: "non-whitespace"},
		{name: "invalid UTF-8 deck", field: 0, value: "\xff.md", reason: "UTF-8"},
		{name: "invalid UTF-8 header", field: 1, value: "\xff", reason: "UTF-8"},
		{name: "deck newline", field: 0, value: "a\nb", reason: "control"},
		{name: "header newline", field: 1, value: "a\nb", reason: "control"},
		{name: "header tab", field: 1, value: "a\tb", reason: "control"},
		{name: "header NUL", field: 1, value: "a\x00b", reason: "control"},
		{name: "header C1 control", field: 1, value: "a\u0085b", reason: "control"},
		{name: "slash", field: 0, value: "a/b", reason: "filename"},
		{name: "backslash", field: 0, value: "a\\b", reason: "filename"},
		{name: "dot", field: 0, value: ".", reason: "filename"},
		{name: "dot dot", field: 0, value: "..", reason: "filename"},
		{name: "empty counter", field: 2, reason: "invalid syntax"},
		{name: "negative", field: 3, value: "-1", reason: "invalid syntax"},
		{name: "plus", field: 4, value: "+1", reason: "invalid syntax"},
		{name: "counter space", field: 2, value: " 1", reason: "invalid syntax"},
		{name: "counter trailing space", field: 2, value: "1 ", reason: "invalid syntax"},
		{name: "fraction", field: 2, value: "1.5", reason: "invalid syntax"},
		{name: "hex", field: 2, value: "0x10", reason: "invalid syntax"},
		{name: "overflow", field: 2, value: "18446744073709551616", reason: "out of range"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			path := filepath.Join(dir, ".canon_d.csv")
			row := []string{"a.md", "Card", "1", "0", "0"}
			row[tt.field] = tt.value
			content := encodeHistory(t, row)

			if err := os.WriteFile(path, []byte(content), 0600); err != nil {
				t.Fatal(err)
			}

			store, err := Open(dir)
			if err == nil {
				store.Close()
				t.Fatal("invalid history accepted")
			}

			for _, text := range []string{path + ":2:", columns[tt.field], tt.reason} {
				if !strings.Contains(err.Error(), text) {
					t.Fatalf("error %q does not include %q", err, text)
				}
			}

			got, err := os.ReadFile(path)
			if err != nil || string(got) != content {
				t.Fatal("invalid history changed")
			}

			if err := os.WriteFile(path, []byte(encodeHistory(t)), 0600); err != nil {
				t.Fatal(err)
			}

			openStore(t, dir)
		})
	}
}

func TestHistoryCSVAcceptedForms(t *testing.T) {
	tests := []struct {
		name, content string
		want          statistics
	}{
		{name: "missing", want: statistics{}},
		{name: "header only", content: encodeHistory(t), want: statistics{}},
		{
			name: "CSV quoting, CRLF, blanks, leading zeros, no final newline",
			content: "\r\n" + strings.Join(columns, ",") + "\r\n\r\n" +
				"\"old,\"\"deck\"\"\",\"Go / Child, \"\"quoted\"\"\",001,002,000",
			want: statistics{
				{Deck: "old,\"deck\"", Header: "Go / Child, \"quoted\""}: {
					Shows: 1, CorrectAnswers: 2,
				},
			},
		},
		{
			name: "case and whitespace are significant",
			content: encodeHistory(t,
				[]string{"a", "Card", "0", "0", "0"},
				[]string{"A", "Card", "1", "0", "0"},
				[]string{"a", " Card ", "2", "0", "0"}),
			want: statistics{
				{Deck: "a", Header: "Card"}:   {},
				{Deck: "A", Header: "Card"}:   {Shows: 1},
				{Deck: "a", Header: " Card "}: {Shows: 2},
			},
		},
		{
			name: "maximum values and unlimited answers per show",
			content: encodeHistory(t, []string{
				"a.md", "Card", "1", "18446744073709551615", "18446744073709551615",
			}),
			want: statistics{
				{Deck: "a.md", Header: "Card"}: {
					Shows:            1,
					CorrectAnswers:   math.MaxUint64,
					IncorrectAnswers: math.MaxUint64,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), ".canon_d.csv")
			if tt.content != "" {
				if err := os.WriteFile(path, []byte(tt.content), 0600); err != nil {
					t.Fatal(err)
				}
			}

			got, err := readCSV(path)
			if err != nil {
				t.Fatal(err)
			}

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Fatalf("statistics mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestHistoryDiagnostics(t *testing.T) {
	header := strings.Join(columns, ",") + "\n"
	tests := []struct {
		name, content, location, reason string
	}{
		{name: "empty", location: ":1:", reason: "header"},
		{name: "blank", content: "\n\n", location: ":1:", reason: "header"},
		{name: "bad header", content: "\nwrong\n", location: ":2:", reason: "header"},
		{
			name: "zero shows correct", content: header + "\na.md,Card,0,1,0\n",
			location: ":3:",
			reason:   "invalid correct_answers: positive answers require shows > 0",
		},
		{
			name: "zero shows incorrect", content: header + "a.md,Card,0,0,1\n",
			location: ":2:",
			reason:   "invalid incorrect_answers: positive answers require shows > 0",
		},
		{
			name:     "duplicate physical lines",
			content:  header + "\na.md,Card,1,0,0\n\na.md,Card,1,0,0\n",
			location: ":5:", reason: "duplicate card identity (first occurrence: ",
		},
		{
			name: "short record", content: header + "\n\na.md,Card,1,0\n",
			location: ":4:", reason: "wrong number of fields",
		},
		{
			name: "malformed quoting", content: header + "\na.md,Car\"d,1,0,0\n",
			location: ":3:", reason: "bare",
		},
		{
			name:     "multiline identity",
			content:  header + "\na.md,\"Card\nother\",1,0,0\n",
			location: ":3:", reason: "invalid card_header: control",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), ".canon_d.csv")
			if err := os.WriteFile(path, []byte(tt.content), 0600); err != nil {
				t.Fatal(err)
			}

			_, err := readCSV(path)
			if err == nil {
				t.Fatal("invalid history accepted")
			}

			for _, text := range []string{path + tt.location, tt.reason} {
				if !strings.Contains(err.Error(), text) {
					t.Fatalf("error %q does not contain %q", err, text)
				}
			}

			if tt.name == "duplicate physical lines" &&
				!strings.Contains(err.Error(), path+":3)") {
				t.Fatalf("missing first duplicate location: %v", err)
			}
		})
	}
}
