package tisasm

import (
	"fmt"
	"strings"
)

type Token struct {
	TokenType TokenType
	Literal   string
	Line      int
}

type TokenType string

const (
	TokenTag         TokenType = "TokenTag"
	TokenRegister    TokenType = "TokenRegister"
	TokenChar        TokenType = "TokenChar"
	TokenMemory      TokenType = "TokenMemory"
	TokenSection     TokenType = "TokenSection"
	TokenNumber      TokenType = "TokenNumber"
	TokenHex         TokenType = "TokenHex"
	TokenString      TokenType = "TokenString"
	TokenInstruction TokenType = "TokenInstruction"
	TokenError       TokenType = "TokenError"
	TokenEof         TokenType = "TokenEof"
)

func (token Token) IsCorrect() bool {
	return token.TokenType != TokenEof && token.TokenType != TokenError
}

func (token Token) IsSection(literal string) bool {
	return token.TokenType == TokenSection && token.Literal == literal
}

func (token Token) IsType(tokenType TokenType) bool {
	return token.TokenType == tokenType
}

func (token Token) IsAnyTypeOf(tokenType ...TokenType) bool {
	for _, t := range tokenType {
		if token.IsType(t) {
			return true
		}
	}
	return false
}

func (token Token) AsInstruction() Instruction {
	if !token.IsType(TokenInstruction) {
		panic("Getting instruction size of something that is not an instruction")
	}
	lowerCase := strings.ToLower(token.Literal)
	instruction, err := GetInstruction(lowerCase)
	if err != nil {
		ShowErrorTokenf(token, "Getting instruction size of undefined instruction: '%s',", token.Literal)
	}
	return instruction
}

func (token Token) String() string {
	return fmt.Sprintf("[%s] '%s' at line %d\n", token.TokenType, token.Literal, token.Line)
}
