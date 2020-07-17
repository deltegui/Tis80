package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"tisasm"
)

type parser struct {
	scanner scanner
	out     *os.File
}

func (prs parser) isCorrect(token Token) bool {
	return token.TokenType != TokenEof && token.TokenType != TokenError
}

func (prs parser) emitHex(literal string, length int) {
	address, err := hex.DecodeString(literal)
	if err != nil {
		fmt.Println("Strange error while decoding hex")
		panic(err)
	}
	if len(address) != length {
		tisasm.ShowError(fmt.Sprintf("Expected hexadecimal to be %d length", length))
	}
	prs.emitBytes(address...)
}

func (prs parser) emitRegister(token Token) {
	integer, err := strconv.Atoi(token.Literal)
	if err != nil {
		tisasm.ShowError("Expected register number.")
	}
	fixed := fmt.Sprintf("%02d", integer)
	prs.emitHex(fixed, 1)
}

func (prs parser) emitMemory(token Token) {
	prs.emitHex(token.Literal, 2)
}

func (prs parser) emitASCII(token Token) {
	data := []byte(token.Literal)
	prs.emitBytes(data...)
	if token.TokenType == TokenString {
		prs.emitBytes(0x00)
	}
}

func (prs parser) emitBytes(bytes ...byte) {
	_, err := prs.out.Write(bytes)
	if err != nil {
		panic(err)
	}
}

func (prs parser) emitSectionStart() {
	prs.emitBytes(0xff, 0xfe, 0x00, 0xfe, 0xff)
}

func (prs parser) emitDataSection() {
	prs.emitSectionStart()
	prs.emitBytes(0x00)
}

func (prs parser) emitCodeSection() {
	prs.emitSectionStart()
	prs.emitBytes(0x01)
	codeMemoryStart := prs.scanner.Scan()
	if codeMemoryStart.TokenType != TokenMemory {
		tisasm.ShowError("Expected code section start memory direction after .code section")
	}
	prs.emitMemory(codeMemoryStart)
}

func (prs parser) parse() {
	token := prs.scanner.Scan()
	if token.TokenType != TokenSection {
		tisasm.ShowError("Expected start of section in top of file")
	}
	switch token.Literal {
	case ".code":
		prs.parseCodeSection()
	case ".data":
		prs.parseDataSection()
	default:
		tisasm.ShowError(fmt.Sprintf("Unknown section %s in top of file", token.Literal))
	}
}

func (prs parser) parseCodeSection() {
	prs.emitCodeSection()
	token := prs.scanner.Scan()
	for prs.isCorrect(token) {
		if token.TokenType != TokenInstruction {
			tisasm.ShowError(fmt.Sprintf("Expected instruction but have [%s] %s", token.TokenType, token.Literal))
		}
		prs.parseInstruction(token)
		token = prs.scanner.Scan()
	}
}

func (prs parser) parseInstruction(token Token) {
	instruction := strings.ToLower(token.Literal)
	switch instruction {
	case "add":
		prs.emitBytes(0x00)
		prs.emitRegister(prs.scanner.Scan())
	default:
		tisasm.ShowError(fmt.Sprintf("Undefined instruction %s!", instruction))
	}
}

func (prs parser) parseDataSection() {
	prs.emitDataSection()
	token := prs.scanner.Scan()
	for prs.isCorrect(token) && token.TokenType != TokenSection {
		if token.TokenType != TokenMemory {
			tisasm.ShowError("Expected memory address inside data section")
		}
		prs.emitMemory(token)
		token = prs.scanner.Scan()
		if token.TokenType != TokenString && token.TokenType != TokenNumber {
			tisasm.ShowError("Expected number or string after memory address inside data section")
		}
		prs.emitASCII(token)
		token = prs.scanner.Scan()
	}
	if token.TokenType != TokenSection && token.Literal != ".code" {
		tisasm.ShowError("Expected .code section after .data section")
	}
	prs.parseCodeSection()
}

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
	}
	parser.parse()
}
