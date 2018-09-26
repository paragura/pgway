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

	pathVariables := map[string]string{}
	pathVariables["path"] = "pathData"

	request := Request{
		HTTPMethod:      http.MethodPost,
		QueryParameters: queryParameters,
		PathVariables:   pathVariables,
	}

	expectedData := map[string]string{}
	expectedData["query"] = "queryData"
	expectedData["path"] = "pathData"

	assert.NoError(t, request.initRequestData())
	assert.Equal(t, expectedData, request.RequestData)
}

func TestRequest_BindWithPostJson(t *testing.T) {
	a := &struct {
		UserId int
		Name   string
	}{}

	request := Request{
		HTTPMethod: http.MethodPost,
		Body:       "{ \"UserId\": 1, \"Name\" : \"namae\" }",
	}
	assert.NoError(t, request.BindWithPostJson(a))
	assert.Equal(t, a.UserId, 1)
	assert.Equal(t, a.Name, "namae")

}
