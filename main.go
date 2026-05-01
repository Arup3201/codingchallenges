package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"unicode/utf8"
)

const (
	BYTE_FLAG      = "-c"
	CHARACTER_FLAG = "-m"
	WORD_FLAG      = "-w"
	LINE_FLAG      = "-l"
)

type wcStruct struct {
	bytes,
	chars,
	words,
	lines uint64
}

func getCounts(f *os.File) (*wcStruct, error) {
	var err error
	var container = make([]byte, 1)
	var bytesCnt,
		charsCnt,
		wordsCnt,
		linesCnt uint64
	var ch []byte         // store bytes to check if they make valid UTF-8 character
	var prevByte byte = 0 // tracks previous byte to detect word or not
	var SPACE_BYTES = []byte{
		'\t',
		'\n',
		'\v',
		'\f',
		'\r',
		' ',
		133, // NEW LINE
		160, // NO-BREAK SPACE
	}
	for {
		_, err = f.Read(container)
		if err == io.EOF {
			// the last word ending with EOF
			if prevByte != 0 && !bytes.Contains(
				SPACE_BYTES,
				[]byte{prevByte},
			) {
				wordsCnt += 1
			}

			break
		} else if err != nil {
			return nil, fmt.Errorf("Read: %s\n", err)
		}

		bytesCnt += 1

		ch = append(ch, container[0])
		if utf8.FullRune(ch) {
			charsCnt += 1
			ch = []byte{}
		}

		if prevByte != 0 && !bytes.Contains(
			SPACE_BYTES,
			[]byte{prevByte},
		) && bytes.Contains(
			SPACE_BYTES,
			container,
		) {
			wordsCnt += 1
		}

		if container[0] == '\n' {
			linesCnt += 1
		}

		prevByte = container[0]
	}

	return &wcStruct{
		bytes: bytesCnt,
		chars: charsCnt,
		words: wordsCnt,
		lines: linesCnt,
	}, nil
}

func main() {
	var (
		hasByteFlag      = false
		hasCharacterFlag = false
		hasWordFlag      = false
		hasLineFlag      = false
	)

	flag.BoolVar(&hasByteFlag, "c", false, "print bytes count")
	flag.BoolVar(&hasCharacterFlag, "m", false, "print characters count")
	flag.BoolVar(&hasWordFlag, "w", false, "print words count")
	flag.BoolVar(&hasLineFlag, "l", false, "print lines count")

	flag.Parse()

	if !(hasByteFlag || hasCharacterFlag || hasWordFlag || hasLineFlag) {
		hasByteFlag = true
		hasWordFlag = true
		hasLineFlag = true
	}

	var filenames = flag.CommandLine.Args()

	var f *os.File
	var err error
	var wcCnts *wcStruct
	for _, fn := range filenames {
		f, err = os.OpenFile(fn, os.O_RDONLY, os.ModePerm)
		if err != nil {
			fmt.Printf("OpenFile: %s\n", err)
			continue
		}
		defer f.Close()

		_, err = f.Seek(0, io.SeekStart)
		if err != nil {
			fmt.Printf("Seek: %s\n", err)
			continue
		}

		wcCnts, err = getCounts(f)
		if err != nil {
			fmt.Printf("getCounts: %s\n", err)
			continue
		}

		if hasByteFlag {
			fmt.Printf("c:%d ", wcCnts.bytes)
		}
		if hasCharacterFlag {
			fmt.Printf("m:%d ", wcCnts.chars)
		}
		if hasWordFlag {
			fmt.Printf("w:%d ", wcCnts.words)
		}
		if hasLineFlag {
			fmt.Printf("l:%d ", wcCnts.lines)
		}

		fmt.Printf("%s\n", fn)
	}

	// Take from STDIN
	if len(filenames) == 0 {
		wcCnts, err = getCounts(os.Stdin)
		if err != nil {
			fmt.Printf("getCounts: %s\n", err)
			return
		}

		if hasByteFlag {
			fmt.Printf("c:%d ", wcCnts.bytes)
		}
		if hasCharacterFlag {
			fmt.Printf("m:%d ", wcCnts.chars)
		}
		if hasWordFlag {
			fmt.Printf("w:%d ", wcCnts.words)
		}
		if hasLineFlag {
			fmt.Printf("l:%d ", wcCnts.lines)
		}

		fmt.Printf("\n")
	}
}
