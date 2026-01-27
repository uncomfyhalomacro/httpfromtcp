package headers

import (
	"errors"
	"fmt"
	"log"
	"slices"
	"strings"
)

type HeaderState int

const (
	ParsingHeaders HeaderState = iota
	Done
)

type Headers map[string]string

func NewHeaders() Headers {
	return Headers{}
}

func indexOfCRLF(data []byte) int {
	for i := 1; i < len(data); i++ {
		if data[i-1] == '\r' && data[i] == '\n' {
			return i - 1
		}
	}
	return -1
}

func buildHeaderFrom(line []byte) (string, string, error) {
	// Trim outer whitespace on the whole line
	s := strings.TrimSpace(string(line))

	colon := strings.IndexByte(s, ':')
	if colon == -1 {
		return "", "", errors.New("invalid header: missing ':'")
	}

	rawKey := s[:colon]
	rawValue := s[colon+1:]

	// Field-name can be surrounded by spaces from the line,
	// but must not have spaces *right before* the colon.
	trimmedKey := strings.TrimSpace(rawKey)
	if trimmedKey == "" {
		return "", "", errors.New("invalid header: empty field-name")
	}

	// If trimming changed the *right* side, that means there was
	// space between the field-name and the colon, which is invalid.
	if strings.HasSuffix(rawKey, " ") || strings.HasSuffix(rawKey, "\t") {
		return "", "", errors.New("invalid header: whitespace before colon")
	}

	whitespace := strings.IndexByte(trimmedKey, ' ')
	if whitespace > -1 {
		return "", "", errors.New("invalid header name: whitespace between non-whitespace characters")
	}

	// Value may have optional whitespace around it
	value := strings.TrimSpace(rawValue)

	trimmedKey = strings.ToLower(trimmedKey)

	if !validAllowedCharsInHeader(trimmedKey) {
		return "", "", fmt.Errorf("invalid header name: header name contains not allowed or invalid characters. got `%v`", trimmedKey)
	}

	return trimmedKey, value, nil
}

func validAllowedCharsInHeader(s string) bool {
	source := []rune(s)

	const allowedChars = "abcdefghijklmnopqrstuvwxyz0123456789!#$%&'*+-.^_`|~"
	runeArray := []rune(allowedChars)

	for _, c := range source {
		if !slices.Contains(runeArray, c) {
			return false
		}
	}
	return true
}

func (h Headers) Get(key string) (value string) {
	value = h[strings.ToLower(key)]
	return
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	log.Printf("Parse called with: %q (%v)", data, []byte(data))
	log.Printf("Parse called with: %q", data)
	if len(data) == 0 {
		return 0, false, nil
	}

	log.Printf("Parse called with: %q (%v)", data, []byte(data))

	i := indexOfCRLF(data)
	log.Printf("indexOfCRLF returned: %d", i)
	if i == -1 {
		return 0, false, nil
	}

	if i == 0 {
		// data starts with CRLF: end of headers
		return 2, true, nil
	}

	line := data[:i]
	key, value, err := buildHeaderFrom(line)
	if err != nil {
		return 0, false, err
	}

	oldValue, exists := h[key]
	if exists {
		newValue := oldValue + "," + value
		h[key] = newValue
	} else {
		h[key] = value
	}

	return i + 2, false, nil
}
