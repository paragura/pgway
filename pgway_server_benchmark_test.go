package pgway

import (
	"net/http"
	"testing"
)

var serverBenchmark = Server{
	Apis: []Api{
		{
			Path:       "/test1/b/test",
			HTTPMethod: http.MethodGet,
			Handler:    func(testStruct testStruct) testResponse { return testResponse{testStruct.UserId} },
		},
		{
			Path:       "/test1/fefe",
			HTTPMethod: http.MethodGet,
			Handler:    func(testStruct testStruct) testResponse { return testResponse{testStruct.UserId} },
		},
		{
			Path:       "/test1/fefe/aaaa",
			HTTPMethod: http.MethodGet,
			Handler:    func(testStruct testStruct) testResponse { return testResponse{testStruct.UserId} },
		},
		{
			Path:       "/test1/bb/test",
			HTTPMethod: http.MethodGet,
			Handler:    func(testStruct testStruct) testResponse { return testResponse{testStruct.UserId} },
		},
		{
			Path:       "/test2/id",
			HTTPMethod: http.MethodGet,
			Handler:    func(testStruct testStruct) testResponse { return testResponse{testStruct.UserId} },
		},
		{
			Path:       "/test2/id2",
			HTTPMethod: http.MethodPost,
			Handler:    func(testStruct testStruct) testResponse { return testResponse{testStruct.UserId} },
		},
	},
	BindingNamingStrategy: BindingStrategyCamelCaseToSnakeCase,
}

func BenchmarkPgwayServer_Handle(b *testing.B) {

	queryParameters := map[string]string{}
	queryParameters["name"] = "namae"

	request := &Request{
		Path:            "/test1/fefe/aaaa",
		HTTPMethod:      http.MethodGet,
		QueryParameters: queryParameters,
		Body:            "{\"body\" : \"body\" }",
	}

	for i := 0; i < b.N; i++ {
		serverBenchmark.handle(request)
	}

}

func BenchmarkPgwayServer_HandleOld(b *testing.B) {

	queryParameters := map[string]string{}
	queryParameters["name"] = "namae"

	request := &Request{
		Path:            "/test1/fefe/aaaa",
		HTTPMethod:      http.MethodGet,
		QueryParameters: queryParameters,
		Body:            "{\"body\" : \"body\" }",
	}

	for i := 0; i < b.N; i++ {
		serverBenchmark.handleOld(request)
	}

}
