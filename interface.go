package nap

import (
	"context"
	"net/http"
	"net/url"

	"github.com/gojek/heimdall/v7/httpclient"
	"github.com/gojek/heimdall/v7/plugins"
	"go.uber.org/zap"
)

//go:generate mockgen --source=client.go -destination=./mocks/client.go -package=mocks
type Client interface {
	Get(ctx context.Context, path string, request interface{}, header http.Header, out interface{}) (int, error)
	Post(ctx context.Context, path string, request BodyProvider, header http.Header, out interface{}) (int, error)
}

type client struct {
	// Logger
	logger *zap.SugaredLogger
	// HTTP Client, wrap by httpClient of heimdall
	httpClient *httpclient.Client
	// Host of service
	baseHost string
	// Response decoder/unmarshal
	// Set default by JsonDecoder. Improve custom own options
	responseDecoder ResponseDecoder
}

func New(baseHost string, proxyURL string, logger *zap.SugaredLogger, opts ...httpclient.Option) Client {
	c := &client{
		logger:          logger,
		baseHost:        baseHost,
		responseDecoder: jsonDecoder{},
		httpClient:      getHTTPClient(proxyURL, opts...),
	}

	return c
}

func getHTTPClient(proxyURL string, opts ...httpclient.Option) *httpclient.Client {
	if proxyURL != "" {
		opts = append(opts, httpclient.WithHTTPClient(&http.Client{
			Transport: getTransport(proxyURL),
		}))
	}

	httpClient := httpclient.NewClient(opts...)
	requestLogger := plugins.NewRequestLogger(nil, nil)
	httpClient.AddPlugin(requestLogger)

	return httpClient
}

func getTransport(proxyURLCfg string) http.RoundTripper {
	proxyURL, _ := url.Parse(proxyURLCfg)
	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	}
	return transport
}
