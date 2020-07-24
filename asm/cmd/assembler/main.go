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
	return strings.Replace(inputPath, ".asm", ".rom", 1)
}

func scannerFromFile(file *os.File) tisasm.Scanner {
	return tisasm.NewFileScanner(bufio.NewReader(file))
}

func main() {
	path := getSourcePath()
	file := tisasm.OpenFile(path)
	defer file.Close()
	outputFile := tisasm.CreateFile(generateOutputFile(path))
	defer outputFile.Close()
	tagReader := tisasm.NewTagReader(scannerFromFile(file))
	tags := tagReader.GetTags()
	file.Seek(0, 0)
	parser := tisasm.NewParser(scannerFromFile(file), outputFile, tags)
	parser.Parse()
}
