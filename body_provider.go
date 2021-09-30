package nap

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"

	goquery "github.com/google/go-querystring/query"
)

// BodyProvider provides Body content for http.Request attachment.
type BodyProvider interface {
	// ContentType returns the Content-Type of the body.
	ContentType() string
	// Body returns the io.Reader body.
	Body() (io.Reader, error)
}

/**
|----------------------------------------------------------------------------
| JSONBodyProvider encodes a JSON tagged struct value as a Body for requests.
| See https://golang.org/pkg/encoding/json/#MarshalIndent for details.
|--------------------------------------------------------------------------*/
type JSONBodyProvider struct {
	Payload interface{}
}

func (p JSONBodyProvider) ContentType() string {
	return jsonContentType
}

func (p JSONBodyProvider) Body() (io.Reader, error) {
	buf := &bytes.Buffer{}
	err := json.NewEncoder(buf).Encode(p.Payload)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

/**
|----------------------------------------------------------------------------
| FormBodyProvider encodes a url tagged struct value as Body for requests.
| See https://godoc.org/github.com/google/go-querystring/query for details.
|--------------------------------------------------------------------------*/
type FormBodyProvider struct {
	Payload interface{}
}

func (p FormBodyProvider) ContentType() string {
	return formContentType
}

func (p FormBodyProvider) Body() (io.Reader, error) {
	values, err := goquery.Values(p.Payload)
	if err != nil {
		return nil, err
	}
	return strings.NewReader(values.Encode()), nil
}
