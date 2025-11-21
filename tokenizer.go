package analysistools

import (
	"fmt"
	"io"
	"bufio"
	"strings"
	"unicode"
	"unicode/utf8"
)



type Token struct {
	Value string
	LineNo int
	WordNo int
}

// Tokenizer breaks a text document down into a list of word tokens
//
// ```
//  tokens, err := Tokenizer("this is a corpus of text\n\tof many words\nover a few lines\r\n")
// ```
func Tokenizer(src string) ([]*Token, error) {
	in := strings.NewReader(src)
	return TokenReader(in)
}

// TokenReader reads a buffer and returns a list of Tokens
func TokenReader(in io.Reader) ([]*Token, error) {
	scanner := bufio.NewScanner(in)
	scanner.Split(scanWordsAndNewLines)
	lineNo := 0
	wordNo := 0
	results := []*Token{}
	for scanner.Scan() {
		token := scanner.Text()
		if token != "\n" {
			results = append(results, &Token{
				Value: strings.TrimSpace(token),
				LineNo: lineNo,
				WordNo: wordNo,
			})
			wordNo++
		} else {
			lineNo++
		}
	}
	if err := scanner.Err(); err != nil {
		return results, fmt.Errorf("scanning error: %s", err)
	}
	return results, nil
}

// scanWordsAndNewLines is a custom split function that returns words and newlines as tokens
func scanWordsAndNewLines(data []byte, atEOF bool) (advance int, token []byte, err error) {
	// Skip leading spaces (except newlines)
	start := 0
	for width := 0; start < len(data); start += width {
		var r rune
		r, width = utf8.DecodeRune(data[start:])
		if r != ' ' && r != '\t' && r != '\r' {
			break
		}
	}

	// If we're at a newline, return it as a token
	if start < len(data) && data[start] == '\n' {
		return start + 1, []byte{'\n'}, nil
	}

	// Scan until word boundary or newline
	for width, i := 0, start; i < len(data); i += width {
		var r rune
		r, width = utf8.DecodeRune(data[i:])
		if unicode.IsSpace(r) {
			if r == '\n' {
				// Return the word before the newline
				if i > start {
					return i, data[start:i], nil
				}
				// Return the newline itself
				return i + width, []byte{'\n'}, nil
			}
			// Return the word before other whitespace
			if i > start {
				return i, data[start:i], nil
			}
		}
	}

	// If we're at EOF and there's remaining data, return it as a word
	if atEOF && len(data) > start {
		return len(data), data[start:], nil
	}

	// Request more data
	return start, nil, nil
}
