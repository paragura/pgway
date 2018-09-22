package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

var initialized = false
var ginLambda *ginadapter.GinLambda

func HandlerByAws(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	if !initialized {
		// stdout and stderr are sent to AWS CloudWatch Logs
		log.Printf("Gin cold start")
		r := gin.Default()
		r.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "pong",
			})
		})

		ginLambda = ginadapter.New(r)
		initialized = true
	}

	// If no name is provided in the HTTP request body, throw an error
	return ginLambda.Proxy(req)
}

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
		return ApiException{ErrorCode: InvalidParameters, Message: "userId 1 is not allowed."}
	}

	response := TestResponse{
		Body: testParam.UserId + "_" + testParam.Name + "_a:" + testParam.AAA,
	}
	return response
}

func PgwayHandler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	apis := api.PgwayApis{
		api.PgwayApi{
			Path:       "/api1",
			HTTPMethod: http.MethodGet,
			Handler:    api1,
		},
	}

	server := api.PgwayServer{
		Apis:                  apis,
		BindingNamingStrategy: api.BindingStrategyCamelCaseToSnakeCase, // you can bind with the query parameters like this "/test?user_id=1" -> UserId = 1
	}

	return server.HandleAPIGateway(req), nil
}

func main() {
	// lambda.Start(HandlerByAws)

	lambda.Start(PgwayHandler)
}
