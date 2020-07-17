package main

type Token struct {
	TokenType TokenType
	Literal   string
}

type TokenType string

const (
	TokenTag         TokenType = "TokenTag"
	TokenRegister    TokenType = "TokenRegister"
	TokenChar        TokenType = "TokenChar"
	TokenMemory      TokenType = "TokenMemory"
	TokenSection     TokenType = "TokenSection"
	TokenNumber      TokenType = "TokenNumber"
	TokenString      TokenType = "TokenString"
	TokenInstruction TokenType = "TokenInstruction"
	TokenError       TokenType = "TokenError"
	TokenEof         TokenType = "TokenEof"
)
