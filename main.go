package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"slices"
	"unicode/utf8"
)

const (
	BYTE_FLAG      = "-c"
	CHARACTER_FLAG = "-m"
	WORD_FLAG      = "-w"
	LINE_FLAG      = "-l"
)

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

	var filenames []string

	var filterFn = func(s []string, cmp func(r string) bool) []string {
		ret := []string{}
		for _, r := range s {
			if cmp(r) {
				ret = append(ret, r)
			}
		}
		return ret
	}

	filenames = filterFn(os.Args[1:], func(r string) bool {
		return !slices.Contains([]string{
			BYTE_FLAG,
			CHARACTER_FLAG,
			WORD_FLAG,
			LINE_FLAG,
		}, r)
	})

	var f *os.File
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

		bytesCnt,
			charsCnt,
			wordsCnt,
			linesCnt = 0, 0, 0, 0
		ch = []byte{}
		prevByte = 0
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
				fmt.Printf("Read: %s\n", err)
				continue
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

		if hasByteFlag {
			fmt.Printf("c:%d ", bytesCnt)
		}
		if hasCharacterFlag {
			fmt.Printf("m:%d ", charsCnt)
		}
		if hasWordFlag {
			fmt.Printf("w:%d ", wordsCnt)
		}
		if hasLineFlag {
			fmt.Printf("l:%d ", linesCnt)
		}

		fmt.Printf("%s\n", fn)
	}

	// Take from STDIN
	if len(filenames) == 0 {
		bytesCnt,
			charsCnt,
			wordsCnt,
			linesCnt = 0, 0, 0, 0
		ch = []byte{}
		prevByte = 0
		for {
			_, err = os.Stdin.Read(container)
			if err == io.EOF {
				if prevByte != 0 && !bytes.Contains(
					SPACE_BYTES,
					[]byte{prevByte},
				) {
					wordsCnt += 1
				}

				break
			} else if err != nil {
				fmt.Printf("Read: %s\n", err)
				continue
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

		if hasByteFlag {
			fmt.Printf("c:%d ", bytesCnt)
		}
		if hasCharacterFlag {
			fmt.Printf("m:%d ", charsCnt)
		}
		if hasWordFlag {
			fmt.Printf("w:%d ", wordsCnt)
		}
		if hasLineFlag {
			fmt.Printf("l:%d ", linesCnt)
		}

		fmt.Printf("\n")
	}
}
