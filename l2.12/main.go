package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

type config struct {
	after    int  // -A N
	before   int  // -B N
	context  int  // -C N
	count    bool // -c
	ignore   bool // -i
	invert   bool // -v
	fixed    bool // -F
	lineNum  bool // -n
	pattern  string
	filename string
}

func parseFlags() config {
	var (
		A = flag.Int("A", 0, "печатать N строк после совпадения")
		B = flag.Int("B", 0, "печатать N строк до совпадения")
		C = flag.Int("C", 0, "печатать N строк вокруг совпадения")
		c = flag.Bool("c", false, "печатать только количество совпавших строк")
		i = flag.Bool("i", false, "игнорировать регистр")
		v = flag.Bool("v", false, "инвертировать выбор")
		F = flag.Bool("F", false, "трактовать шаблон как фиксированную строку")
		n = flag.Bool("n", false, "печатать номер строки")
	)
	flag.Parse()

	pattern, filename := getArgs()
	cfg := config{
		after:    *A,
		before:   *B,
		context:  *C,
		count:    *c,
		ignore:   *i,
		invert:   *v,
		fixed:    *F,
		lineNum:  *n,
		pattern:  pattern,
		filename: filename,
	}
	return cfg
}

func getArgs() (string, string) {
	args := flag.Args()
	if len(args) > 1 {
		return args[0], args[1]
	}
	if len(args) == 1 {
		return args[0], ""
	}
	return "", ""
}

func main() {
	cfg := parseFlags()
	if cfg.pattern == "" {
		fmt.Fprintln(os.Stderr, "grep: pattern is empty")
		os.Exit(2)
	}

	lines, err := readLines(cfg.filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "grep: %v\n", err)
		os.Exit(2)
	}

	searchLines := make([]string, len(lines))
	copy(searchLines, lines)

	pattern := cfg.pattern
	if cfg.ignore {
		pattern = strings.ToLower(pattern)
		for i, s := range searchLines {
			searchLines[i] = strings.ToLower(s)
		}
	}

	var matches []int
	if cfg.fixed {
		matches = fixedMatches(searchLines, pattern)
	} else {
		matches, err = regexMatches(searchLines, pattern)
		if err != nil {
			fmt.Fprintf(os.Stderr, "grep: regex error: %v\n", err)
			os.Exit(2)
		}
	}

	if cfg.invert {
		// Инвертируем: все индексы, не входящие в matches
		m := make([]bool, len(lines))
		for _, idx := range matches {
			m[idx] = true
		}
		matches = nil
		for i, ok := range m {
			if !ok {
				matches = append(matches, i)
			}
		}
	}

	if cfg.count {
		fmt.Println(len(matches))
		return
	}

	intervals := buildIntervals(matches, lines, cfg)
	printResult(lines, intervals, cfg)
}

func readLines(filename string) ([]string, error) {
	var r io.Reader
	if filename == "" {
		r = os.Stdin
	} else {
		f, err := os.Open(filename)
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

func fixedMatches(lines []string, pattern string) []int {
	var res []int
	for i, line := range lines {
		if strings.Contains(line, pattern) {
			res = append(res, i)
		}
	}
	return res
}

func regexMatches(lines []string, pattern string) ([]int, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}
	var res []int
	for i, line := range lines {
		if re.MatchString(line) {
			res = append(res, i)
		}
	}
	return res, nil
}

func buildIntervals(matches []int, lines []string, cfg config) [][2]int {
	if len(matches) == 0 {
		return nil
	}
	before := cfg.before
	after := cfg.after
	if cfg.context > 0 {
		before = cfg.context
		after = cfg.context
	}
	var intervals [][2]int
	for _, pos := range matches {
		l := pos - before
		if l < 0 {
			l = 0
		}
		r := pos + after
		if r >= len(lines) {
			r = len(lines) - 1
		}
		if len(intervals) == 0 || l > intervals[len(intervals)-1][1]+1 {
			intervals = append(intervals, [2]int{l, r})
		} else {
			if r > intervals[len(intervals)-1][1] {
				intervals[len(intervals)-1][1] = r
			}
		}
	}
	return intervals
}

func printResult(lines []string, intervals [][2]int, cfg config) {
	for _, iv := range intervals {
		for i := iv[0]; i <= iv[1]; i++ {
			if cfg.lineNum {
				fmt.Printf("%d:", i+1)
			}
			fmt.Println(lines[i])
		}
	}
}
