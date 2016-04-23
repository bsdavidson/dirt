// FIXME: handle (RuneError, 0) for empty string and (RuneError, 1) for bad utf-8
// FIXME: parser should parse on a separate goroutine and emit tokens and errors by channel

package parser

import (
	"fmt"
	"unicode/utf8"
)

const eof rune = -1

type Token struct {
	Value string
	Start int
	Width int
}

type parser struct {
	input  string
	pos    int
	start  int
	tokens []Token
}

func (p *parser) appendToken(t Token) {
	p.tokens = append(p.tokens, t)
}

func (p *parser) nextRune() (r rune, w int) {
	if p.pos >= len(p.input) {
		return eof, 0
	}
	return utf8.DecodeRuneInString(p.input[p.pos:])
}

type parseFn func(p *parser) (parseFn, error)

func isQuote(r rune) bool {
	switch r {
	case '"', '\'':
		return true
	default:
		return false
	}
}

func isSpace(r rune) bool {
	switch r {
	case ' ', '\t', '\r', '\n':
		return true
	default:
		return false
	}
}

func parseQuotedString(p *parser) (parseFn, error) {
	t := Token{Start: p.pos}
	quote, w := p.nextRune()
	if !isQuote(quote) {
		return nil, fmt.Errorf("expected quote, got '%c' (%v) at %d", quote, quote, p.pos)
	}
	p.pos += w
	t.Width += w
	if err := readTokenUntilQuote(&t, p, quote); err != nil {
		return nil, err
	}
	p.appendToken(t)
	return parseSpaces, nil
}

func parseSpaces(p *parser) (parseFn, error) {
	for {
		r, w := p.nextRune()
		if r == eof {
			return nil, nil
		} else if !isSpace(r) {
			break
		}
		p.pos += w
		p.start = p.pos
	}
	return parseString, nil
}

func parseString(p *parser) (parseFn, error) {
	t := Token{Start: p.pos}
	for {
		r, w := p.nextRune()
		if r == eof || isSpace(r) {
			break
		} else if p.pos == p.start && isQuote(r) {
			return parseQuotedString, nil
		}
		p.pos += w
		t.Width += w
		t.Value += string(r)
	}
	if t.Width > 0 {
		p.appendToken(t)
	}
	return parseSpaces, nil
}

func readTokenUntilQuote(t *Token, p *parser, quote rune) error {
	escaped := false
	for {
		r, w := p.nextRune()
		p.pos += w
		t.Width += w
		if r == eof {
			return fmt.Errorf("unexpected end of line, expected %c", quote)
		}
		if escaped {
			escaped = false
		} else if r == '\\' {
			escaped = true
			continue
		} else if r == quote {
			p.pos += w
			return nil
		}
		t.Value += string(r)
	}
}

// Parse splits a string into a series of tokens. Tokens are strings of
// characters separated by whitespace. Tokens can be quoted with single or
// double quotes if they include whitespace.
func Parse(input string) ([]Token, error) {
	p := &parser{input: input}
	fn := parseSpaces
	for {
		var err error
		fn, err = fn(p)
		if err != nil {
			return nil, err
		} else if fn == nil {
			break
		}
	}
	return p.tokens, nil
}
