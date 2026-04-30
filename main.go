package main

import (
	"fmt"
	"io"
	"os"
	"slices"
)

func main() {
	/*
		Usage: wc -c -l FILE
	*/

	if len(os.Args) >= 3 {
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

		var b = make([]byte, 1)
		var bytesCnt, linesCnt uint64
		for {
			_, err = f.Read(b)
			if err == io.EOF {
				break
			} else if err != nil {
				fmt.Printf("Error encountered while scanning the file: %s\n", err)
				return
			}

			if b[0] == '\n' {
				linesCnt += 1
			}

			bytesCnt += 1
		}

		if slices.Contains(options, "-c") {
			fmt.Printf("%d ", bytesCnt)
		}
		if slices.Contains(options, "-l") {
			fmt.Printf("%d ", linesCnt)
		}

		fmt.Printf(" %s\n", filename)
	}
}
