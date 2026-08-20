package search

import (
	"bufio"
	"html"
	"strconv"
	"strings"
)

type Result struct {
	Book    string
	Chapter int
	Verse   int
	Text    string
}

func Find(data, query string) []Result {
	query = strings.ToLower(strings.TrimSpace(query))

	if query == "" {
		return nil
	}

	var results []Result

	scanner := bufio.NewScanner(strings.NewReader(data))
	scanner.Buffer(make([]byte, 64*1024), 1024*1024)

	for scanner.Scan() {
		fields := strings.SplitN(scanner.Text(), "\t", 5)
		if len(fields) != 5 {
			continue
		}

		text := html.UnescapeString(fields[4])

		if !strings.Contains(strings.ToLower(text), query) {
			continue
		}

		chapter, err := strconv.Atoi(fields[2])
		if err != nil {
			continue
		}

		verse, err := strconv.Atoi(fields[3])
		if err != nil {
			continue
		}

		results = append(results, Result{
			Book:    fields[1],
			Chapter: chapter,
			Verse:   verse,
			Text:    text,
		})
	}

	return results
}
