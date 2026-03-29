package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	var (
		fields    string        // -f
		delimiter string = "\t" // -d
		sOnly     bool          // -s
	)

	args := os.Args[1:]
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-f":
			if i+1 >= len(args) {
				die("missing argument for -f")
			}
			fields = args[i+1]
			i++

		case "-d":
			if i+1 >= len(args) {
				die("missing argument for -d")
			}
			delimiter = args[i+1]
			i++

		case "-s":
			sOnly = true

		default:
			if len(args[i]) > 0 && args[i][0] == '-' {
				die("unknown flag: " + args[i])
			}
		}
	}

	if fields == "" {
		fields = "1"
	}

	fieldSet, err := parseFields(fields)
	if err != nil {
		die("invalid fields: " + err.Error())
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()

		if sOnly && !strings.Contains(line, delimiter) {
			continue
		}

		parts := strings.Split(line, delimiter)
		var output []string
		for i := 0; i < len(parts); i++ {
			if fieldSet[i+1] {
				output = append(output, parts[i])
			}
		}
		if len(output) > 0 {
			fmt.Println(strings.Join(output, delimiter))
		}
	}

	if err := scanner.Err(); err != nil {
		die("read error: " + err.Error())
	}
}

func parseFields(spec string) (map[int]bool, error) {
	fields := make(map[int]bool)
	re := regexp.MustCompile(`^\s*(\d+)\s*-\s*(\d+)\s*$`)

	for _, p := range strings.Split(spec, ",") {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}

		if re.MatchString(p) {
			m := re.FindStringSubmatch(p)
			start, errA := strconv.Atoi(m[1])
			end, errB := strconv.Atoi(m[2])
			if errA != nil || errB != nil || start <= 0 || end <= 0 || start > end {
				return nil, fmt.Errorf("invalid range: %s", p)
			}
			for i := start; i <= end; i++ {
				fields[i] = true
			}
		} else {
			n, err := strconv.Atoi(p)
			if err != nil || n <= 0 {
				return nil, fmt.Errorf("invalid field: %s", p)
			}
			fields[n] = true
		}
	}

	return fields, nil
}

func die(msg string) {
	fmt.Fprintf(os.Stderr, "cututil: %s\n", msg)
	os.Exit(1)
}
