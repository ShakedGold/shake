package lexer

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"regexp"
	"shake/queue"
	"unicode"
)

type TokenType int

const (
	TokenUnknown TokenType = iota
	TokenOperation
	TokenKeyword
	TokenIdentifier
	TokenIdentifierType
	TokenNumber
	TokenPunctuation
	TokenSemicolon
)

var tokenNames = map[TokenType]string{
	TokenUnknown:        "Unknown",
	TokenOperation:      "Operation",
	TokenKeyword:        "Keyword",
	TokenIdentifier:     "Identifier",
	TokenIdentifierType: "Type",
	TokenNumber:         "Number",
	TokenPunctuation:    "Punctuation",
	TokenSemicolon:      "Semicolon",
}

func (tt TokenType) String() string {
	return tokenNames[tt]
}

// Custom MarshalJSON method for TokenType
func (tt TokenType) MarshalJSON() ([]byte, error) {
	// Return the string representation as JSON
	return json.Marshal(tt.String())
}

type Token struct {
	Type       TokenType
	Value      string
	LineNumber uint64
}

func (t Token) GetBinaryPrecedence() (int, error) {
	if t.Type != TokenOperation {
		return 0, errors.New("Not an operation")
	}
	switch t.Value {
	case "+":
	case "-":
		return 0, nil
	case "*":
	case "/":
		return 1, nil
	}
	return 0, errors.New("Not a supported operation")
}

// Define the keywords
var keywords = map[string]TokenType{
	"if":     TokenKeyword,
	"for":    TokenKeyword,
	"fn":     TokenKeyword,
	"return": TokenKeyword,
}

func Lex(reader *bytes.Reader) (*queue.Queue[Token], error) {
	var tokens []Token

	// Define the regular expressions for different token types
	identifierRegexp := regexp.MustCompile(`^[a-zA-Z_]`) // No colon in identifier regex
	integerRegexp := regexp.MustCompile(`^[0-9]+`)
	operationRegexp := regexp.MustCompile(`^[\+\-\*/=]`)
	punctuationRegexp := regexp.MustCompile(`^[\(\)\{\},]`)

	var lineNumber uint64 = 1
	wasPreviousColon := false
	for {
		byteResult, err := reader.ReadByte()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		checkNL := func(b byte) (bool, error) {
			if b == '\r' {
				b, err := reader.ReadByte()
				if err != nil {
					return false, err
				}
				if b == '\n' {
					lineNumber++
					return true, nil
				}
				err = reader.UnreadByte()
				if err != nil {
					return false, err
				}
			} else if b == '\n' {
				lineNumber++
				return true, nil
			}
			return false, nil
		}

		if byteResult == '/' {
			byteResult, err := reader.ReadByte()
			if err != nil {
				return nil, err
			}
			if byteResult == '/' {
				for {
					byteResult, err := reader.ReadByte()
					if err != nil {
						return nil, err
					}
					shouldContinue, err := checkNL(byteResult)
					if err != nil {
						return nil, err
					}
					if shouldContinue {
						break
					}
				}
				continue
			}
		}

		shouldContinue, err := checkNL(byteResult)
		if err != nil {
			return nil, err
		}
		if shouldContinue {
			continue
		}

		// Match ; (end statement)
		if byteResult == ';' {
			tokens = append(tokens, Token{Type: TokenSemicolon, Value: ";", LineNumber: lineNumber})
			continue
		}

		char := string(byteResult)

		// Skip whitespace
		if unicode.IsSpace(rune(byteResult)) {
			continue
		}

		// Match numbers (integers)
		if integerRegexp.MatchString(char) {
			// Read number
			number, err := Scan(reader, integerRegexp, true)
			if err != nil {
				return nil, err
			}
			tokens = append(tokens, Token{Type: TokenNumber, Value: number, LineNumber: lineNumber})
			continue
		}

		// Match identifiers (variable names, function names)
		if identifierRegexp.MatchString(char) {
			// Read identifier
			allowedAfterFirst, err := regexp.Compile("^[a-zA-Z0-9_]")
			if err != nil {
				return nil, err
			}
			identifier, err := Scan(reader, allowedAfterFirst, true)
			if err != nil {
				return nil, err
			}

			// Check if it's a keyword
			if tokenType, ok := keywords[identifier]; ok {
				// It's a keyword
				tokens = append(tokens, Token{Type: tokenType, Value: identifier, LineNumber: lineNumber})
			} else {
				// Regular identifier
				// if wasBeforeColon true, treat this as type
				if wasPreviousColon {
					tokens = append(tokens, Token{Type: TokenIdentifierType, Value: identifier, LineNumber: lineNumber})
					wasPreviousColon = false
				} else {
					tokens = append(tokens, Token{Type: TokenIdentifier, Value: identifier, LineNumber: lineNumber})
				}
			}

			continue
		}

		// Match operations (+, -, *, /, =)
		if operationRegexp.MatchString(char) {
			tokens = append(tokens, Token{Type: TokenOperation, Value: char, LineNumber: lineNumber})
			continue
		}

		// Match punctuation (parentheses, braces, semicolons, etc.)
		if punctuationRegexp.MatchString(char) {
			tokens = append(tokens, Token{Type: TokenPunctuation, Value: char, LineNumber: lineNumber})
			continue
		}

		// Handle colon (possibly for types)
		if char == ":" {
			wasPreviousColon = true
			continue
		}

		// If no match, add an unknown token
		tokens = append(tokens, Token{Type: TokenUnknown, Value: char, LineNumber: lineNumber})
	}

	return queue.NewQueueFromSlice(tokens), nil
}

// Scan function to read and match allowed characters for a token
func Scan(reader *bytes.Reader, allowed *regexp.Regexp, unread bool) (string, error) {
	if unread {
		err := reader.UnreadByte()
		if err != nil {
			return "", err
		}
	}
	firstByte, err := reader.ReadByte()
	if err != nil {
		return "", err
	}
	value := string(firstByte)
	if !allowed.MatchString(value) {
		err = reader.UnreadByte()
		if err != nil {
			return value, err
		}
		return value, nil
	}

	// Read until the next character is not allowed
	for reader.Size() != 0 {
		nextByte, err := reader.ReadByte()
		if err != nil {
			if err == io.EOF {
				return value, nil
			}
			return value, err
		}

		char := string(nextByte)
		if !allowed.MatchString(char) {
			err = reader.UnreadByte()
			if err != nil {
				return value, err
			}
			break
		}
		value += char
	}

	return value, nil
}
