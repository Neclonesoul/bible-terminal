package tui

import (
	"bufio"
	"fmt"
	"html"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"

	bibledata "github.com/Neclonesoul/bible-terminal/internal/data"
	biblesearch "github.com/Neclonesoul/bible-terminal/internal/search"
	"golang.org/x/term"
)

type Verse struct {
	Book    string
	Chapter int
	Verse   int
	Text    string
}

type Reader struct {
	Book    string
	Chapter int
	Edition string
	Mode    string
	Scroll  int
}

func New() *Reader {
	return &Reader{
		Book:    "Genesis",
		Chapter: 1,
		Edition: "kjv",
		Mode:    "read",
	}
}

func clear() {
	fmt.Print("\033[2J\033[H")
}

func dimensions() (int, int) {
	w, h, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil || w < 35 || h < 12 {
		return 80, 24
	}

	if w > 160 {
		w = 160
	}

	return w, h
}

func runeLen(s string) int {
	return utf8.RuneCountInString(s)
}

func trimWidth(s string, width int) string {
	if width <= 0 {
		return ""
	}

	r := []rune(s)

	if len(r) <= width {
		return s
	}

	if width == 1 {
		return "…"
	}

	return string(r[:width-1]) + "…"
}

func pad(s string, width int) string {
	s = trimWidth(s, width)

	if n := width - runeLen(s); n > 0 {
		s += strings.Repeat(" ", n)
	}

	return s
}

func center(s string, width int) string {
	s = trimWidth(s, width)

	n := width - runeLen(s)
	if n <= 0 {
		return s
	}

	left := n / 2
	right := n - left

	return strings.Repeat(" ", left) +
		s +
		strings.Repeat(" ", right)
}

func wrap(text string, width int) []string {
	if width < 8 {
		width = 8
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

func editionData(edition string) string {
	if edition == "1611" {
		return bibledata.KJV1611
	}

	return bibledata.KJV
}

func editionTitle(edition string) string {
	if edition == "1611" {
		return "KING JAMES · 1611"
	}

	return "KING JAMES VERSION"
}

func parseBooks(data string) []string {
	var books []string
	seen := make(map[string]bool)

	scanner := bufio.NewScanner(strings.NewReader(data))
	scanner.Buffer(make([]byte, 64*1024), 1024*1024)

	for scanner.Scan() {
		fields := strings.SplitN(scanner.Text(), "\t", 5)

		if len(fields) != 5 {
			continue
		}

		book := fields[1]

		if !seen[book] {
			seen[book] = true
			books = append(books, book)
		}
	}

	return books
}

func load(data, book string, chapter int) []Verse {
	var verses []Verse

	scanner := bufio.NewScanner(strings.NewReader(data))
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

		verses = append(verses, Verse{
			Book:    fields[1],
			Chapter: ch,
			Verse:   v,
			Text:    html.UnescapeString(fields[4]),
		})
	}

	return verses
}

func chapterExists(data, book string, chapter int) bool {
	return len(load(data, book, chapter)) > 0
}

func bookExists(data, book string) bool {
	books := parseBooks(data)

	for _, b := range books {
		if strings.EqualFold(b, book) {
			return true
		}
	}

	return false
}

func roman(n int) string {
	values := []struct {
		n int
		s string
	}{
		{1000, "M"}, {900, "CM"},
		{500, "D"}, {400, "CD"},
		{100, "C"}, {90, "XC"},
		{50, "L"}, {40, "XL"},
		{10, "X"}, {9, "IX"},
		{5, "V"}, {4, "IV"},
		{1, "I"},
	}

	var b strings.Builder

	for _, item := range values {
		for n >= item.n {
			b.WriteString(item.s)
			n -= item.n
		}
	}

	return b.String()
}

func canonicalBook(input string, data string) string {
	key := strings.ToLower(
		strings.Join(strings.Fields(input), " "),
	)

	aliases := map[string]string{
		"gen": "Genesis",
		"ge":  "Genesis",
		"exo": "Exodus",
		"ex":  "Exodus",
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

	if canonical, ok := aliases[key]; ok {
		if bookExists(data, canonical) {
			return canonical
		}
	}

	for _, book := range parseBooks(data) {
		if strings.EqualFold(book, input) {
			return book
		}
	}

	return strings.TrimSpace(input)
}

func verseLines(verses []Verse, width int) []string {
	var lines []string

	textWidth := width - 8

	if textWidth < 20 {
		textWidth = width - 6
	}

	for _, verse := range verses {
		prefix := fmt.Sprintf("%d  ", verse.Verse)
		indent := strings.Repeat(" ", runeLen(prefix))

		bodyWidth := textWidth - runeLen(prefix)

		for i, line := range wrap(verse.Text, bodyWidth) {
			if i == 0 {
				lines = append(
					lines,
					"  "+prefix+line,
				)
			} else {
				lines = append(
					lines,
					"  "+indent+line,
				)
			}
		}

		lines = append(lines, "")
	}

	return lines
}

func comparisonLines(kjv, old []Verse, width int) []string {
	// Narrow terminals use stacked comparison.
	if width < 100 {
		var lines []string

		kmap := make(map[int]Verse)
		omap := make(map[int]Verse)

		for _, v := range kjv {
			kmap[v.Verse] = v
		}

		for _, v := range old {
			omap[v.Verse] = v
		}

		max := 0

		for n := range kmap {
			if n > max {
				max = n
			}
		}

		for n := range omap {
			if n > max {
				max = n
			}
		}

		for n := 1; n <= max; n++ {
			if v, ok := kmap[n]; ok {
				lines = append(
					lines,
					fmt.Sprintf("  KJV · %d", n),
				)

				for _, line := range wrap(v.Text, width-8) {
					lines = append(lines, "  "+line)
				}
			}

			lines = append(
				lines,
				"  "+strings.Repeat("─", width-8),
			)

			if v, ok := omap[n]; ok {
				lines = append(
					lines,
					fmt.Sprintf("  1611 · %d", n),
				)

				for _, line := range wrap(v.Text, width-8) {
					lines = append(lines, "  "+line)
				}
			}

			lines = append(lines, "")
		}

		return lines
	}

	// Wide terminals get real columns.
	inner := width - 6
	separator := " │ "
	col := (inner - runeLen(separator)) / 2

	kmap := make(map[int]Verse)
	omap := make(map[int]Verse)

	max := 0

	for _, v := range kjv {
		kmap[v.Verse] = v
		if v.Verse > max {
			max = v.Verse
		}
	}

	for _, v := range old {
		omap[v.Verse] = v
		if v.Verse > max {
			max = v.Verse
		}
	}

	var lines []string

	lines = append(
		lines,
		"  "+pad("KJV", col)+separator+"KJV 1611",
	)

	lines = append(
		lines,
		"  "+
			strings.Repeat("─", col)+
			"─┼─"+
			strings.Repeat("─", col),
	)

	for n := 1; n <= max; n++ {
		var left []string
		var right []string

		if v, ok := kmap[n]; ok {
			left = wrap(
				fmt.Sprintf("%d  %s", n, v.Text),
				col,
			)
		}

		if v, ok := omap[n]; ok {
			right = wrap(
				fmt.Sprintf("%d  %s", n, v.Text),
				col,
			)
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

			lines = append(
				lines,
				"  "+pad(l, col)+separator+r,
			)
		}

		lines = append(lines, "")
	}

	return lines
}

func (r *Reader) contentLines(width int) []string {
	if r.Mode == "compare" {
		return comparisonLines(
			load(bibledata.KJV, r.Book, r.Chapter),
			load(bibledata.KJV1611, r.Book, r.Chapter),
			width,
		)
	}

	return verseLines(
		load(
			editionData(r.Edition),
			r.Book,
			r.Chapter,
		),
		width,
	)
}

func (r *Reader) nextChapter() {
	data := editionData(r.Edition)

	if r.Mode == "compare" {
		data = bibledata.KJV
	}

	if chapterExists(data, r.Book, r.Chapter+1) {
		r.Chapter++
		r.Scroll = 0
		return
	}

	books := parseBooks(data)

	for i, book := range books {
		if strings.EqualFold(book, r.Book) &&
			i+1 < len(books) {

			r.Book = books[i+1]
			r.Chapter = 1
			r.Scroll = 0
			return
		}
	}
}

func (r *Reader) previousChapter() {
	if r.Chapter > 1 {
		r.Chapter--
		r.Scroll = 0
		return
	}

	data := editionData(r.Edition)

	if r.Mode == "compare" {
		data = bibledata.KJV
	}

	books := parseBooks(data)

	for i, book := range books {
		if !strings.EqualFold(book, r.Book) || i == 0 {
			continue
		}

		previous := books[i-1]
		chapter := 1

		for chapterExists(data, previous, chapter+1) {
			chapter++
		}

		r.Book = previous
		r.Chapter = chapter
		r.Scroll = 0
		return
	}
}

func (r *Reader) toggleEdition() {
	if r.Edition == "kjv" {
		r.Edition = "1611"
		r.Mode = "read"
		r.Scroll = 0
		return
	}

	// A 1611-only book may not exist in the 66-book KJV.
	if !bookExists(bibledata.KJV, r.Book) {
		r.Book = "Genesis"
		r.Chapter = 1
	}

	r.Edition = "kjv"
	r.Mode = "read"
	r.Scroll = 0
}

func (r *Reader) toggleComparison() {
	if r.Mode == "compare" {
		r.Mode = "read"
		r.Scroll = 0
		return
	}

	// Comparison only makes sense for books shared by both editions.
	if !bookExists(bibledata.KJV, r.Book) ||
		!bookExists(bibledata.KJV1611, r.Book) {
		return
	}

	r.Mode = "compare"
	r.Scroll = 0
}

func (r *Reader) gotoReference(input string) {
	parts := strings.Fields(strings.TrimSpace(input))

	if len(parts) < 2 {
		return
	}

	last := parts[len(parts)-1]

	// Accept John 3:16 but TUI currently opens the chapter.
	if strings.Contains(last, ":") {
		last = strings.SplitN(last, ":", 2)[0]
	}

	chapter, err := strconv.Atoi(last)

	if err != nil || chapter < 1 {
		return
	}

	data := editionData(r.Edition)

	if r.Mode == "compare" {
		data = bibledata.KJV
	}

	book := canonicalBook(
		strings.Join(parts[:len(parts)-1], " "),
		data,
	)

	if !chapterExists(data, book, chapter) {
		return
	}

	r.Book = book
	r.Chapter = chapter
	r.Scroll = 0
}

func (r *Reader) render() {
	clear()

	width, height := dimensions()

	if width < 35 {
		width = 35
	}

	inner := width - 4

	headerHeight := 6
	footerHeight := 4

	viewportHeight := height - headerHeight - footerHeight

	if viewportHeight < 5 {
		viewportHeight = 5
	}

	lines := r.contentLines(width)

	maxScroll := len(lines) - viewportHeight

	if maxScroll < 0 {
		maxScroll = 0
	}

	if r.Scroll < 0 {
		r.Scroll = 0
	}

	if r.Scroll > maxScroll {
		r.Scroll = maxScroll
	}

	end := r.Scroll + viewportHeight

	if end > len(lines) {
		end = len(lines)
	}

	fmt.Println(
		"┌" + strings.Repeat("─", width-2) + "┐",
	)

	fmt.Printf(
		"│ %s │\n",
		center("THE HOLY BIBLE", inner),
	)

	subtitle := editionTitle(r.Edition)

	if r.Mode == "compare" {
		subtitle = "KJV · COMPARISON · 1611"
	}

	fmt.Printf(
		"│ %s │\n",
		center(subtitle, inner),
	)

	fmt.Println(
		"├" + strings.Repeat("─", width-2) + "┤",
	)

	title := fmt.Sprintf(
		"%s %s",
		strings.ToUpper(r.Book),
		roman(r.Chapter),
	)

	fmt.Printf(
		"│ %s │\n",
		center(title, inner),
	)

	fmt.Printf(
		"│ %s │\n",
		pad("", inner),
	)

	for i := r.Scroll; i < end; i++ {
		fmt.Printf(
			"│ %s │\n",
			pad(lines[i], inner),
		)
	}

	for i := end - r.Scroll; i < viewportHeight; i++ {
		fmt.Printf(
			"│ %s │\n",
			pad("", inner),
		)
	}

	fmt.Println(
		"├" + strings.Repeat("─", width-2) + "┤",
	)

	percent := 0

	if maxScroll > 0 {
		percent = r.Scroll * 100 / maxScroll
	}

	mode := strings.ToUpper(r.Edition)

	if r.Mode == "compare" {
		mode = "COMPARE"
	}

	status := fmt.Sprintf(
		" %s %d · %s · %d%%",
		strings.ToUpper(r.Book),
		r.Chapter,
		mode,
		percent,
	)

	fmt.Printf(
		"│ %s │\n",
		pad(status, inner),
	)

	controls := "↑↓ scroll · ←→ chapter · b books · v · c · g · / · q"

	fmt.Printf(
		"│ %s │\n",
		center(controls, inner),
	)

	fmt.Println(
		"└" + strings.Repeat("─", width-2) + "┘",
	)
}

func prompt(fd int, state *term.State, label string) string {
	term.Restore(fd, state)

	fmt.Print("\033[?25h")
	fmt.Print("\n" + label)

	reader := bufio.NewReader(os.Stdin)

	line, _ := reader.ReadString('\n')

	fmt.Print("\033[?25l")

	return strings.TrimSpace(line)
}

func chooseBook(fd int, state *term.State, r *Reader) {
	term.Restore(fd, state)

	clear()
	fmt.Print("\033[?25h")

	data := editionData(r.Edition)

	if r.Mode == "compare" {
		data = bibledata.KJV
	}

	books := parseBooks(data)

	fmt.Println("THE HOLY BIBLE")
	fmt.Println()
	fmt.Println("BOOKS")
	fmt.Println(strings.Repeat("─", 40))

	for i, book := range books {
		fmt.Printf("%2d  %s\n", i+1, book)
	}

	fmt.Println()
	fmt.Print("Book number or name: ")

	input := bufio.NewReader(os.Stdin)
	value, _ := input.ReadString('\n')

	value = strings.TrimSpace(value)

	if n, err := strconv.Atoi(value); err == nil {
		if n >= 1 && n <= len(books) {
			r.Book = books[n-1]
			r.Chapter = 1
			r.Scroll = 0
		}
	} else {
		book := canonicalBook(value, data)

		if bookExists(data, book) {
			r.Book = book
			r.Chapter = 1
			r.Scroll = 0
		}
	}

	fmt.Print("\033[?25l")
}

func searchUI(fd int, state *term.State, r *Reader) {
	term.Restore(fd, state)

	clear()
	fmt.Print("\033[?25h")

	fmt.Printf(
		"Search %s: ",
		editionTitle(r.Edition),
	)

	input := bufio.NewReader(os.Stdin)

	query, _ := input.ReadString('\n')
	query = strings.TrimSpace(query)

	if query == "" {
		fmt.Print("\033[?25l")
		return
	}

	data := editionData(r.Edition)
	results := biblesearch.Find(data, query)

	clear()

	fmt.Printf(
		"SEARCH · %s · %q\n",
		editionTitle(r.Edition),
		query,
	)

	fmt.Println(strings.Repeat("═", 60))

	if len(results) == 0 {
		fmt.Println("No results.")
		fmt.Println()
		fmt.Print("Press Enter...")

		input.ReadString('\n')
		fmt.Print("\033[?25l")
		return
	}

	limit := len(results)

	if limit > 25 {
		limit = 25
	}

	for i := 0; i < limit; i++ {
		result := results[i]

		fmt.Printf(
			"%2d  %s %d:%d\n",
			i+1,
			result.Book,
			result.Chapter,
			result.Verse,
		)

		preview := trimWidth(
			strings.Join(strings.Fields(result.Text), " "),
			70,
		)

		fmt.Printf("    %s\n\n", preview)
	}

	if len(results) > limit {
		fmt.Printf(
			"Showing first %d of %d results.\n\n",
			limit,
			len(results),
		)
	}

	fmt.Print("Open result number (Enter to cancel): ")

	value, _ := input.ReadString('\n')
	value = strings.TrimSpace(value)

	n, err := strconv.Atoi(value)

	if err == nil && n >= 1 && n <= limit {
		result := results[n-1]

		r.Book = result.Book
		r.Chapter = result.Chapter
		r.Scroll = 0
		r.Mode = "read"
	}

	fmt.Print("\033[?25l")
}

func Run() {
	fd := int(os.Stdin.Fd())

	if !term.IsTerminal(fd) {
		fmt.Fprintln(
			os.Stderr,
			"Interactive mode requires a terminal.",
		)
		return
	}

	originalState, err := term.MakeRaw(fd)

	if err != nil {
		fmt.Fprintln(
			os.Stderr,
			"Could not enter raw terminal mode:",
			err,
		)
		return
	}

	defer func() {
		term.Restore(fd, originalState)
		fmt.Print("\033[?25h")
		clear()
	}()

	fmt.Print("\033[?25l")

	reader := New()
	input := bufio.NewReader(os.Stdin)

	for {
		reader.render()

		key, err := input.ReadByte()

		if err != nil {
			return
		}

		switch key {
		case 'q':
			return

		case 'j':
			reader.Scroll++

		case 'k':
			reader.Scroll--

		case 'l':
			reader.nextChapter()

		case 'h':
			reader.previousChapter()

		case 'v':
			reader.toggleEdition()

		case 'c':
			reader.toggleComparison()

		case 'g':
			ref := prompt(
				fd,
				originalState,
				"Go to: ",
			)

			if _, err := term.MakeRaw(fd); err == nil {
				fmt.Print("\033[?25l")
			}

			reader.gotoReference(ref)

		case 'b':
			chooseBook(
				fd,
				originalState,
				reader,
			)

			if _, err := term.MakeRaw(fd); err == nil {
				fmt.Print("\033[?25l")
			}

		case '/':
			searchUI(
				fd,
				originalState,
				reader,
			)

			if _, err := term.MakeRaw(fd); err == nil {
				fmt.Print("\033[?25l")
			}

		case 27:
			second, err := input.ReadByte()

			if err != nil || second != '[' {
				continue
			}

			third, err := input.ReadByte()

			if err != nil {
				continue
			}

			switch third {
			case 'A':
				reader.Scroll--

			case 'B':
				reader.Scroll++

			case 'C':
				reader.nextChapter()

			case 'D':
				reader.previousChapter()
			}
		}
	}
}
