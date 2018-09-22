package pgway

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestPgwayRequest_InitRequestData_Get(t *testing.T) {
	queryParameters := map[string]string{}
	queryParameters["query"] = "queryData"

	request := Request{
		HTTPMethod:      http.MethodGet,
		QueryParameters: queryParameters,
		Body:            "{ 'body' : 'bodyData' }",
	}

	expectedData := map[string]string{}
	expectedData["query"] = "queryData"

	assert.NoError(t, request.initRequestData())
	assert.Equal(t, expectedData, request.RequestData)
}

func TestPgwayRequest_InitRequestData_POST(t *testing.T) {
	queryParameters := map[string]string{}
	queryParameters["query"] = "queryData"

	request := Request{
		HTTPMethod:      http.MethodPost,
		QueryParameters: queryParameters,
		Body:            "{ \"body\" : \"bodyData\" }",
	}

	expectedData := map[string]string{}
	expectedData["query"] = "queryData"
	expectedData["body"] = "bodyData"

	assert.NoError(t, request.initRequestData())
	assert.Equal(t, expectedData, request.RequestData)
}
