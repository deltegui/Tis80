package tisasm

import (
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
)

type Parser struct {
	scanner Scanner
	out     *os.File
	tags    map[string]string
}

func NewParser(scanner Scanner, out *os.File, tags map[string]string) Parser {
	return Parser{
		scanner,
		out,
		tags,
	}
}

func (prs Parser) Parse() {
	token := prs.scanner.Scan()
	if token.TokenType != TokenSection {
		ShowError("Expected start of section in top of file")
	}
	switch token.Literal {
	case ".data":
		prs.parseDataSection()
	case ".code":
		prs.parseCodeSection()
	default:
		ShowErrorf("Unknown section %s in top of file", token.Literal)
	}
}

func (prs Parser) parseDataSection() {
	prs.emitDataSection()
	token := prs.scanner.Scan()
	for token.IsCorrect() && !token.IsType(TokenSection) {
		if !token.IsType(TokenMemory) {
			ShowError("Expected memory address inside data section")
		}
		prs.emitMemory(token)
		token = prs.scanner.Scan()
		switch token.TokenType {
		case TokenString:
			prs.emitBytes(0x01)
			prs.emitASCII(token)
		case TokenNumber:
			prs.emitBytes(0x02)
			prs.emitNumber(token)
		default:
			ShowError("Expected number or string after memory address inside data section")
		}
		token = prs.scanner.Scan()
	}
	if !token.IsSection(".code") {
		ShowError("Expected .code section after .data section")
	}
	prs.emitBytes(0x00, 0x00, 0x00)
	prs.parseCodeSection()
}

func (prs Parser) parseCodeSection() {
	prs.emitCodeSection()
	token := prs.scanner.Scan()
	for token.IsCorrect() {
		switch token.TokenType {
		case TokenTag: // Do nothing
		case TokenInstruction:
			prs.parseInstruction(token)
		default:
			ShowErrorf("Expected instruction but have %s", token)
		}
		token = prs.scanner.Scan()
	}
	if token.IsType(TokenError) {
		ShowErrorf("Error while reading section code %s", token)
	}
}

func (prs *Parser) parseInstruction(token Token) {
	instruction := token.AsInstruction()
	prs.emitBytes(instruction.OpCode)
	instruction.ParseParams(*prs)
}

func (prs *Parser) emitCodeSection() {
	prs.emitSectionStart()
	prs.emitBytes(0x01)
	codeMemoryStart := prs.scanner.Scan()
	if !codeMemoryStart.IsType(TokenMemory) {
		ShowError("Expected code section start memory direction after .code section")
	}
	prs.emitMemory(codeMemoryStart)
}

func (prs Parser) emitDataSection() {
	prs.emitSectionStart()
	prs.emitBytes(0x00)
}

func (prs Parser) emitSectionStart() {
	prs.emitBytes(0xff, 0xfe, 0xfe, 0xff)
}

func (prs Parser) emitASCII(token Token) {
	data := []byte(token.Literal)
	prs.emitBytes(data...)
	prs.emitBytes(0x00)
}

func (prs Parser) emitNumber(token Token) {
	if !token.IsType(TokenNumber) {
		ShowError("Expected token to be number")
	}
	integer, err := strconv.Atoi(token.Literal)
	if err != nil {
		ShowErrorf("Error while decoding as number literal: %s", token.Literal)
	}
	if integer >= 256 {
		ShowError("Integer must be under 256")
	}
	b := byte(integer & 0xff)
	prs.emitBytes(b)
}

func (prs Parser) emitJumpDest(token Token) {
	switch token.TokenType {
	case TokenMemory:
		prs.emitMemory(token)
	case TokenInstruction:
		prs.emitTag(token)
	default:
		ShowError("Expected jump destination to be a tag or memory address")
	}
}

func (prs Parser) emitTag(token Token) {
	if token.TokenType != TokenInstruction {
		ShowError("Expected tag")
	}
	val, ok := prs.tags[token.Literal]
	if !ok {
		ShowErrorf("Expected tag %s to be defined", token.Literal)
	}
	prs.emitMemory(Token{
		Literal:   val,
		TokenType: TokenMemory,
	})
}

func (prs Parser) emitRegister(token Token) {
	integer, err := strconv.Atoi(token.Literal)
	if err != nil {
		ShowErrorf("Expected register number, got %s", token)
	}
	fixed := fmt.Sprintf("%02d", integer)
	prs.emitHex(fixed, 1)
}

func (prs Parser) emitMemory(token Token) {
	prs.emitHex(token.Literal, 2)
}

func (prs Parser) emitHex(literal string, length int) {
	address, err := hex.DecodeString(literal)
	if err != nil {
		ShowErrorf("Error while decoding as hexadecimal literal in emitHex: %s", literal)
	}
	if len(address) != length {
		ShowErrorf("Expected hexadecimal to be %d length", length)
	}
	prs.emitBytes(address...)
}

func (prs Parser) emitBytes(bytes ...byte) {
	_, err := prs.out.Write(bytes)
	if err != nil {
		panic(err)
	}
}
