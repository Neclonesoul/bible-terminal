package main

import (
	"bufio"
	"fmt"
	"html"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"

	bibledata "github.com/Neclonesoul/bible-terminal/internal/data"
	biblesearch "github.com/Neclonesoul/bible-terminal/internal/search"
	"github.com/Neclonesoul/bible-terminal/internal/speech"
	bibletui "github.com/Neclonesoul/bible-terminal/internal/tui"
)

type Verse struct {
	Edition string
	Book    string
	Chapter int
	Verse   int
	Text    string
}

func usage() {
	fmt.Println(`BIBLE TERMINAL

Usage:
  bible BOOK CHAPTER
  bible BOOK CHAPTER:VERSE
  bible --kjv BOOK CHAPTER:VERSE
  bible --1611 BOOK CHAPTER:VERSE
  bible --compare BOOK CHAPTER:VERSE
  bible --sound BOOK CHAPTER:VERSE
  bible --search "PHRASE"

Add --sound (or --voice, --speak) to any reference to hear it read aloud.

Examples:
  bible Psalm 23
  bible Jn 3:16
  bible --1611 John 3:16
  bible --compare John 3:16
  bible --sound Psalm 23`)
}

func canonicalBook(book string) string {
	key := strings.ToLower(strings.Join(strings.Fields(book), " "))

	aliases := map[string]string{
		"gen": "Genesis",
		"ge":  "Genesis",
		"ex":  "Exodus",
		"exo": "Exodus",
		"lev": "Leviticus",
		"num": "Numbers",

		"deut": "Deuteronomy",
		"dt":   "Deuteronomy",
		"josh": "Joshua",
		"judg": "Judges",

		"psalm":  "Psalms",
		"psalms": "Psalms",
		"ps":     "Psalms",
		"psa":    "Psalms",

		"prov": "Proverbs",
		"pr":   "Proverbs",

		"eccl": "Ecclesiastes",
		"ecc":  "Ecclesiastes",

		"song":          "Song of Solomon",
		"sos":           "Song of Solomon",
		"song of songs": "Song of Solomon",

		"isa":  "Isaiah",
		"jer":  "Jeremiah",
		"lam":  "Lamentations",
		"ezek": "Ezekiel",
		"dan":  "Daniel",

		"matt": "Matthew",
		"mt":   "Matthew",
		"mk":   "Mark",
		"mrk":  "Mark",
		"lk":   "Luke",

		"jn":  "John",
		"jhn": "John",

		"rom": "Romans",

		"1cor":  "1 Corinthians",
		"1 cor": "1 Corinthians",
		"2cor":  "2 Corinthians",
		"2 cor": "2 Corinthians",

		"gal":  "Galatians",
		"eph":  "Ephesians",
		"phil": "Philippians",
		"col":  "Colossians",

		"1thess":  "1 Thessalonians",
		"1 thess": "1 Thessalonians",
		"2thess":  "2 Thessalonians",
		"2 thess": "2 Thessalonians",

		"1tim":  "1 Timothy",
		"1 tim": "1 Timothy",
		"2tim":  "2 Timothy",
		"2 tim": "2 Timothy",

		"heb": "Hebrews",
		"jas": "James",

		"1pet":  "1 Peter",
		"1 pet": "1 Peter",
		"2pet":  "2 Peter",
		"2 pet": "2 Peter",

		"1jn":  "1 John",
		"1 jn": "1 John",
		"2jn":  "2 John",
		"2 jn": "2 John",
		"3jn":  "3 John",
		"3 jn": "3 John",

		"rev": "Revelation",
	}

	if v, ok := aliases[key]; ok {
		return v
	}

	return strings.TrimSpace(book)
}

func parseReference(args []string) (string, int, int, bool, error) {
	if len(args) == 0 {
		return "", 0, 0, false, fmt.Errorf("missing reference")
	}

	last := args[len(args)-1]

	// A book name on its own ("john", "1 john", "song of solomon")
	// opens chapter 1.
	if len(args) == 1 || (last[0] < '0' || last[0] > '9') {
		return canonicalBook(strings.Join(args, " ")), 1, 0, false, nil
	}

	book := canonicalBook(strings.Join(args[:len(args)-1], " "))

	if strings.Contains(last, ":") {
		parts := strings.SplitN(last, ":", 2)

		chapter, err := strconv.Atoi(parts[0])
		if err != nil {
			return "", 0, 0, false, fmt.Errorf("invalid chapter")
		}

		verse, err := strconv.Atoi(parts[1])
		if err != nil {
			return "", 0, 0, false, fmt.Errorf("invalid verse")
		}

		return book, chapter, verse, true, nil
	}

	chapter, err := strconv.Atoi(last)
	if err != nil {
		return "", 0, 0, false, fmt.Errorf("invalid chapter")
	}

	return book, chapter, 0, false, nil
}

func load(rawData, edition, book string, chapter, verse int, specific bool) ([]Verse, error) {
	var results []Verse

	scanner := bufio.NewScanner(strings.NewReader(rawData))
	scanner.Buffer(make([]byte, 64*1024), 1024*1024)

	for scanner.Scan() {
		fields := strings.SplitN(scanner.Text(), "\t", 5)
		if len(fields) != 5 {
			continue
		}

		if !strings.EqualFold(fields[1], book) {
			continue
		}

		ch, err := strconv.Atoi(fields[2])
		if err != nil || ch != chapter {
			continue
		}

		v, err := strconv.Atoi(fields[3])
		if err != nil {
			continue
		}

		if specific && v != verse {
			continue
		}

		results = append(results, Verse{
			Edition: edition,
			Book:    fields[1],
			Chapter: ch,
			Verse:   v,
			Text:    html.UnescapeString(fields[4]),
		})
	}

	return results, scanner.Err()
}

func terminalWidth() int {
	if env := os.Getenv("COLUMNS"); env != "" {
		if n, err := strconv.Atoi(env); err == nil && n >= 40 {
			if n > 180 {
				return 180
			}
			return n
		}
	}

	out, err := exec.Command("tput", "cols").Output()
	if err == nil {
		if n, err := strconv.Atoi(strings.TrimSpace(string(out))); err == nil && n >= 40 {
			if n > 180 {
				return 180
			}
			return n
		}
	}

	return 80
}

func runeLen(s string) int {
	return utf8.RuneCountInString(s)
}

func padRight(s string, width int) string {
	n := width - runeLen(s)
	if n <= 0 {
		return s
	}
	return s + strings.Repeat(" ", n)
}

func wrap(text string, width int) []string {
	if width < 10 {
		width = 10
	}

	words := strings.Fields(text)
	if len(words) == 0 {
		return []string{""}
	}

	var lines []string
	current := words[0]

	for _, word := range words[1:] {
		if runeLen(current)+1+runeLen(word) <= width {
			current += " " + word
		} else {
			lines = append(lines, current)
			current = word
		}
	}

	lines = append(lines, current)
	return lines
}

func verseLines(v Verse, width int) []string {
	prefix := fmt.Sprintf("%d  ", v.Verse)
	bodyWidth := width - runeLen(prefix)

	if bodyWidth < 10 {
		bodyWidth = 10
	}

	body := wrap(v.Text, bodyWidth)
	lines := make([]string, len(body))

	for i, line := range body {
		if i == 0 {
			lines[i] = prefix + line
		} else {
			lines[i] = strings.Repeat(" ", runeLen(prefix)) + line
		}
	}

	return lines
}

func printEdition(title string, verses []Verse) {
	width := terminalWidth()

	fmt.Printf("\n%s\n", title)
	fmt.Println(strings.Repeat("─", width))

	for _, v := range verses {
		for _, line := range verseLines(v, width) {
			fmt.Println(line)
		}
		fmt.Println()
	}
}

func verseMap(verses []Verse) map[int]Verse {
	m := make(map[int]Verse)
	for _, v := range verses {
		m[v.Verse] = v
	}
	return m
}

func allVerseNumbers(a, b []Verse) []int {
	seen := make(map[int]bool)

	for _, v := range a {
		seen[v.Verse] = true
	}
	for _, v := range b {
		seen[v.Verse] = true
	}

	var nums []int
	for n := range seen {
		nums = append(nums, n)
	}

	sort.Ints(nums)
	return nums
}

func printStackedComparison(book string, chapter int, kjv, old []Verse) {
	width := terminalWidth()
	k := verseMap(kjv)
	o := verseMap(old)

	fmt.Printf("\nTHE HOLY BIBLE · COMPARISON · %s %d\n",
		strings.ToUpper(book), chapter)
	fmt.Println(strings.Repeat("═", width))

	for _, n := range allVerseNumbers(kjv, old) {
		if v, ok := k[n]; ok {
			fmt.Printf("\nKJV · %d\n", n)
			for _, line := range wrap(v.Text, width) {
				fmt.Println(line)
			}
		}

		fmt.Println(strings.Repeat("─", width))

		if v, ok := o[n]; ok {
			fmt.Printf("1611 · %d\n", n)
			for _, line := range wrap(v.Text, width) {
				fmt.Println(line)
			}
		}

		fmt.Println(strings.Repeat("═", width))
	}
}

func printSideBySide(book string, chapter int, kjv, old []Verse) {
	width := terminalWidth()

	separator := " │ "
	colWidth := (width - runeLen(separator)) / 2

	k := verseMap(kjv)
	o := verseMap(old)

	fmt.Printf("\nTHE HOLY BIBLE · COMPARISON · %s %d\n",
		strings.ToUpper(book), chapter)

	fmt.Println(strings.Repeat("═", width))

	fmt.Printf("%s%s%s\n",
		padRight("KJV", colWidth),
		separator,
		"KJV 1611",
	)

	fmt.Printf("%s%s%s\n",
		strings.Repeat("─", colWidth),
		"─┼─",
		strings.Repeat("─", colWidth),
	)

	for _, n := range allVerseNumbers(kjv, old) {
		var left []string
		var right []string

		if v, ok := k[n]; ok {
			left = verseLines(v, colWidth)
		}
		if v, ok := o[n]; ok {
			right = verseLines(v, colWidth)
		}

		height := len(left)
		if len(right) > height {
			height = len(right)
		}

		for i := 0; i < height; i++ {
			l := ""
			r := ""

			if i < len(left) {
				l = left[i]
			}
			if i < len(right) {
				r = right[i]
			}

			fmt.Printf("%s%s%s\n",
				padRight(l, colWidth),
				separator,
				r,
			)
		}

		fmt.Printf("%s%s%s\n",
			strings.Repeat("·", colWidth),
			"─┼─",
			strings.Repeat("·", colWidth),
		)
	}
}

func printComparison(book string, chapter int, kjv, old []Verse) {
	if terminalWidth() >= 100 {
		printSideBySide(book, chapter, kjv, old)
		return
	}

	printStackedComparison(book, chapter, kjv, old)
}

func printSearchResults(edition, query string, results []biblesearch.Result) {
	width := terminalWidth()

	fmt.Printf("\nSEARCH · %s · %q\n", edition, query)
	fmt.Println(strings.Repeat("═", width))

	if len(results) == 0 {
		fmt.Println("No verses found.")
		return
	}

	for _, r := range results {
		ref := fmt.Sprintf("%s %d:%d", r.Book, r.Chapter, r.Verse)

		fmt.Printf("\n%s\n", ref)

		for _, line := range wrap(r.Text, width) {
			fmt.Println(line)
		}
	}

	fmt.Printf("\n%s\n", strings.Repeat("─", width))
	fmt.Printf("%d result", len(results))

	if len(results) != 1 {
		fmt.Print("s")
	}

	fmt.Println()
}

// speakPassage reads the passage aloud. The standard KJV wording is
// always used, because speech engines cannot pronounce 1611 spelling.
func speakPassage(book string, chapter, verse int, specific bool) {
	verses, err := load(bibledata.KJV, "KJV", book, chapter, verse, specific)
	if err != nil || len(verses) == 0 {
		fmt.Fprintf(os.Stderr,
			"Cannot read aloud: %s %d has no modern-spelling KJV text.\n",
			book, chapter)
		os.Exit(1)
	}

	announce := ""
	if !specific {
		announce = fmt.Sprintf("%s, chapter %d", verses[0].Book, chapter)
	}

	texts := make([]string, len(verses))
	for i, v := range verses {
		texts[i] = v.Text
	}

	if err := speech.Speak(speech.Passage(announce, texts)); err != nil {
		fmt.Fprintln(os.Stderr, "Cannot read aloud:", err)
		os.Exit(1)
	}
}

func main() {

	// SEARCH_MODE_HANDLER_V1
	// Handle full-text search before normal reference parsing.
	rawArgs := os.Args[1:]

	searchMode := false
	searchEdition := "kjv"
	searchQuery := ""

	for i := 0; i < len(rawArgs); i++ {
		switch rawArgs[i] {
		case "--1611":
			searchEdition = "1611"

		case "--kjv":
			searchEdition = "kjv"

		case "--search":
			searchMode = true

			if i+1 < len(rawArgs) {
				searchQuery = rawArgs[i+1]
				i++
			}
		}
	}

	if searchMode {
		if strings.TrimSpace(searchQuery) == "" {
			fmt.Fprintln(os.Stderr, "Error: --search requires a search phrase")
			os.Exit(1)
		}

		if searchEdition == "1611" {
			results := biblesearch.Find(bibledata.KJV1611, searchQuery)
			printSearchResults("KJV 1611", searchQuery, results)
			return
		}

		results := biblesearch.Find(bibledata.KJV, searchQuery)
		printSearchResults("KJV", searchQuery, results)
		return
	}

	speak := false

	var args []string
	for _, arg := range os.Args[1:] {
		if arg == "--sound" || arg == "--voice" || arg == "--speak" {
			speak = true
			continue
		}
		args = append(args, arg)
	}

	if len(args) == 0 {
		if speak {
			fmt.Fprintln(os.Stderr, "Error: --sound requires a reference")
			usage()
			os.Exit(1)
		}

		bibletui.Run()
		return
	}

	mode := "kjv"

	switch args[0] {
	case "--kjv":
		args = args[1:]

	case "--1611":
		mode = "1611"
		args = args[1:]

	case "--compare":
		mode = "compare"
		args = args[1:]

	case "-h", "--help":
		usage()
		return
	}

	book, chapter, verse, specific, err := parseReference(args)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		usage()
		os.Exit(1)
	}

	switch mode {
	case "kjv":
		verses, err := load(
			bibledata.KJV,
			"KJV",
			book,
			chapter,
			verse,
			specific,
		)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		if len(verses) == 0 {
			fmt.Fprintf(os.Stderr, "Reference not found: %s %d\n", book, chapter)
			os.Exit(1)
		}

		printEdition(
			fmt.Sprintf("THE HOLY BIBLE · KJV · %s %d",
				strings.ToUpper(book), chapter),
			verses,
		)

	case "1611":
		verses, err := load(
			bibledata.KJV1611,
			"KJV 1611",
			book,
			chapter,
			verse,
			specific,
		)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		if len(verses) == 0 {
			fmt.Fprintf(os.Stderr, "Reference not found: %s %d\n", book, chapter)
			os.Exit(1)
		}

		printEdition(
			fmt.Sprintf("THE HOLY BIBLE · 1611 · %s %d",
				strings.ToUpper(book), chapter),
			verses,
		)

	case "compare":
		kjv, err := load(
			bibledata.KJV,
			"KJV",
			book,
			chapter,
			verse,
			specific,
		)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		old, err := load(
			bibledata.KJV1611,
			"KJV 1611",
			book,
			chapter,
			verse,
			specific,
		)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		if len(kjv) == 0 && len(old) == 0 {
			fmt.Fprintf(os.Stderr, "Reference not found: %s %d\n", book, chapter)
			os.Exit(1)
		}

		printComparison(book, chapter, kjv, old)
	}

	if speak {
		speakPassage(book, chapter, verse, specific)
	}
}
