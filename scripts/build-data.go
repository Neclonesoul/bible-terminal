package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

type FlexibleNumber string

func (n *FlexibleNumber) UnmarshalJSON(data []byte) error {
	value := strings.Trim(string(data), "\"")
	*n = FlexibleNumber(value)
	return nil
}

func (n FlexibleNumber) String() string {
	return string(n)
}

type Verse struct {
	Verse FlexibleNumber `json:"verse"`
	Text  string         `json:"text"`
}

type Chapter struct {
	Chapter FlexibleNumber `json:"chapter"`
	Verses  []Verse        `json:"verses"`
}

type Book struct {
	Book     string    `json:"book"`
	Chapters []Chapter `json:"chapters"`
}

type Source struct {
	Edition      string
	Base         string
	Output       string
	RemoveSpaces bool
}

func download(rawURL string) ([]byte, error) {
	resp, err := http.Get(rawURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %s for %s", resp.Status, rawURL)
	}

	return io.ReadAll(resp.Body)
}

func clean(s string) string {
	s = strings.ReplaceAll(s, "\t", " ")
	s = strings.ReplaceAll(s, "\r", " ")
	s = strings.ReplaceAll(s, "\n", " ")
	return strings.TrimSpace(s)
}

func fetchBookNames(base string) ([]string, error) {
	raw, err := download(base + "/Books.json")
	if err != nil {
		return nil, err
	}

	var books []string
	if err := json.Unmarshal(raw, &books); err != nil {
		return nil, err
	}

	return books, nil
}

func fetchBook(source Source, name string) (Book, error) {
	filenameName := name

	if source.RemoveSpaces {
		filenameName = strings.ReplaceAll(filenameName, " ", "")
	}

	filename := url.PathEscape(filenameName + ".json")
	rawURL := source.Base + "/" + filename

	raw, err := download(rawURL)
	if err != nil {
		return Book{}, err
	}

	var book Book
	if err := json.Unmarshal(raw, &book); err != nil {
		return Book{}, err
	}

	return book, nil
}

func buildEdition(source Source) error {
	fmt.Printf("Fetching %s book list...\n", source.Edition)

	books, err := fetchBookNames(source.Base)
	if err != nil {
		return err
	}

	fmt.Printf("  books: %d\n", len(books))

	if err := os.MkdirAll(filepath.Dir(source.Output), 0755); err != nil {
		return err
	}

	f, err := os.Create(source.Output)
	if err != nil {
		return err
	}
	defer f.Close()

	totalVerses := 0

	for i, name := range books {
		fmt.Printf("  [%02d/%02d] %s\n", i+1, len(books), name)

		book, err := fetchBook(source, name)
		if err != nil {
			return fmt.Errorf("%s: %w", name, err)
		}

		for _, chapter := range book.Chapters {
			for _, verse := range chapter.Verses {
				fields := []string{
					source.Edition,
					clean(book.Book),
					chapter.Chapter.String(),
					verse.Verse.String(),
					clean(verse.Text),
				}

				if _, err := fmt.Fprintln(f, strings.Join(fields, "\t")); err != nil {
					return err
				}

				totalVerses++
			}
		}
	}

	fmt.Printf("  verses: %d\n", totalVerses)
	fmt.Printf("  wrote: %s\n", source.Output)

	return nil
}

func main() {
	sources := []Source{
		{
			Edition:      "kjv",
			Base:         "https://raw.githubusercontent.com/aruljohn/Bible-kjv/master",
			Output:       "data/kjv.tsv",
			RemoveSpaces: true,
		},
		{
			Edition:      "1611",
			Base:         "https://raw.githubusercontent.com/aruljohn/Bible-kjv-1611/main",
			Output:       "data/kjv1611.tsv",
			RemoveSpaces: false,
		},
	}

	for _, source := range sources {
		if err := buildEdition(source); err != nil {
			fmt.Fprintf(os.Stderr, "\n%s: %v\n", source.Edition, err)
			os.Exit(1)
		}
	}

	fmt.Println("\nBible data build complete.")
}
