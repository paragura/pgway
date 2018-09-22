package pgway

import (
	"github.com/aws/aws-lambda-go/events"
)

func (req *Request) WithAwsAPIGatewayProxyRequest(baseRequest *events.APIGatewayProxyRequest) {
	req.Path = baseRequest.Path
	req.QueryParameters = baseRequest.QueryStringParameters
	req.HTTPMethod = baseRequest.HTTPMethod
	req.Headers = baseRequest.Headers
	req.Body = baseRequest.Body

}

func (resp *Response) CreateAwsAPIGatewayProxyResponse() events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: resp.StatusCode,
		Headers:    resp.Headers,
		Body:       resp.Body,
	}
}

func (server *Server) HandleAPIGateway(apiGatewayRequest events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {

	request := &Request{}
	request.WithAwsAPIGatewayProxyRequest(&apiGatewayRequest)
	response := server.handle(request)
	return response.CreateAwsAPIGatewayProxyResponse()
}
