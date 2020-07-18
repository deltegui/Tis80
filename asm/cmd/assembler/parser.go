package main

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"strings"
	"tisasm"
)

type parser struct {
	scanner   scanner
	out       *os.File
	tags      map[string]string
	codeStart int
	line      int
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

func (prs parser) emitNumber(token Token) {
	if token.TokenType != TokenNumber {
		tisasm.ShowError("Expected token to be number")
	}
	prs.emitRegister(token)
}

func (prs parser) emitMemory(token Token) {
	prs.emitHex(token.Literal, 2)
}

func (prs parser) emitASCII(token Token) {
	data := []byte(token.Literal)
	prs.emitBytes(data...)
	prs.emitBytes(0x00)
}

func (prs parser) emitBytes(bytes ...byte) {
	_, err := prs.out.Write(bytes)
	if err != nil {
		panic(err)
	}
}

func (prs parser) emitJumpDest(token Token) {
	switch token.TokenType {
	case TokenMemory:
		prs.emitMemory(token)
	case TokenInstruction:
		prs.emitTag(token)
	default:
		tisasm.ShowError("Expected jump destination to be a tag or memory address")
	}
}

func (prs parser) emitSectionStart() {
	prs.emitBytes(0xff, 0xfe, 0x00, 0xfe, 0xff)
}

func (prs parser) emitDataSection() {
	prs.emitSectionStart()
	prs.emitBytes(0x00)
}

func (prs *parser) emitCodeSection() {
	prs.emitSectionStart()
	prs.emitBytes(0x01)
	codeMemoryStart := prs.scanner.Scan()
	if codeMemoryStart.TokenType != TokenMemory {
		tisasm.ShowError("Expected code section start memory direction after .code section")
	}
	prs.emitMemory(codeMemoryStart)
	prs.setLineStart(codeMemoryStart.Literal)
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
		fmt.Printf("[%s] %s\n", token.TokenType, token.Literal)
		switch token.TokenType {
		case TokenTag:
			prs.defineTag(token.Literal)
		case TokenInstruction:
			prs.parseInstruction(token)
		default:
			tisasm.ShowError(fmt.Sprintf("Expected instruction but have [%s] %s", token.TokenType, token.Literal))
		}
		token = prs.scanner.Scan()
	}
	if token.TokenType == TokenError {
		fmt.Printf("Error while reading section code %s: %s\n", token.TokenType, token.Literal)
	}
}

func (prs parser) defineTag(tag string) {
	prs.tags[tag] = fmt.Sprintf("%04x", prs.codeStart+prs.line)
}

func (prs parser) emitTag(token Token) {
	if token.TokenType != TokenInstruction {
		tisasm.ShowError("Expected tag")
	}
	val, ok := prs.tags[token.Literal]
	if !ok {
		tisasm.ShowError(fmt.Sprintf("Expected tag %s to be defined", token.Literal))
	}
	prs.emitMemory(Token{
		Literal:   val,
		TokenType: TokenMemory,
	})
}

func (prs *parser) setLineStart(memory string) {
	bytes, err := hex.DecodeString(memory)
	if err != nil {
		tisasm.ShowError("Line start invaild")
	}
	prs.codeStart = int(binary.BigEndian.Uint16(bytes))
}

func (prs *parser) advanceLines(number int) {
	prs.line += number
}

func (prs *parser) parseInstruction(token Token) {
	instruction := strings.ToLower(token.Literal)
	switch instruction {
	case "add":
		prs.emitBytes(0x00)
		prs.emitRegister(prs.scanner.Scan())
		prs.advanceLines(2)
	case "movi":
		prs.emitBytes(0x01)
		prs.emitRegister(prs.scanner.Scan())
		prs.emitNumber(prs.scanner.Scan())
		prs.advanceLines(3)
	case "jne":
		prs.emitBytes(0x02)
		prs.emitJumpDest(prs.scanner.Scan())
		prs.advanceLines(3)
	case "hlt":
		prs.emitBytes(0x03)
		prs.advanceLines(1)
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
