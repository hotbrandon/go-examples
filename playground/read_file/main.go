package main

import (
	"io"
	"log"
	"os"
)

// On Windows (and most OS-level file readers), n > 0, err == io.EOF usually
// does not happen — instead, EOF appears on the next read.
// Still, the Go idiom accounts for both possibilities, since other readers
// (like bytes.Reader, strings.Reader, or compressed readers) do return n > 0, err == io.EOF
func main() {
	if len(os.Args) < 2 {
		log.Fatal("usage: go run main.go <filename>")
	}
	filename := os.Args[1]

	buffer := make([]byte, 149)
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	for {
		n, err := file.Read(buffer)
		//  ref. https://pkg.go.dev/io#Reader
		if n > 0 {
			os.Stdout.Write(buffer[:n])
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}
	}
}
