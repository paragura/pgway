package main

import (
	"github.com/paragura/pgway"
	"net/http"
)

type TestParam struct {
	UserId string
	Name   string
	AAA    string `pgway_binding:testParam` // you can bind with the query parameters like this "/test?testParam=1" -> AAA = 1
}

type TestResponse struct {
	Body string
}

func api1(testParam TestParam) interface{} {

	if testParam.UserId == "1" {
		return pgway.ApiException{ErrorCode: pgway.InvalidParameters, Message: "userId 1 is not allowed."}
	}

	response := TestResponse{
		Body: testParam.UserId + "_" + testParam.Name + "_a:" + testParam.AAA,
	}
	return response
}

func main() {
	// lambda.Start(HandlerByAws)

	apis := pgway.Apis{
		pgway.Api{
			Path:       "/api1/:user_id",
			HTTPMethod: http.MethodGet,
			Handler:    api1,
		},
	}

	server := pgway.Server{
		Apis:                  apis,
		BindingNamingStrategy: pgway.BindingStrategyCamelCaseToSnakeCase, // you can bind with the query parameters like this "/test?user_id=1" -> UserId = 1
	}
	server.BootHttpServerWithDefaultConfig()
}
