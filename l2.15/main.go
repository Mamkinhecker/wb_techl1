package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
)

type builtinFunc func(args []string) error

var builtins = map[string]builtinFunc{
	"cd":   builtinCd,
	"pwd":  builtinPwd,
	"echo": builtinEcho,
	"kill": builtinKill,
	"ps":   builtinPs,
}

func main() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		for range sigChan {
			fmt.Println("\nInterrupted")
		}
	}()

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("$ ")
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if line == "" {
			fmt.Print("$ ")
			continue
		}
		if err := executeLine(line); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		fmt.Print("$ ")
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func executeLine(line string) error {
	commands, err := parsePipeline(line)
	if err != nil {
		return err
	}
	return executePipeline(commands)
}

func parsePipeline(line string) ([][]string, error) {
	parts := strings.Split(line, "|")
	var res [][]string
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			return nil, errors.New("empty command in pipeline")
		}
		args, err := splitArgs(part)
		if err != nil {
			return nil, fmt.Errorf("parsing arguments: %w", err)
		}
		res = append(res, args)
	}
	return res, nil
}

func executePipeline(commands [][]string) error {
	if len(commands) == 0 {
		return nil
	}
	if len(commands) == 1 {
		return executeCommand(commands[0], os.Stdin, os.Stdout)
	}

	pipes := make([]*os.File, 2*(len(commands)-1))
	for i := 0; i < len(commands)-1; i++ {
		r, w, err := os.Pipe()
		if err != nil {
			return err
		}
		pipes[2*i] = r
		pipes[2*i+1] = w
	}
	cmds := make([]*exec.Cmd, len(commands))
	for i, args := range commands {
		var stdin io.Reader = os.Stdin
		var stdout io.Writer = os.Stdout
		if i > 0 {
			stdin = pipes[2*(i-1)]
		}
		if i < len(commands)-1 {
			stdout = pipes[2*i+1]
		}
		cmd, err := buildCommand(args, stdin, stdout)
		if err != nil {
			return err
		}
		if cmd == nil {
			return fmt.Errorf("builtin command %s cannot be used in pipeline", args[0])
		}
		cmds[i] = cmd
	}

	for _, cmd := range cmds {
		if err := cmd.Start(); err != nil {
			return err
		}
	}

	for i := 0; i < len(commands)-1; i++ {
		pipes[2*i+1].Close()
	}

	var lastErr error
	for _, cmd := range cmds {
		if err := cmd.Wait(); err != nil {
			lastErr = err
		}
	}

	for i := 0; i < len(commands)-1; i++ {
		pipes[2*i].Close()
	}
	return lastErr
}

func buildCommand(args []string, stdin io.Reader, stdout io.Writer) (*exec.Cmd, error) {
	if len(args) == 0 {
		return nil, errors.New("empty command")
	}
	cmdName := args[0]
	if _, ok := builtins[cmdName]; ok {
		return nil, nil
	}
	cmd := exec.Command(cmdName, args[1:]...)
	cmd.Stdin = stdin
	cmd.Stdout = stdout
	cmd.Stderr = os.Stderr
	return cmd, nil
}

func executeCommand(args []string, stdin io.Reader, stdout io.Writer) error {
	if len(args) == 0 {
		return nil
	}
	cmdName := args[0]
	if fn, ok := builtins[cmdName]; ok {
		return fn(args[1:])
	}
	cmd := exec.Command(cmdName, args[1:]...)
	cmd.Stdin = stdin
	cmd.Stdout = stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func builtinCd(args []string) error {
	var dir string
	if len(args) == 0 {
		var err error
		dir, err = os.UserHomeDir()
		if err != nil {
			return err
		}
	} else {
		dir = args[0]
	}
	return os.Chdir(dir)
}

func builtinPwd(args []string) error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	fmt.Println(dir)
	return nil
}

func builtinEcho(args []string) error {
	fmt.Println(strings.Join(args, " "))
	return nil
}

func builtinKill(args []string) error {
	if len(args) == 0 {
		return errors.New("kill: missing PID")
	}
	pid, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("kill: invalid PID %q", args[0])
	}
	proc, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	return proc.Signal(syscall.SIGTERM)
}

func builtinPs(args []string) error {
	cmd := exec.Command("ps", "aux")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func splitArgs(s string) ([]string, error) {
	var args []string
	var cur []rune
	inQuote := false
	quoteChar := rune(0)
	for _, ch := range s {
		if !inQuote && (ch == ' ' || ch == '\t') {
			if len(cur) > 0 {
				args = append(args, string(cur))
				cur = nil
			}
			continue
		}
		if !inQuote && (ch == '"' || ch == '\'') {
			inQuote = true
			quoteChar = ch
			continue
		}
		if inQuote && ch == quoteChar {
			inQuote = false
			quoteChar = 0
			continue
		}
		cur = append(cur, ch)
	}
	if len(cur) > 0 {
		args = append(args, string(cur))
	}
	if inQuote {
		return nil, fmt.Errorf("unclosed quote at position %d", len(s)-1)
	}
	return args, nil
}
