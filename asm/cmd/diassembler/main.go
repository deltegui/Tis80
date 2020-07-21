package main

import (
	"log"
	"os"
	"tisasm"
)

func main() {
	binaryFile := tisasm.OpenFile(getSourcePath())
	defer binaryFile.Close()
	diassembler := tisasm.NewDiassembler(binaryFile)
	diassembler.Diasemble()
}

func getSourcePath() string {
	if len(os.Args) != 2 {
		log.Fatalln("You should provide an assembly file")
	}
	return os.Args[1]
}
