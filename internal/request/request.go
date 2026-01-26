package request

import (
	"fmt"
	"io"
	"slices"
	"strings"
)

const bufferSize = 8

type RequestState int

const (
	Initialized RequestState = iota
	Done
)

type Request struct {
	RequestLine  RequestLine
	RequestState RequestState
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func parseRequestLine(data []byte) (int, error) {
	for i := 0; i < len(data); i++ {
		if i > 0 && data[i-1] == '\r' && data[i] == '\n' {
			// consume everything through '\n'
			return i + 1, nil
		}
	}
	return 0, nil

}

func getRawLines(reader io.Reader) (int, []string, error) {
	data := make([]byte, 8)
	requestLines := []string{}
	currentLine := ""
	consumedN := 0
	for {
		n, err := reader.Read(data)
		if err == io.EOF {
			break
		}
		start := 0
		for i := 0; i < n; i++ {
			j := i - 1
			if j >= 0 && j <= len(data) {
				if data[j] == '\r' && data[i] == '\n' {
					part := string(data[start:i])
					currentLine += part
					requestLines = append(requestLines, currentLine)
					currentLine = ""
					start = i + 1
				}
			}

		}
		lastPart := string(data[start:n])
		currentLine += lastPart
		consumedN += n
	}

	if len(currentLine) > 0 {
		requestLines = append(requestLines, currentLine)
	}
	return consumedN, requestLines, nil
}

func buildRequestLine(line string) (*RequestLine, error) {
	httpVersions := []string{"HTTP/1.1"}
	parts := strings.Fields(line)
	if len(parts) == 0 {
		err := fmt.Errorf("Empty request line: length %v\n", len(parts))
		return nil, err
	}
	if len(parts) == 2 {
		method := parts[0]
		route := parts[1]
		version := "1.1" // HTTP/1.1
		requestLine := RequestLine{
			HttpVersion:   version,
			Method:        method,
			RequestTarget: route,
		}
		return &requestLine, nil
	}
	if len(parts) == 3 {
		method := parts[0]
		route := parts[1]
		version := parts[2]
		if slices.Contains(httpVersions, version) {
			version = strings.Split(version, "/")[1]
			requestLine := RequestLine{
				HttpVersion:   version,
				Method:        method,
				RequestTarget: route,
			}
			return &requestLine, nil
		}
	}
	err := fmt.Errorf("Not a valid request string: %s\n", line)
	return nil, err
}

func (r *Request) parse(data []byte) (int, error) {
	switch r.RequestState {
	case Initialized:
		n, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}
		if n == 0 {
			return 0, nil
		}
		if n > 0 {
			line := string(data[:n])
			requestLine, err := buildRequestLine(line)
			if err != nil {
				return 0, err
			}
			r.RequestLine = *requestLine
			r.RequestState = Done
			return n, nil
		}
	case Done:
		return 0, fmt.Errorf("trying to read data in a done state")
	default:
		return 0, fmt.Errorf("Unknown state: %v", r.RequestState)
	}
	return 0, nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	buf := make([]byte, bufferSize)
	readToIndex := 0
	var requestLine RequestLine
	request := Request{
		RequestLine:  requestLine,
		RequestState: Initialized,
	}
	for {
		if request.RequestState == Done {
			break
		}
		if request.RequestState == Initialized {
			if readToIndex == len(buf) {
				newBuf := make([]byte, len(buf)*2)
				copy(newBuf, buf[:readToIndex])
				buf = newBuf
			}
			numBytesRead, err := reader.Read(buf[readToIndex:])
			if err == io.EOF {
				break
			}
			if err != nil {
				return nil, err
			}
			readToIndex += numBytesRead
			numBytesParsed, err := request.parse(buf[:readToIndex])
			if err != nil {
				return nil, err
			}
			tbuf := make([]byte, readToIndex)
			copy(tbuf, buf[numBytesParsed:readToIndex])
			buf = tbuf
			readToIndex -= numBytesParsed
		}
	}
	return &request, nil
}
