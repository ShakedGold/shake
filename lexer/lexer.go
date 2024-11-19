package lexer

import (
	"bytes"
	"io"
	"regexp"
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
	TokenNewLine
)

var tokenNames = map[TokenType]string{
	TokenUnknown:        "Unknown",
	TokenOperation:      "Operation",
	TokenKeyword:        "Keyword",
	TokenIdentifier:     "Identifier",
	TokenIdentifierType: "Type",
	TokenNumber:         "Number",
	TokenPunctuation:    "Punctuation",
	TokenNewLine:        "NewLine",
}

func (tt TokenType) String() string {
	return tokenNames[tt]
}

type Token struct {
	Type  TokenType
	Value string
}

// Define the keywords
var keywords = map[string]TokenType{
	"if":  TokenKeyword,
	"for": TokenKeyword,
}

func Lex(reader *bytes.Reader) ([]Token, error) {
	var tokens []Token

	// Define the regular expressions for different token types
	identifierRegexp := regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*`) // No colon in identifier regex
	integerRegexp := regexp.MustCompile(`^[0-9]+`)
	operationRegexp := regexp.MustCompile(`^[\+\-\*/=]`)
	punctuationRegexp := regexp.MustCompile(`^[\(\)\{\};,]`)

	wasPreviousColon := false
	for {
		byteResult, err := reader.ReadByte()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		if byteResult == '\n' {
			tokens = append(tokens, Token{Type: TokenNewLine, Value: "NL"})
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
			tokens = append(tokens, Token{Type: TokenNumber, Value: number})
			continue
		}

		// Match identifiers (variable names, function names)
		if identifierRegexp.MatchString(char) {
			// Read identifier
			identifier, err := Scan(reader, identifierRegexp, true)
			if err != nil {
				return nil, err
			}

			// Check if it's a keyword
			if tokenType, ok := keywords[identifier]; ok {
				// It's a keyword
				tokens = append(tokens, Token{Type: tokenType, Value: identifier})
			} else {
				// Regular identifier
				// if wasBeforeColon true, treat this as type
				if wasPreviousColon {
					tokens = append(tokens, Token{Type: TokenIdentifierType, Value: identifier})
					wasPreviousColon = false
				} else {
					tokens = append(tokens, Token{Type: TokenIdentifier, Value: identifier})
				}
			}

			continue
		}

		// Match operations (+, -, *, /, =)
		if operationRegexp.MatchString(char) {
			tokens = append(tokens, Token{Type: TokenOperation, Value: char})
			continue
		}

		// Match punctuation (parentheses, braces, semicolons, etc.)
		if punctuationRegexp.MatchString(char) {
			tokens = append(tokens, Token{Type: TokenPunctuation, Value: char})
			continue
		}

		// Handle colon (possibly for types)
		if char == ":" {
			wasPreviousColon = true
			continue
		}

		// If no match, add an unknown token
		tokens = append(tokens, Token{Type: TokenUnknown, Value: char})
	}

	return tokens, nil
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
