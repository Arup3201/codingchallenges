package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"slices"
	"unicode/utf8"
)

func main() {
	/*
		Usage: wc -c -m -w -l FILE
	*/

	var options []string
	var filename string

	options = os.Args[1 : len(os.Args)-1]
	filename = os.Args[len(os.Args)-1]

	var f *os.File
	var err error
	f, err = os.OpenFile(filename, os.O_RDONLY, os.ModePerm)
	if err != nil {
		fmt.Printf("Error while opening file: %s\n", err)
		return
	}
	defer f.Close()

	_, err = f.Seek(0, io.SeekStart)
	if err != nil {
		fmt.Printf("Error encountered while starting: %s\n", err)
		return
	}

	var container = make([]byte, 1)
	var bytesCnt,
		charsCnt,
		wordsCnt,
		linesCnt uint64
	var ch []byte // single or multi-byte character (UTF-8)
	var prevByte byte = 0
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
			if prevByte != 0 && !bytes.Contains(
				SPACE_BYTES,
				[]byte{prevByte},
			) {
				wordsCnt += 1
			}

			break
		} else if err != nil {
			fmt.Printf("Error encountered while scanning the file: %s\n", err)
			return
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

	if len(options) == 0 || slices.Contains(options, "-c") {
		fmt.Printf("c:%d ", bytesCnt)
	}
	if slices.Contains(options, "-m") {
		fmt.Printf("m:%d ", charsCnt)
	}
	if len(options) == 0 || slices.Contains(options, "-w") {
		fmt.Printf("w:%d ", wordsCnt)
	}
	if len(options) == 0 || slices.Contains(options, "-l") {
		fmt.Printf("l:%d ", linesCnt)
	}

	fmt.Printf("%s\n", filename)
}
