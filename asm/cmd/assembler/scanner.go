package main

import (
	"bufio"
	"io"
	"tisasm"
	"unicode"
)

type scanner struct {
	reader     *bufio.Reader
	word       []rune
	lastReaded rune
}

func newScanner(reader *bufio.Reader) scanner {
	return scanner{reader, []rune{}, ' '}
}

func (scn scanner) skipWhitespaces() {
	for {
		if scn.isAtEnd() {
			return
		}
		c := scn.current()
		if c == '\t' || c == '\r' || c == ' ' || c == '\n' {
			scn.consume()
		} else if c == ';' {
			scn.consumeUntil('\n')
		} else {
			return
		}
	}
}

func (scn scanner) consumeUntil(end rune) {
	c := scn.current()
	for !scn.isAtEnd() && c != end {
		c = scn.consume()
	}
}

func (scn *scanner) consume() rune {
	readed := scn.skip()
	scn.word = append(scn.word, readed)
	return readed
}

func (scn scanner) skip() rune {
	readed, _, err := scn.reader.ReadRune()
	if err != nil {
		panic(err)
	}
	scn.lastReaded = readed
	return readed
}

func (scn scanner) skipExpected(expected rune, msg string) {
	if expected != scn.skip() {
		tisasm.ShowError(msg)
	}
}

func (scn scanner) current() rune {
	runes, err := scn.reader.Peek(1)
	if err != nil {
		return scn.lastReaded
	}
	return rune(runes[0])
}

func (scn scanner) isNumeric() bool {
	return !scn.isAtEnd() && unicode.IsDigit(scn.current())
}

func (scn scanner) isLetter() bool {
	if scn.isAtEnd() {
		return false
	}
	c := scn.current()
	return unicode.IsLetter(c) && c != '"'
}

func (scn scanner) isInstruction() bool {
	return scn.isLetter() && !scn.isRegister()
}

func (scn scanner) isRegister() bool {
	c := scn.current()
	return c == 'r' || c == 'R'
}

func (scn scanner) scanNumber() Token {
	for scn.isNumeric() {
		scn.consume()
	}
	if scn.current() == '.' {
		scn.consume()
		for scn.isNumeric() {
			scn.consume()
		}
	}
	return scn.createToken(TokenNumber)
}

func (scn scanner) scanSection() Token {
	scn.consume() // Consume dot
	for scn.isLetter() {
		scn.consume()
	}
	return scn.createToken(TokenSection)
}

func (scn scanner) scanInstruction() Token {
	for scn.isLetter() {
		scn.consume()
	}
	return scn.createToken(TokenInstruction)
}

func (scn scanner) isAtEnd() bool {
	_, err := scn.reader.Peek(1)
	return err == io.EOF
}

func (scn scanner) scanString() Token {
	scn.skipExpected('"', "Expected '\"' at start of string")
	for !scn.isAtEnd() && scn.current() != '"' {
		if scn.current() == '\n' || scn.isAtEnd() {
			return scn.createError("Unterminated string")
		}
		scn.consume()
	}
	if !scn.isAtEnd() {
		scn.skipExpected('"', "Expected '\"' at end of string")
	}
	return scn.createToken(TokenString)
}

func (scn scanner) scanDirection() Token {
	scn.skipExpected('$', "Expected '$' at start of memory direction")
	for i := 0; i < 4; i++ {
		if scn.isAtEnd() || scn.current() == '\n' {
			return scn.createError("Unterminated memory direction")
		}
		scn.consume()
	}
	return scn.createToken(TokenMemory)
}

func (scn scanner) scanCharacter() Token {
	scn.skipExpected('\'', "Expected ''' at start of character")
	scn.consume()
	scn.skipExpected('"', "Expected '\"' at end of string")
	return scn.createToken(TokenChar)
}

func (scn scanner) scanRegister() Token {
	startRegister := scn.skip()
	if startRegister != 'R' && startRegister != 'r' {
		tisasm.ShowError("Expected register to start with 'R' or 'r'")
	}
	scn.consume()
	return scn.createToken(TokenRegister)
}

func (scn scanner) scanTag() Token {
	scn.skipExpected(':', "Expected ':' before tag")
	for !scn.isAtEnd() && scn.current() != '\n' {
		scn.consume()
	}
	return scn.createToken(TokenTag)
}

func (scn scanner) createToken(tokenType TokenType) Token {
	return Token{
		TokenType: tokenType,
		Literal:   string(scn.word),
	}
}

func (scn scanner) createError(msg string) Token {
	return Token{
		TokenType: TokenError,
		Literal:   msg,
	}
}

func (scn scanner) Scan() Token {
	scn.skipWhitespaces()
	if scn.isAtEnd() {
		return scn.createToken(TokenEof)
	}
	scn.word = []rune{}
	if scn.isNumeric() {
		return scn.scanNumber()
	}
	if scn.isInstruction() {
		return scn.scanInstruction()
	}
	if scn.isRegister() {
		return scn.scanRegister()
	}
	switch scn.current() {
	case '.':
		return scn.scanSection()
	case '"':
		return scn.scanString()
	case '$':
		return scn.scanDirection()
	case '\'':
		return scn.scanCharacter()
	case ':':
		return scn.scanTag()
	default:
		return scn.createError("Unknown token")
	}
}
