package steamvdf

import (
	"bytes"
	"fmt"
	"io"
)

// TokenType identifies the kind of a Token produced by the text tokenizer.
type TokenType byte

// Token types of the text format.
const (
	TokenTypeString     TokenType = iota // a quoted string
	TokenTypeChildStart                  // "{"
	TokenTypeChildEnd                    // "}"
)

// Token is one lexical element of a text-format document.
type Token struct {
	Type  TokenType
	Value string // the string contents, with escape sequences kept as written
	Line  int    // 1-based line the token starts on
}

// readText parses a text-format document. Values keep their escape sequences
// exactly as written. Duplicate keys within an object replace the earlier
// value in place, so a document that repeats its root key parses to the last
// occurrence.
func readText(data []byte) (KeyValue, error) {
	tokens, err := tokenize(data)
	if err != nil {
		return KeyValue{}, err
	}

	p := textParser{tokens: tokens}
	var root KeyValue
	if err := p.parseMembers(&root, 0); err != nil {
		return KeyValue{}, err
	}

	switch len(root.Children) {
	case 0:
		return KeyValue{}, nil
	case 1:
		return root.Children[0], nil
	default:
		return KeyValue{}, ErrMultipleRoots
	}
}

// tokenize splits data into tokens. Comments run from "/" or "#" to the end
// of the line, which also skips "#include" and "#base" directives. A trailing
// null byte, which Steam appends to text documents, is ignored.
func tokenize(data []byte) ([]Token, error) {
	var tokens []Token
	line := 1

	for i := 0; i < len(data); i++ {
		switch c := data[i]; c {
		case '"':
			end, ok := scanString(data, i)
			if !ok {
				return nil, fmt.Errorf("vdf: unterminated string starting on line %d: %w", line, io.ErrUnexpectedEOF)
			}
			tokens = append(tokens, Token{Type: TokenTypeString, Value: string(data[i+1 : end]), Line: line})
			line += bytes.Count(data[i:end], []byte{'\n'})
			i = end
		case '{':
			tokens = append(tokens, Token{Type: TokenTypeChildStart, Line: line})
		case '}':
			tokens = append(tokens, Token{Type: TokenTypeChildEnd, Line: line})
		case '/', '#':
			for i+1 < len(data) && data[i+1] != '\n' {
				i++
			}
		case '\n':
			line++
		case ' ', '\t', '\r', 0:
		default:
			return nil, fmt.Errorf("vdf: unexpected character %q on line %d", c, line)
		}
	}

	return tokens, nil
}

// scanString returns the index of the quote that closes the string opening at
// data[open]. A backslash escapes the byte that follows it, so an escaped quote
// does not close the string.
func scanString(data []byte, open int) (int, bool) {
	for j := open + 1; j < len(data); j++ {
		switch data[j] {
		case '\\':
			j++
		case '"':
			return j, true
		}
	}
	return 0, false
}

// textParser is a cursor over the token stream.
type textParser struct {
	tokens []Token
	pos    int
}

func (p *textParser) next() (Token, bool) {
	if p.pos >= len(p.tokens) {
		return Token{}, false
	}
	t := p.tokens[p.pos]
	p.pos++
	return t, true
}

// parseMembers reads key/value pairs into parent. At depth 0 it consumes the
// whole token stream; deeper it stops at the closing brace of parent.
func (p *textParser) parseMembers(parent *KeyValue, depth int) error {
	for {
		key, ok := p.next()
		if !ok {
			if depth == 0 {
				return nil
			}
			return fmt.Errorf("vdf: missing closing brace for %q: %w", parent.Key, io.ErrUnexpectedEOF)
		}

		switch key.Type {
		case TokenTypeChildEnd:
			if depth == 0 {
				return fmt.Errorf("vdf: unexpected '}' on line %d", key.Line)
			}
			return nil
		case TokenTypeChildStart:
			return fmt.Errorf("vdf: unexpected '{' on line %d", key.Line)
		}

		value, ok := p.next()
		if !ok {
			return fmt.Errorf("vdf: key %q on line %d has no value: %w", key.Value, key.Line, io.ErrUnexpectedEOF)
		}

		child := KeyValue{Key: key.Value}
		switch value.Type {
		case TokenTypeString:
			child.Value = value.Value
		case TokenTypeChildStart:
			if err := p.parseMembers(&child, depth+1); err != nil {
				return err
			}
		case TokenTypeChildEnd:
			return fmt.Errorf("vdf: key %q on line %d has no value", key.Value, key.Line)
		}

		parent.SetChild(child)
	}
}
