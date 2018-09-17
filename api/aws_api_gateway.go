package api

import (
	"github.com/aws/aws-lambda-go/events"
)

func (req *PgwayRequest) WithAwsAPIGatewayProxyRequest(baseRequest *events.APIGatewayProxyRequest) {
	req.Path = baseRequest.Path
	req.QueryParameters = baseRequest.QueryStringParameters
	req.HTTPMethod = baseRequest.HTTPMethod
	req.Headers = baseRequest.Headers
	req.Body = baseRequest.Body

}

func (resp *PgwayResponse) CreateAwsAPIGatewayProxyResponse() events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: resp.StatusCode,
		Headers:    resp.Headers,
		Body:       resp.Body,
	}
}

func (server *PgwayServer) HandleAPIGateway(apiGatewayRequest events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {

	request := &PgwayRequest{}
	request.WithAwsAPIGatewayProxyRequest(&apiGatewayRequest)
	response := server.handle(request)
	return response.CreateAwsAPIGatewayProxyResponse()
}
