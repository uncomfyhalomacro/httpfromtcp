package headers

import (
	"errors"
	"strings"
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

	return trimmedKey, value, nil
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	if len(data) == 0 {
		return 0, false, nil
	}

	i := indexOfCRLF(data)
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

	h[key] = value
	return i + 2, false, nil
}
