package main

import (
	"net/http"
	"pgway/api"
	"pgway/model"
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
		return model.ApiException{ErrorCode: model.InvalidParameters, Message: "userId 1 is not allowed."}
	}

	response := TestResponse{
		Body: testParam.UserId + "_" + testParam.Name + "_a:" + testParam.AAA,
	}
	return response
}

func main() {
	// lambda.Start(HandlerByAws)

	apis := api.PgwayApis{
		api.PgwayApi{
			Path:       "/api1/:user_id",
			HTTPMethod: http.MethodGet,
			Handler:    api1,
		},
	}

	server := api.PgwayServer{
		Apis:                  apis,
		BindingNamingStrategy: api.BindingStrategyCamelCaseToSnakeCase, // you can bind with the query parameters like this "/test?user_id=1" -> UserId = 1
	}
	server.BootHttpServerWithDefaultConfig()
}
