package nap

import (
	"encoding/json"
	"net/http"
)

// ResponseDecoder decodes http responses into struct values.
type ResponseDecoder interface {
	// Decode decodes the response into the value pointed to by out.
	Decode(resp *http.Response, out interface{}) error
}

/**
|----------------------------------------------------------------------------
| jsonDecoder decodes http response JSON into a JSON-tagged struct value.
| + Decode: Caller must provide a non-nil v and close the resp.Body.
|--------------------------------------------------------------------------*/
type jsonDecoder struct {
}

func (d jsonDecoder) Decode(resp *http.Response, out interface{}) error {
	return json.NewDecoder(resp.Body).Decode(out)
}
