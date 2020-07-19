package tisasm

import (
	"bufio"
	"io"
	"unicode"
)

type FileScanner struct {
	reader     *bufio.Reader
	word       []rune
	lastReaded rune
}

func NewFileScanner(reader *bufio.Reader) FileScanner {
	return FileScanner{reader, []rune{}, ' '}
}

func (scn FileScanner) skipWhitespaces() {
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

func (scn FileScanner) consumeUntil(end rune) {
	c := scn.current()
	for !scn.isAtEnd() && c != end {
		c = scn.consume()
	}
}

func (scn *FileScanner) consume() rune {
	readed := scn.skip()
	scn.word = append(scn.word, readed)
	return readed
}

func (scn FileScanner) skip() rune {
	readed, _, err := scn.reader.ReadRune()
	if err != nil {
		panic(err)
	}
	scn.lastReaded = readed
	return readed
}

func (scn FileScanner) skipExpected(expected rune, msg string) {
	if expected != scn.skip() {
		ShowError(msg)
	}
}

func (scn FileScanner) current() rune {
	runes, err := scn.reader.Peek(1)
	if err != nil {
		return scn.lastReaded
	}
	return rune(runes[0])
}

func (scn FileScanner) isNumeric() bool {
	return !scn.isAtEnd() && unicode.IsDigit(scn.current())
}

func (scn FileScanner) isLetter() bool {
	if scn.isAtEnd() {
		return false
	}
	c := scn.current()
	return unicode.IsLetter(c) && c != '"'
}

func (scn FileScanner) isInstruction() bool {
	return scn.isLetter() && !scn.isRegister()
}

func (scn FileScanner) isRegister() bool {
	c := scn.current()
	return c == 'r' || c == 'R'
}

func (scn FileScanner) scanNumber() Token {
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

func (scn FileScanner) scanSection() Token {
	scn.consume() // Consume dot
	for scn.isLetter() {
		scn.consume()
	}
	return scn.createToken(TokenSection)
}

func (scn FileScanner) scanInstruction() Token {
	for scn.isLetter() {
		scn.consume()
	}
	return scn.createToken(TokenInstruction)
}

func (scn FileScanner) isAtEnd() bool {
	_, err := scn.reader.Peek(1)
	return err == io.EOF
}

func (scn FileScanner) scanString() Token {
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

func (scn FileScanner) scanDirection() Token {
	scn.skipExpected('$', "Expected '$' at start of memory direction")
	for i := 0; i < 4; i++ {
		if scn.isAtEnd() || scn.current() == '\n' {
			return scn.createError("Unterminated memory direction")
		}
		scn.consume()
	}
	return scn.createToken(TokenMemory)
}

func (scn FileScanner) scanCharacter() Token {
	scn.skipExpected('\'', "Expected ''' at start of character")
	scn.consume()
	scn.skipExpected('"', "Expected '\"' at end of string")
	return scn.createToken(TokenChar)
}

func (scn FileScanner) scanRegister() Token {
	startRegister := scn.skip()
	if startRegister != 'R' && startRegister != 'r' {
		ShowError("Expected register to start with 'R' or 'r'")
	}
	scn.consume()
	return scn.createToken(TokenRegister)
}

func (scn FileScanner) scanTag() Token {
	scn.skipExpected(':', "Expected ':' before tag")
	for !scn.isAtEnd() && scn.current() != '\n' {
		scn.consume()
	}
	return scn.createToken(TokenTag)
}

func (scn FileScanner) createToken(tokenType TokenType) Token {
	return Token{
		TokenType: tokenType,
		Literal:   string(scn.word),
	}
}

func (scn FileScanner) createError(msg string) Token {
	return Token{
		TokenType: TokenError,
		Literal:   msg,
	}
}

func (scn FileScanner) Scan() Token {
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

func (scn FileScanner) Advance(times int) {
	for i := 0; i < times; i++ {
		token := scn.Scan()
		if !token.IsCorrect() {
			break
		}
	}
}
