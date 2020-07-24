package tisasm

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
)

type tagReader struct {
	line            int
	codeStart       int
	tags            map[string]string
	scn             Scanner
	isInCodeSection bool
}

func NewTagReader(inner Scanner) tagReader {
	return tagReader{
		line:            0,
		codeStart:       0,
		tags:            make(map[string]string),
		scn:             inner,
		isInCodeSection: false,
	}
}

func (reader tagReader) GetTags() map[string]string {
	reader.readTags()
	return reader.tags
}

func (reader tagReader) readTags() {
	token := reader.scn.Scan()
	for token.IsCorrect() && !reader.isInCodeSection {
		if token.IsSection(".code") {
			reader.isInCodeSection = true
			reader.setCodeStart(reader.scn.Scan())
		}
		token = reader.scn.Scan()
	}
	for token.IsCorrect() && reader.isInCodeSection {
		size := reader.processToken(token)
		reader.checkMemoryLimit()
		reader.scn.Advance(size - 1)
		token = reader.scn.Scan()
	}
}

func (reader tagReader) checkMemoryLimit() {
	position := reader.codeStart + reader.line
	if position > MemoryLimit {
		ShowError("Memory limit exceed.")
	}
}

func (reader *tagReader) setCodeStart(token Token) {
	if !token.IsType(TokenMemory) {
		ShowErrorToken(token, "Expected code start to be a memory direction")
	}
	memory := token.Literal
	bytes, err := hex.DecodeString(memory)
	if err != nil {
		ShowErrorToken(token, "Invalid memory format")
	}
	reader.codeStart = int(binary.BigEndian.Uint16(bytes))
}

func (reader *tagReader) processToken(token Token) int {
	switch token.TokenType {
	case TokenTag:
		reader.defineTag(token)
		return 1
	case TokenInstruction:
		instruction := token.AsInstruction()
		reader.line += instruction.MemorySize
		return instruction.TokenSize
	default:
		return 0
	}
}

func (reader tagReader) defineTag(token Token) {
	reader.tags[token.Literal] = fmt.Sprintf("%04x", reader.codeStart+reader.line)
}
