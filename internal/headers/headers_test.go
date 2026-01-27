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
	assert.Equal(t, "localhost:42069", headers.Get("Host"))
	assert.Equal(t, 23, n)
	assert.False(t, done)

	// Test: Invalid spacing header
	headers = NewHeaders()
	data = []byte("       Host : localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

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
	assert.False(t, done)
	assert.Equal(t, "\"application/json\"", headers.Get("Content-Type"))
	total := n
	for {
		n, done, err = headers.Parse(data[total:])
		require.NoError(t, err)
		total += n
		if done {
			break
		}
	}
	assert.Equal(t, "\"application/json\"", headers.Get("Content-Type"))
	assert.Equal(t, "localhost:42069", headers.Get("Host"))
	assert.Equal(t, 73, total)
	assert.True(t, done)

	// Test: Invalid header with whitespace between non-whitespace characters
	headers = NewHeaders()
	data = []byte("Content Type: \"application/json\"\r\n       Host: localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)

	// Test: Invalid header with latex characters
	headers = NewHeaders()
	data = []byte("Co≆tent Type: \"application/json\"\r\n       Host: localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)

	// Test: Valid headers with same name but different values
	headers = NewHeaders()
	data = []byte("Set-Fav: ice cream\r\n  Set-Fav: pork chops\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	assert.False(t, done)
	assert.Equal(t, "ice cream", headers.Get("Set-Fav"))
	total = n
	for {
		n, done, err = headers.Parse(data[total:])
		require.NoError(t, err)
		total += n
		if done {
			break
		}
	}
	assert.True(t, done)
	assert.Equal(t, "ice cream,pork chops", headers.Get("Set-Fav"))

}
