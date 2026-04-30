package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	/*
		Usage: wc -c FILE
	*/

	if len(os.Args) >= 3 {
		var option string
		var filename string

		option = os.Args[1]
		filename = os.Args[2]

		var f *os.File
		var err error
		f, err = os.OpenFile(filename, os.O_RDONLY, os.ModePerm)
		if err != nil {
			fmt.Printf("Error while opening file: %s\n", err)
			return
		}
		defer f.Close()

		if option == "-c" {

			_, err = f.Seek(0, io.SeekStart)
			if err != nil {
				fmt.Printf("Error encountered while starting: %s\n", err)
				return
			}

			var b = make([]byte, 1)
			var cnt uint64 = 0
			for {
				_, err = f.Read(b)
				if err == io.EOF {
					break
				} else if err != nil {
					fmt.Printf("Error encountered while scanning the file: %s\n", err)
					return
				}
				cnt += 1
			}

			fmt.Printf("%d %s\n", cnt, filename)
		}
	}
}
