package request

import (
	"fmt"
	"io"
	"regexp"
	"strings"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	request, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	requestLine, err := ParseRequestLine(string(request))
	if err != nil {
		return nil, err
	}
	return &Request{
		RequestLine: requestLine,
	}, nil
}

func ParseRequestLine(requestLine string) (RequestLine, error) {
	parts := strings.Split(requestLine, "\r\n")
	line := parts[0]
	requestLineParts := strings.Fields(line)
	if len(requestLineParts) != 3 {
		return RequestLine{}, io.ErrUnexpectedEOF
	}

	httpVersion := strings.Split(requestLineParts[2], "/")

	if httpVersion[1] != "1.1" {
		return RequestLine{}, io.ErrUnexpectedEOF
	}

	request := RequestLine{
		Method:        requestLineParts[0],
		RequestTarget: requestLineParts[1],
		HttpVersion:   httpVersion[1],
	}

	match, err := regexp.MatchString(`^[A-Z]+$`, request.Method)
	if err != nil || !match {
		return RequestLine{}, fmt.Errorf("error matching method")
	}

	return request, nil
}
