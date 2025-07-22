package core

import (
	"fmt"
	"strconv"
	"strings"
)

type HTTPStatus int

const (
	HTTP_STATUS_OK HTTPStatus = iota
	HTTP_STATUS_UNKNOWN
)

func StringToHTTPStatus(s string) HTTPStatus {
	switch s {
	case "OK":
		return HTTP_STATUS_OK
	default:
		return HTTP_STATUS_UNKNOWN
	}
}

func (status HTTPStatus) ToString() string {
	switch status {
	case HTTP_STATUS_OK:
		return "OK"
	default:
		return "UNKNOWN"
	}
}

type HTTPResponse struct {
	Version    string
	StatusCode int
	Status     HTTPStatus
}

func ParseHTTPResponse(s string) (*HTTPResponse, error) {
	parts := strings.Split(s, "\r\n")
	firstLine := strings.Split(parts[0], " ")

	version := firstLine[0]
	statusCode, err := strconv.Atoi(firstLine[1])
	if err != nil {
		return &HTTPResponse{}, err
	}

	status := firstLine[2]

	return &HTTPResponse{
		Version:    version,
		StatusCode: statusCode,
		Status:     StringToHTTPStatus(status),
	}, nil
}

func (resp *HTTPResponse) ToString() string {
	return fmt.Sprintf("%s %d %s\r\n\r\n", resp.Version, resp.StatusCode, resp.Status.ToString())
}
