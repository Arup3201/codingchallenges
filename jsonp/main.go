package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

const (
	OPEN  = '{'
	CLOSE = '}'
)

const (
	TK_OPEN = iota
	TK_CLOSE
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: jsonp <input.json>\n")
		os.Exit(1)
	}

	jsonFile := os.Args[1]
	splits := strings.Split(jsonFile, ".")
	if strings.ToLower(splits[len(splits)-1]) != "json" {
		fmt.Printf("Invalid file format\n")
		fmt.Printf("Usage: jsonp <input.json>\n")
		os.Exit(1)
	}

	f, err := os.OpenFile(jsonFile, os.O_RDONLY, os.ModePerm)
	if err != nil {
		fmt.Printf("OpenFile: %s\n", err)
		os.Exit(1)
	}

	tokens := []int{}
	rd := bufio.NewReader(f)
	for {
		ch, _, err := rd.ReadRune()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			fmt.Printf("ReadRune: %s\n", err)
			os.Exit(1)
		}

		switch ch {
		case OPEN:
			tokens = append(tokens, TK_OPEN)
		case CLOSE:
			tokens = append(tokens, TK_CLOSE)
		}
	}

	if len(tokens) < 2 {
		fmt.Printf("Invalid JSON format\n")
		os.Exit(1)
	}

	n := len(tokens)
	if tokens[0] != TK_OPEN || tokens[n-1] != TK_CLOSE {
		fmt.Printf("Invalid JSON format\n")
		os.Exit(1)
	}

	os.Exit(0)
}
