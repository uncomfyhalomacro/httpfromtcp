package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHeaderParserReader(t *testing.T) {
	// Test: Valid single header
	headers := NewHeaders()
	data := []byte("Host: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)

	val, ok := headers.Get("Host")
	assert.True(t, ok)
	assert.Equal(t, "localhost:42069", val)
	assert.Equal(t, 23, n)
	assert.False(t, done)

	// Test: Invalid spacing header
	headers = NewHeaders()
	data = []byte("       Host : localhost:42069       \r\n\r\n")
	n, _, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)

	// Test: Valid single header with extra whitespace
	headers = NewHeaders()
	data = []byte("       Host: localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	assert.Equal(t, 37, n)
	assert.False(t, done)

	// Test: Valid 2 headers
	headers = NewHeaders()
	data = []byte("Content-Type: \"application/json\"\r\n       Host: localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	assert.Equal(t, 34, n)

	val, ok = headers.Get("Content-Type")
	assert.True(t, ok)
	assert.Equal(t, "\"application/json\"", val)

	total := n
	for !done {
		n, done, err = headers.Parse(data[total:])
		require.NoError(t, err)
		total += n
	}

	val, ok = headers.Get("Content-Type")
	assert.True(t, ok)
	assert.Equal(t, "\"application/json\"", val)

	val, ok = headers.Get("Host")
	assert.True(t, ok)
	assert.Equal(t, "localhost:42069", val)
	assert.Equal(t, 73, total)
	assert.True(t, done)

	// Test: Invalid headers (Whitespace/Latex)
	invalidData := [][]byte{
		[]byte("Content Type: \"application/json\"\r\n"),
		[]byte("Co≆tent Type: \"application/json\"\r\n"),
	}
	for _, d := range invalidData {
		headers = NewHeaders()
		_, _, err = headers.Parse(d)
		require.Error(t, err)
	}

	// Test: Valid headers with same name but different values
	headers = NewHeaders()
	data = []byte("Set-Fav: ice cream\r\n  Set-Fav: pork chops\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)

	// Assuming Get returns the concatenated string or first value depending on your logic
	val, ok = headers.Get("Set-Fav")
	assert.True(t, ok)
	// Update this expectation based on your specific Get implementation logic
	assert.Contains(t, val, "ice cream")

	total = n
	for !done {
		n, done, err = headers.Parse(data[total:])
		require.NoError(t, err)
		total += n
	}
	assert.True(t, done)

	val, ok = headers.Get("Set-Fav")
	assert.True(t, ok)
	assert.Equal(t, "ice cream,pork chops", val)
}
