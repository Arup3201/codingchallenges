package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

var ErrStackUnderflow = errors.New("stack underflow")

var UsageInstruction = "Usage: jsonp <input.json>\n"

const (
	LEFT_BRACE   = '{'
	RIGHT_BRACE  = '}'
	DOUBLE_QUOTE = '"'
	COLON        = ':'
	COMMA        = ','
)

const (
	TOKEN_OBJ_OPEN = iota
	TOKEN_OBJ_CLOSE
	TOKEN_KEY
	TOKEN_COLON
	TOKEN_VALUE
	TOKEN_COMMA
)

func exitWithError(format string, a ...string) {
	if len(a) == 0 {
		fmt.Println(format)
	} else {
		fmt.Printf(format, a)
	}
	os.Exit(1)
}

func main() {
	if len(os.Args) < 2 {
		exitWithError(UsageInstruction)
	}

	jsonFile := os.Args[1]
	splits := strings.Split(jsonFile, ".")
	if strings.ToLower(splits[len(splits)-1]) != "json" {
		exitWithError("only json file can be parsed\n%s", UsageInstruction)
	}

	f, err := os.OpenFile(jsonFile, os.O_RDONLY, os.ModePerm)
	if err != nil {
		exitWithError("OpenFile: %s\n", err.Error())
	}

	var tokens = []int{}
	var tokenCount = 0
	var pushToken = func(tkn int) {
		tokens = append(tokens, tkn)
		tokenCount += 1
	}

	var ch rune
	rd := bufio.NewReader(f)
	for {
		ch, _, err = rd.ReadRune()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			exitWithError("ReadRune: %s\n", err.Error())
		}

		switch ch {
		case LEFT_BRACE:
			pushToken(TOKEN_OBJ_OPEN)
		case RIGHT_BRACE:
			pushToken(TOKEN_OBJ_CLOSE)
		case DOUBLE_QUOTE:
			/*
				Start reading the string after double-quotes(")
				untill another double-quote is found.
			*/
			ch, _, err = rd.ReadRune()
			if err != nil {
				exitWithError("missing closing \" for string value")
			}
			for ch != DOUBLE_QUOTE {
				ch, _, err = rd.ReadRune()
				if err != nil {
					exitWithError("missing closing \" for string value")
				}
				if ch == DOUBLE_QUOTE {
					break
				}
			}

			// string before colon is key and after colon is value
			if tokens[tokenCount-1] == TOKEN_COLON {
				pushToken(TOKEN_VALUE)
			} else {
				pushToken(TOKEN_KEY)
			}
		case COLON:
			pushToken(TOKEN_COLON)
		case COMMA:
			pushToken(TOKEN_COMMA)
		}
	}

	if tokenCount < 2 ||
		(tokens[0] != TOKEN_OBJ_OPEN && tokens[tokenCount-1] != TOKEN_OBJ_CLOSE) {
		exitWithError("invalid json")
	}

	for i := 1; i < tokenCount; i++ {
		if (tokens[i-1] == TOKEN_OBJ_OPEN &&
			(tokens[i] != TOKEN_KEY && tokens[i] != TOKEN_OBJ_CLOSE)) || // \{"key" / {}
			(tokens[i-1] == TOKEN_KEY && tokens[i] != TOKEN_COLON) || // key:
			(tokens[i-1] == TOKEN_COLON && tokens[i] != TOKEN_VALUE) || // :value
			(tokens[i-1] == TOKEN_VALUE && (tokens[i] != TOKEN_COMMA && tokens[i] != TOKEN_OBJ_CLOSE)) || // value, / value}
			(tokens[i-1] == TOKEN_COMMA && (tokens[i] != TOKEN_KEY)) || // ,key
			false {
			exitWithError("invalid json")
		}

	}

	fmt.Println("valid json")
}
