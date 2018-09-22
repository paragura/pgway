# PGWay: light weight WebFramework like Spring-Boot in Go
PGWay is a web framework written in Go (Golang). 

You can use this for:
* Http Server
* APIGateway
* RealtimeCommunication-WebSocket or WebRTC(now developing...)




- [Support Feature](#support-feature)
- [Quick Start](#quick-start)
- [Http Sample](#http-sample)
- [AWSApiGateway Sample](#aws-apigateway-sample)

個人的にはすごい使いやすいと思います。

## support feature
* High Speed Routing (maybe faster than GIN :))
* Auto binding Model and response
* Path Parameters ( in future,RegExp Validation will be added..)
... and improving now!

## Quick Start

you can install like below:

```
go get -t paragura/pgway
```

## Http Sample
this is the sample code for PGWay as http server

```
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
		pgway.Api{
			Path:       "/api1/:user_id",
			HTTPMethod: http.MethodPost,
			Handler:    api1,
		},
	}

	server := pgway.Server{
		Apis:                  apis,
		BindingNamingStrategy: pgway.BindingStrategyCamelCaseToSnakeCase, // you can bind with the query parameters like this "/test?user_id=1" -> UserId = 1
	}
	server.BootHttpServerWithDefaultConfig()
}


```

check sample/http_sample/main.go

## Aws ApiGateway Sample
this is the sample code for PGWay as api gateway server.
```

type TestParam struct {
	UserId string
	Name   string
	AAA    string `pgway_binding:testParam` // you can bind with the query parameters like this "/test?testParam=1" -> AAA = 1
}

type TestResponse struct {
	Body string
}
//
// auto binding from query parameters with TestParam! 
func api1(testParam TestParam) interface{} {

	if testParam.UserId == "1" {
		return pgway.ApiException{ErrorCode: pgway.InvalidParameters, Message: "userId 1 is not allowed."}
	}

	response := TestResponse{
		Body: testParam.UserId + "_" + testParam.Name + "_a:" + testParam.AAA,
	}
	return response
}

func PgwayHandler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	apis := pgway.Apis{
		pgway.Api{
			Path:       "/api1",
			HTTPMethod: http.MethodGet,
			Handler:    api1,
		},
	}

	server := pgway.Server{
		Apis:                  apis,
		BindingNamingStrategy: pgway.BindingStrategyCamelCaseToSnakeCase, // you can bind with the query parameters like this "/test?user_id=1" -> UserId = 1
	}

	return server.HandleAPIGateway(req), nil
}

func main() {

	lambda.Start(PgwayHandler)
}
```

check sample/api_gateway_sample/main.go

