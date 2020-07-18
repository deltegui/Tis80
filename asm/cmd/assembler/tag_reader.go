package main

type tagReader struct {
	line      int
	codeStart int
	tags      map[string]string
	scn       scanner
	isInCodeSection bool
}

func newTagReader(inner scanner) tagReader {
	tagScanner{
		line:      0,
		codeStart: 0,
		tags:      make(map[string]string),
		scn:       inner,
		isInCodeSection: false,
	}
}

func (reader tagReader) GetTags() map[string]string {
	readTags()
	return reader.tags
}

func (reader tagReader) readTags() {
	token := reader.scn.Scan()
	codeSection := false
	for token.TokenType != TokenEof && token.TokenType != TokenError {
		if token.TokenType == TokenSection && token.Literal == '.code' {
			codeSection = true
		}
		if !codeSection {
			continue
		}
		if token.TokenType == TokenTag {
			reader.defineTag(token)
		} else {
			advanceInstruction
		}
	}
}
