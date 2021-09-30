package nap

import (
	"net/url"

	goquery "github.com/google/go-querystring/query"
)

// Path extends the fullURL with the given path by resolving the reference to
// an absolute URL. If parsing errors occur, the fullURL is left unmodified.
func (c *client) buildFullURL(path string) string {
	return c.baseHost + path
}

// buildQueryParamUrl parses url tagged query structs using go-querystring to
// encode them to url.Values and format them onto the url.RawQuery. Any
// query parsing or encoding errors are returned.
func buildQueryParamUrl(fullPath string, queryStruct interface{}) (string, error) {
	reqURL, err := url.Parse(fullPath)
	if err != nil {
		return "", err
	}

	// encodes query structs into a url.Values map and merges maps
	queryValues, err := goquery.Values(queryStruct)
	if err != nil {
		return "", err
	}

	// url.Values format to a sorted "url encoded" string, e.g. "key=val&foo=bar"
	reqURL.RawQuery = queryValues.Encode()

	return reqURL.String(), nil
}
