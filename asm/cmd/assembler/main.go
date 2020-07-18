package main

import (
	"bufio"
	"log"
	"os"
	"strings"
	"tisasm"
)

func getSourcePath() string {
	if len(os.Args) != 2 {
		log.Fatalln("You should provide an assembly file")
	}
	return os.Args[1]
}

func generateOutputFile(inputPath string) string {
	return strings.Replace(inputPath, ".asm", ".bin", 1)
}

func main() {
	path := getSourcePath()
	file := tisasm.OpenFile(path)
	defer file.Close()
	outputFile := tisasm.CreateFile(generateOutputFile(path))
	defer outputFile.Close()
	scn := newScanner(bufio.NewReader(file))
	parser := parser{
		scanner: scn,
		out:     outputFile,
		tags:    make(map[string]string),
		line:    0,
	}
	parser.parse()
}
