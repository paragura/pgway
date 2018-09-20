package api

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

type testStruct struct {
	UserId string
	Name   string
	Body   string
}

type testResponse struct {
	Body string
}

func TestPgwayServer_Handle(t *testing.T) {

	api1 := PgwayApi{
		Path:       "/test1/fe",
		HTTPMethod: http.MethodGet,
		Handler:    func(testStruct testStruct) testResponse { return testResponse{testStruct.UserId} },
	}

	api2 := PgwayApi{
		Path:       "/test2/aa",
		HTTPMethod: http.MethodGet,
		Handler:    func(testStruct testStruct) testResponse { return testResponse{testStruct.UserId} },
	}

	api3 := PgwayApi{
		Path:       "/test2/:body",
		HTTPMethod: http.MethodPost,
		Handler: func(testStruct testStruct) testResponse {
			return testResponse{testStruct.UserId + "_" + testStruct.Name + "_" + testStruct.Body}
		},
	}

	queryParameters := map[string]string{}
	queryParameters["user_id"] = "1"

	request := &PgwayRequest{
		Path:            "/test2/fefe",
		HTTPMethod:      http.MethodPost,
		QueryParameters: queryParameters,
		Body:            "{\"name\" : \"namae\" }",
	}

	server := PgwayServer{
		Apis:                  []PgwayApi{api1, api2, api3},
		BindingNamingStrategy: BindingStrategyCamelCaseToSnakeCase,
	}

	response := server.handle(request)

	testResponse := &testResponse{}

	err := json.Unmarshal([]byte(response.Body), testResponse)
	assert.NoError(t, err)
	assert.Equal(t, "1_namae_fefe", testResponse.Body)

}
