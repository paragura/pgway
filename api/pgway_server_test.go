package api

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

type TestStruct struct {
	UserId string
	Name   string
}

type TestResponse struct {
	Body string
}

func TestPgwayServer_Handle(t *testing.T) {

	api1 := PgwayApi{
		Path:       "/test1",
		HTTPMethod: http.MethodGet,
		Handler:    func(testStruct TestStruct) TestResponse { return TestResponse{testStruct.UserId} },
	}

	api2 := PgwayApi{
		Path:       "/test2",
		HTTPMethod: http.MethodGet,
		Handler:    func(testStruct TestStruct) TestResponse { return TestResponse{testStruct.UserId} },
	}

	api3 := PgwayApi{
		Path:       "/test2",
		HTTPMethod: http.MethodPost,
		Handler: func(testStruct TestStruct) TestResponse {
			return TestResponse{testStruct.UserId + "_" + testStruct.Name}
		},
	}

	queryParameters := map[string]string{}
	queryParameters["user_id"] = "1"

	request := &PgwayRequest{
		Path:            "/test2",
		HTTPMethod:      http.MethodPost,
		QueryParameters: queryParameters,
		Body:            "{\"name\" : \"namae\" }",
	}

	server := PgwayServer{
		Apis:                  []PgwayApi{api1, api2, api3},
		BindingNamingStrategy: BindingStrategyCamelCaseToSnakeCase,
	}

	response := server.handle(request)

	testResponse := &TestResponse{}

	err := json.Unmarshal([]byte(response.Body), testResponse)
	assert.NoError(t, err)
	assert.Equal(t, "1_namae", testResponse.Body)

}
