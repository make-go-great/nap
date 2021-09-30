package nap

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

type TestHTTPClientSuite struct {
	suite.Suite
	logger *zap.SugaredLogger
	client Client
}

func (suite *TestHTTPClientSuite) SetupTest() {
	logger, _ := zap.NewProduction()
	suite.logger = logger.Sugar()
}

func TestTestHTTPClientSuite(t *testing.T) {
	suite.Run(t, new(TestHTTPClientSuite))
}

func (suite *TestHTTPClientSuite) Test_Get_RequestSuccess_ExistResponse_ShouldReturnHTTP200() {
	// ARRANGE
	// Mock request
	path := "/orders/status"
	request := mockRequest()
	header := http.Header{}
	out := OutputTest{}

	// Dummy handler http server
	dummyHandler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(suite.T(), http.MethodGet, r.Method)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{ "return_code": 1, "return_message" :  "Giao dịch thành công", "zp_trans_id" : "zp_001_001"}`))
	}

	server := httptest.NewServer(http.HandlerFunc(dummyHandler))
	defer server.Close()
	suite.client = New(server.URL, "", suite.logger)

	// ACTION
	httpCode, err := suite.client.Get(context.Background(), path, request, header, &out)

	// ASSERT
	assert.Equal(suite.T(), http.StatusOK, httpCode)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), mockGETResponse(), out)
}

func (suite *TestHTTPClientSuite) Test_Get_RequestSuccess_NilResponse_ShouldReturnHTTP200() {
	// ARRANGE
	// Mock request
	path := "/orders/status"
	request := mockRequest()
	header := http.Header{}
	out := OutputTest{}

	// Dummy handler http server
	dummyHandler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(suite.T(), http.MethodGet, r.Method)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(``))
	}

	server := httptest.NewServer(http.HandlerFunc(dummyHandler))
	defer server.Close()
	suite.client = New(server.URL, "", suite.logger)

	// ACTION
	httpCode, err := suite.client.Get(context.Background(), path, request, header, nil)

	// ASSERT
	assert.Equal(suite.T(), http.StatusOK, httpCode)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), OutputTest{}, out)
}

func (suite *TestHTTPClientSuite) Test_Get_RequestSuccess_CustomHeader_ShouldReturnHTTP200() {
	// ARRANGE
	// Mock request
	path := "/orders/status"
	request := mockRequest()
	header := http.Header{}
	header.Add("Content-Type", "application/json")
	out := OutputTest{}

	// Dummy handler http server
	dummyHandler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(suite.T(), http.MethodGet, r.Method)
		assert.Equal(suite.T(), "application/json", r.Header.Get("Content-Type"))
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(``))
	}

	server := httptest.NewServer(http.HandlerFunc(dummyHandler))
	defer server.Close()
	suite.client = New(server.URL, "", suite.logger)

	// ACTION
	httpCode, err := suite.client.Get(context.Background(), path, request, header, nil)

	// ASSERT
	assert.Equal(suite.T(), http.StatusOK, httpCode)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), OutputTest{}, out)
}

func (suite *TestHTTPClientSuite) Test_Get_ParseRequestFail_ShouldReturnError() {
	// ARRANGE
	// Mock request
	path := "%zzzzz"
	request := mockRequest()
	header := http.Header{}
	out := OutputTest{}

	// Dummy handler http server
	dummyHandler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(suite.T(), http.MethodGet, r.Method)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(``))
	}

	server := httptest.NewServer(http.HandlerFunc(dummyHandler))
	defer server.Close()
	suite.client = New(server.URL, "", suite.logger)

	// ACTION
	httpCode, err := suite.client.Get(context.Background(), path, request, header, nil)

	// ASSERT
	assert.Equal(suite.T(), http.StatusBadRequest, httpCode)
	assert.Contains(suite.T(), err.Error(), "invalid port")
	assert.Equal(suite.T(), OutputTest{}, out)
}

func (suite *TestHTTPClientSuite) Test_Get_RequestRequestInvalid_ShouldReturnError() {
	// ARRANGE
	// Mock request
	path := "/orders/status"
	request := 5
	header := http.Header{}
	out := OutputTest{}

	// Dummy handler http server
	dummyHandler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(suite.T(), http.MethodGet, r.Method)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(``))
	}

	server := httptest.NewServer(http.HandlerFunc(dummyHandler))
	defer server.Close()
	suite.client = New(server.URL, "", suite.logger)

	// ACTION
	httpCode, err := suite.client.Get(context.Background(), path, request, header, nil)

	// ASSERT
	assert.Equal(suite.T(), http.StatusBadRequest, httpCode)
	assert.EqualError(suite.T(), err, "query: Values() expects struct input. Got int")
	assert.Equal(suite.T(), OutputTest{}, out)
}

func (suite *TestHTTPClientSuite) Test_Get_MakeRequestError_ShouldReturnError() {
	// ARRANGE
	// Mock request
	path := "/orders/status"
	request := mockRequest()
	header := http.Header{}
	out := OutputTest{}

	// Dummy handler http server
	suite.client = New("localhost", "", suite.logger)

	// ACTION
	httpCode, err := suite.client.Get(context.Background(), path, request, header, nil)

	// ASSERT
	assert.Equal(suite.T(), http.StatusBadGateway, httpCode)
	assert.Contains(suite.T(), err.Error(), "unsupported protocol scheme")
	assert.Equal(suite.T(), OutputTest{}, out)
}

type OutputTest struct {
	ReturnCode    int32  `json:"return_code"`
	ReturnMessage string `json:"return_message"`
	ZPTransID     string `json:"zp_trans_id"`
}

func mockRequest() interface{} {
	return struct {
		AppID      int32  `url:"app_id" json:"app_id"`
		AppTransID string `url:"app_trans_id" json:"app_trans_id"`
		Mac        string `url:"mac" json:"mac"`
	}{
		AppID:      1,
		AppTransID: "mmf_transid_210822001",
		Mac:        "MFwwDQYJKoZIhvcNAQEBBQADSwAwSAJBAPl9eHJltu48w1P",
	}
}

func mockGETResponse() OutputTest {
	return OutputTest{
		ReturnCode:    1,
		ReturnMessage: "Giao dịch thành công",
		ZPTransID:     "zp_001_001",
	}
}

func (suite *TestHTTPClientSuite) Test_Post_RequestSuccess_JsonBody_ShouldReturnHTTP200() {
	// ARRANGE
	path := "/api/v2/orders/create"
	bodyProvider := mockJSONBody()
	header := http.Header{}
	out := OutputTest{}

	// Dummy handler http server
	dummyHandler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(suite.T(), http.MethodPost, r.Method)
		assert.Equal(suite.T(), jsonContentType, r.Header.Get(headerContentTypeKey))
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{ "return_code": 1, "return_message" :  "Giao dịch thành công", "zp_trans_id" : "zp_001_001"}`))
	}

	server := httptest.NewServer(http.HandlerFunc(dummyHandler))
	defer server.Close()
	suite.client = New(server.URL, "", suite.logger)

	// ACTION
	httpCode, err := suite.client.Post(context.Background(), path, bodyProvider, header, &out)

	// ASSERT
	assert.Equal(suite.T(), http.StatusOK, httpCode)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), mockGETResponse(), out)
}

func (suite *TestHTTPClientSuite) Test_Post_RequestSuccess_FormBody_ShouldReturnHTTP200() {
	// ARRANGE
	path := "/api/v2/orders/create"
	bodyProvider := mockFormURLEncodeBody()
	header := http.Header{}
	out := OutputTest{}

	// Dummy handler http server
	dummyHandler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(suite.T(), http.MethodPost, r.Method)
		assert.Equal(suite.T(), formContentType, r.Header.Get(headerContentTypeKey))
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{ "return_code": 1, "return_message" :  "Giao dịch thành công", "zp_trans_id" : "zp_001_001"}`))
	}

	server := httptest.NewServer(http.HandlerFunc(dummyHandler))
	defer server.Close()
	suite.client = New(server.URL, "", suite.logger)

	// ACTION
	httpCode, err := suite.client.Post(context.Background(), path, bodyProvider, header, &out)

	// ASSERT
	assert.Equal(suite.T(), http.StatusOK, httpCode)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), mockGETResponse(), out)
}

func (suite *TestHTTPClientSuite) Test_Post_RequestSuccess_EmptyResponse_ShouldReturnHTTP200() {
	// ARRANGE
	path := "/api/v2/orders/create"
	bodyProvider := mockJSONBody()
	header := http.Header{}
	out := OutputTest{}

	// Dummy handler http server
	dummyHandler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(suite.T(), http.MethodPost, r.Method)
		assert.Equal(suite.T(), jsonContentType, r.Header.Get(headerContentTypeKey))
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(``))
	}

	server := httptest.NewServer(http.HandlerFunc(dummyHandler))
	defer server.Close()
	suite.client = New(server.URL, "", suite.logger)

	// ACTION
	httpCode, err := suite.client.Post(context.Background(), path, bodyProvider, header, nil)

	// ASSERT
	assert.Equal(suite.T(), http.StatusOK, httpCode)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), OutputTest{}, out)
}

func (suite *TestHTTPClientSuite) Test_Post_RequestFailure_BodyNil_ShouldReturnHTTP400() {
	// ARRANGE
	path := "/api/v2/orders/create"
	bodyProvider := mockFormURLEncodeBody()
	bodyProvider.Payload = "vinhha_test"
	header := http.Header{}
	out := OutputTest{}

	server := httptest.NewServer(http.HandlerFunc(nil))
	defer server.Close()
	suite.client = New(server.URL, "", suite.logger)

	// ACTION
	httpCode, err := suite.client.Post(context.Background(), path, bodyProvider, header, nil)

	// ASSERT
	assert.Equal(suite.T(), http.StatusBadRequest, httpCode)
	assert.Contains(suite.T(), err.Error(), "query: Values() expects struct input. Got string")
	assert.Equal(suite.T(), OutputTest{}, out)
}

func (suite *TestHTTPClientSuite) Test_Post_RequestFailure_InvalidPath_ShouldReturnHTTP502() {
	// ARRANGE
	path := "%zzzzz"
	bodyProvider := mockFormURLEncodeBody()
	header := http.Header{}
	out := OutputTest{}

	server := httptest.NewServer(http.HandlerFunc(nil))
	defer server.Close()
	suite.client = New(server.URL, "", suite.logger)

	// ACTION
	httpCode, err := suite.client.Post(context.Background(), path, bodyProvider, header, nil)

	// ASSERT
	assert.Equal(suite.T(), http.StatusBadGateway, httpCode)
	assert.Contains(suite.T(), err.Error(), "invalid port")
	assert.Equal(suite.T(), OutputTest{}, out)
}

func mockJSONBody() JSONBodyProvider {
	return JSONBodyProvider{
		Payload: mockRequest(),
	}
}

func mockFormURLEncodeBody() FormBodyProvider {
	return FormBodyProvider{
		Payload: mockRequest(),
	}
}

func (suite *TestHTTPClientSuite) TestInitClientWithProxy() {
	// ARRANGE
	path := "/api/v2/orders/create"
	bodyProvider := mockFormURLEncodeBody()
	header := http.Header{}
	out := OutputTest{}

	// Dummy handler http server
	dummyHandler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(suite.T(), http.MethodPost, r.Method)
		assert.Equal(suite.T(), formContentType, r.Header.Get(headerContentTypeKey))
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{ "return_code": 1, "return_message" :  "Giao dịch thành công", "zp_trans_id" : "zp_001_001"}`))
	}

	server := httptest.NewServer(http.HandlerFunc(dummyHandler))
	defer server.Close()
	suite.client = New(server.URL, server.URL, suite.logger)

	// ACTION
	httpCode, err := suite.client.Post(context.Background(), path, bodyProvider, header, &out)

	// ASSERT
	assert.Equal(suite.T(), http.StatusOK, httpCode)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), mockGETResponse(), out)
}
