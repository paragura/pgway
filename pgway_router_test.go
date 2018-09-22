package pgway

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

var server = Server{
	Apis: []Api{
		{
			Path:       "/test1/:fefe/test",
			HTTPMethod: http.MethodGet,
			Handler:    func(testStruct testStruct) testResponse { return testResponse{testStruct.UserId} },
		},
		{
			Path:       "/test1/:fefe",
			HTTPMethod: http.MethodGet,
			Handler:    func(testStruct testStruct) testResponse { return testResponse{testStruct.UserId} },
		},
		{
			Path:       "/test1/:fefe/aaaa",
			HTTPMethod: http.MethodGet,
			Handler:    func(testStruct testStruct) testResponse { return testResponse{testStruct.UserId} },
		},
		{
			Path:       "/test1/bb/test",
			HTTPMethod: http.MethodGet,
			Handler:    func(testStruct testStruct) testResponse { return testResponse{testStruct.UserId} },
		},
		{
			Path:       "/test2/:id",
			HTTPMethod: http.MethodGet,
			Handler:    func(testStruct testStruct) testResponse { return testResponse{testStruct.UserId} },
		},
		{
			Path:       "/test2/:id2",
			HTTPMethod: http.MethodPost,
			Handler:    func(testStruct testStruct) testResponse { return testResponse{testStruct.UserId} },
		},
	},
	BindingNamingStrategy: BindingStrategyCamelCaseToSnakeCase,
}

var tree = server.BuildRoutingTree()
var request = &Request{
	Path:       "/test1/aaa/test",
	HTTPMethod: http.MethodGet,
}

func BenchmarkPgwayRouter_trace(b *testing.B) {
	for i := 0; i < b.N; i++ {
		tree.tracePath(request)
	}
}

func TestPgwayRouter_trace(t *testing.T) {
	api, pathVariables := tree.tracePath(request)

	expected := map[string]string{
		"fefe": "aaa",
	}
	assert.Equal(t, server.Apis[0].Path, api.Path)
	assert.Equal(t, expected, pathVariables)
}
