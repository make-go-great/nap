package nap

import (
	"context"
	"net/http"
)

func (c *client) Get(ctx context.Context, path string, request interface{}, header http.Header, out interface{}) (int, error) {
	// Get Full URL (host/endpoint)
	fullPath := c.buildFullURL(path)

	// Compose host, endpoint and request into query URL (ex: key=val&foo=bar)
	queryURL, err := buildQueryParamUrl(fullPath, request)
	if err != nil {
		return http.StatusBadRequest, err
	}

	// Make request
	response, err := c.httpClient.Get(queryURL, header)
	if err != nil {
		if response != nil {
			return response.StatusCode, err
		}
		return http.StatusBadGateway, err
	}
	defer func() {
		err = response.Body.Close()
		if err != nil {
			c.logger.Errorw("failed to close body", "err", err)
		}
	}()

	// Skip parse response if not need return result
	if out == nil {
		return response.StatusCode, nil
	}

	// Parse response to out struct
	err = c.responseDecoder.Decode(response, out)

	return response.StatusCode, err
}

func (c *client) Post(ctx context.Context, path string, bodyProvider BodyProvider, header http.Header, out interface{}) (int, error) {
	// Get Full URL (host/endpoint)
	fullPath := c.buildFullURL(path)

	// Get body
	body, err := bodyProvider.Body()
	if err != nil {
		return http.StatusBadRequest, err
	}

	header.Add(headerContentTypeKey, bodyProvider.ContentType())

	// Make request
	response, err := c.httpClient.Post(fullPath, body, header)
	if err != nil {
		if response != nil {
			return response.StatusCode, err
		}
		return http.StatusBadGateway, err
	}
	defer func() {
		err = response.Body.Close()
		if err != nil {
			c.logger.Errorw("failed to close body", "err", err)
		}
	}()

	// Skip parse response if not need return result
	if out == nil {
		return response.StatusCode, nil
	}

	// Parse response to out struct
	err = c.responseDecoder.Decode(response, out)

	return response.StatusCode, err
}
