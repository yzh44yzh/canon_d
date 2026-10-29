// Command check validates repository formatting, line widths, and documentation links.
package main

import (
	"bytes"
	"fmt"
	"go/format"
	"io/fs"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var inlineLink = regexp.MustCompile(`\]\(([^\s)]+)\)`)

func main() {
	failed := false
	report := func(path string, line int, message string) {
		fmt.Fprintf(os.Stderr, "%s:%d: %s\n", path, line, message)
		failed = true
	}

	err := filepath.WalkDir(".", func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if entry.IsDir() {
			hidden := strings.HasPrefix(entry.Name(), ".")
			if path != "." && (hidden || entry.Name() == "vendor") {
				return filepath.SkipDir
			}

			return nil
		}

		isGo := filepath.Ext(path) == ".go"
		docLocation := filepath.Dir(path) == "." ||
			strings.HasPrefix(filepath.ToSlash(path), "docs/")
		isDoc := filepath.Ext(path) == ".md" && docLocation
		if !isGo && !isDoc {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		if isGo {
			checkGo(path, content, report)
		} else {
			checkLinks(path, string(content), report)
		}

		return nil
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		failed = true
	}

	if failed {
		os.Exit(1)
	}
}

func checkGo(path string, content []byte, report func(string, int, string)) {
	formatted, err := format.Source(content)
	if err != nil {
		report(path, 1, err.Error())
	} else if !bytes.Equal(content, formatted) {
		report(path, 1, "run gofmt -w "+path)
	}

	for index, line := range strings.Split(string(content), "\n") {
		columns := 0

		for _, char := range line {
			if char == '\t' {
				columns += 8 - columns%8
			} else {
				columns++
			}
		}

		if columns > 100 {
			report(path, index+1, fmt.Sprintf("%d columns; maximum is 100", columns))
		}
	}
}

func checkLinks(path, content string, report func(string, int, string)) {
	for index, line := range strings.Split(content, "\n") {
		for _, match := range inlineLink.FindAllStringSubmatch(line, -1) {
			target := strings.Trim(match[1], "<>")
			link, err := url.Parse(target)
			if err != nil {
				report(path, index+1, "invalid link: "+target)

				continue
			}

			if link.IsAbs() || link.Host != "" || link.Path == "" {
				continue
			}

			local := filepath.Join(filepath.Dir(path), filepath.FromSlash(link.Path))
			if _, err := os.Stat(local); err != nil {
				report(path, index+1, fmt.Sprintf("link %q: %v", target, err))
			}
		}
	}
}
