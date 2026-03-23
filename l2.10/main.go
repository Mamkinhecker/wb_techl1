package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
)

type config struct {
	column   int    // -k N (1-индекс)
	numeric  bool   // -n
	reverse  bool   // -r
	unique   bool   // -u
	month    bool   // -M
	trim     bool   // -b (игнорировать хвостовые пробелы)
	check    bool   // -c (проверить сортировку)
	human    bool   // -h
	filename string // имя файла (если не указан, читаем stdin)
}

func main() {
	cfg := ParseFlags()

	lines, err := ReadLines(cfg.filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading input: %v\n", err)
		os.Exit(1)
	}

	if cfg.check {
		if !isSorted(lines, cfg) {
			fmt.Fprintln(os.Stderr, "input is not sorted")
			os.Exit(1)
		}
		return
	}

	sortLines(lines, cfg)

	if cfg.unique {
		lines = unique(lines, cfg)
	}

	writeLines(lines)
}

func ParseFlags() config {
	var (
		col     = flag.Int("k", 0, "сортировать по колонке N")
		numeric = flag.Bool("n", false, "числовая сортировка")
		reverse = flag.Bool("r", false, "сортировка в обратном порядке")
		unique  = flag.Bool("u", false, "выводить только уникальные строки")
		month   = flag.Bool("M", false, "сортировка по названию месяца")
		trim    = flag.Bool("b", false, "игнорировать хвостовые пробелы")
		check   = flag.Bool("c", false, "проверить сортировку")
		human   = flag.Bool("h", false, "сортировка по человекочитаемым размерам")
	)
	flag.Parse()

	cfg := config{
		column:   *col,
		numeric:  *numeric,
		reverse:  *reverse,
		unique:   *unique,
		month:    *month,
		trim:     *trim,
		check:    *check,
		human:    *human,
		filename: getFilename(),
	}
	return cfg
}

func getFilename() string {
	args := flag.Args()
	if len(args) > 0 {
		return args[0]
	}
	return ""
}

func ReadLines(fileName string) ([]string, error) {
	var r io.Reader
	if fileName == "" {
		r = os.Stdin
	} else {
		f, err := os.Open(fileName)
		if err != nil {
			return nil, err
		}
		defer f.Close()

		r = f
	}

	var lines []string
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return lines, nil
}

func writeLines(lines []string) {
	for _, s := range lines {
		fmt.Println(s)
	}
}

func isSorted(lines []string, cfg config) bool {
	for i := 1; i < len(lines); i++ {
		if less(lines[i], lines[i-1], cfg) {
			return false
		}
	}
	return true
}

func sortLines(lines []string, cfg config) {
	sort.Slice(lines, func(i, j int) bool {
		return less(lines[i], lines[j], cfg)
	})
}

func less(a, b string, cfg config) bool {
	valA, valB := getCompareValues(a, b, cfg)
	cmp := compare(valA, valB, cfg)
	if cfg.reverse {
		return cmp > 0
	}
	return cmp < 0
}

func getCompareValues(a, b string, cfg config) (string, string) {
	if cfg.column == 0 {
		return a, b
	}
	col := cfg.column - 1
	colsA := strings.Split(a, "\t")
	colsB := strings.Split(b, "\t")
	valA := ""
	if col < len(colsA) {
		valA = colsA[col]
	}
	valB := ""
	if col < len(colsB) {
		valB = colsB[col]
	}
	if cfg.trim {
		valA = strings.TrimRight(valA, " \t")
		valB = strings.TrimRight(valB, " \t")
	}
	return valA, valB
}

func compare(a, b string, cfg config) int {
	if cfg.numeric {
		return compareNumeric(a, b)
	}
	if cfg.human {
		return compareHuman(a, b)
	}
	if cfg.month {
		return compareMonth(a, b)
	}
	return strings.Compare(a, b)
}

func compareNumeric(a, b string) int {
	fa, errA := strconv.ParseFloat(a, 64)
	fb, errB := strconv.ParseFloat(b, 64)
	if errA != nil && errB != nil {
		return strings.Compare(a, b)
	}
	if errA != nil {
		return -1
	}
	if errB != nil {
		return 1
	}

	switch {
	case fa < fb:
		return -1
	case fa > fb:
		return 1
	default:
		return 0
	}
}

func parseHumanSize(s string) (float64, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, fmt.Errorf("empty")
	}

	suffix := s[len(s)-1]
	var numStr string
	var factor float64
	switch suffix {
	case 'K', 'k':
		factor = 1024
		numStr = s[:len(s)-1]
	case 'M', 'm':
		factor = 1024 * 1024
		numStr = s[:len(s)-1]
	case 'G', 'g':
		factor = 1024 * 1024 * 1024
		numStr = s[:len(s)-1]
	case 'T', 't':
		factor = 1024 * 1024 * 1024 * 1024
		numStr = s[:len(s)-1]
	case 'P', 'p':
		factor = 1024 * 1024 * 1024 * 1024 * 1024
		numStr = s[:len(s)-1]
	case 'E', 'e':
		factor = 1024 * 1024 * 1024 * 1024 * 1024 * 1024
		numStr = s[:len(s)-1]
	default:
		factor = 1
		numStr = s
	}
	num, err := strconv.ParseFloat(numStr, 64)
	if err != nil {
		return 0, err
	}
	return num * factor, nil
}

func compareHuman(a, b string) int {
	va, errA := parseHumanSize(a)
	vb, errB := parseHumanSize(b)
	if errA != nil && errB != nil {
		return strings.Compare(a, b)
	}
	if errA != nil {
		return -1
	}
	if errB != nil {
		return 1
	}
	switch {
	case va < vb:
		return -1
	case va > vb:
		return 1
	default:
		return 0
	}
}

func monthValue(s string) int {
	months := map[string]int{
		"Jan": 1, "Feb": 2, "Mar": 3, "Apr": 4,
		"May": 5, "Jun": 6, "Jul": 7, "Aug": 8,
		"Sep": 9, "Oct": 10, "Nov": 11, "Dec": 12,
	}

	trimmed := strings.TrimSpace(s)
	if len(trimmed) >= 3 {
		key := trimmed[:3]
		if val, ok := months[key]; ok {
			return val
		}
	}
	return 0
}

func compareMonth(a, b string) int {
	ma := monthValue(a)
	mb := monthValue(b)
	if ma == 0 && mb == 0 {
		return strings.Compare(a, b)
	}
	if ma == 0 {
		return -1
	}
	if mb == 0 {
		return 1
	}
	switch {
	case ma < mb:
		return -1
	case ma > mb:
		return 1
	default:
		return 0
	}
}

func unique(lines []string, cfg config) []string {
	if len(lines) == 0 {
		return lines
	}
	res := []string{lines[0]}
	for i := 1; i < len(lines); i++ {
		if !equal(lines[i], lines[i-1], cfg) {
			res = append(res, lines[i])
		}
	}
	return res
}

func equal(a, b string, cfg config) bool {
	return !less(a, b, cfg) && !less(b, a, cfg)
}
