package pgway

import (
	"net/http"
	"strings"
)

type Server struct {
	Apis                      Apis
	ValidationFailedProcessor func([]string) interface{} // called when validation failed
	BindingNamingStrategy     BindingNamingStrategy
	Tree                      RouteTree
	Compiled                  bool
}

func (server *Server) BuildRoutingTree() RouteTree {
	tree := RouteTree{Nodes: map[string]*RouteNode{}}
	for _, api := range server.Apis {
		tree.addRoute(api)
	}
	return tree
}

//
// setup server (build routing tree..)
func (server *Server) Compile() {
	server.Tree = server.BuildRoutingTree()
	server.Compiled = true
}

func (server *Server) handle(request *Request) Response {
	if !server.Compiled {
		server.Compile()
	}
	request.Path = UrlSanitize(request.Path)
	apiPtr, pathVariables := server.Tree.tracePath(request)

	if apiPtr == nil {
		//
		// api not found
		return server.createResponse(ApiException{ErrorCode: ApiNotFound})
	}
	request.PathVariables = pathVariables
	if err := request.initRequestData(); err != nil {
		exception := ApiException{Message: "unsupported type post data.:" + request.Body, Error: err, ErrorCode: InternalServerError}
		return server.createResponse(exception)
	}
	//
	// found handler
	resp := server.Exec(apiPtr, request)
	return server.createResponse(resp)

}

func (server *Server) createResponse(baseResponse interface{}) Response {

	response := Response{
		StatusCode: http.StatusOK,
	}

	if exception, ok := baseResponse.(ApiException); ok {
		//
		// error response
		response.StatusCode = exception.ErrorCode.HttpStatus
		body, err := CreateJsonString(exception.ErrorCode)
		if err != nil {
			panic(err)
		}
		response.Body = body
	} else {
		//
		// normal response
		body, err := CreateJsonString(baseResponse)
		if err != nil {
			return server.createResponse(ApiException{Error: err, ErrorCode: InternalServerError})
		}
		response.Body = body
	}
	return response
}

func (server *Server) Exec(api *Api, request *Request) interface{} { //

	validationFailedFunc := server.ValidationFailedProcessor
	if validationFailedFunc == nil {
		validationFailedFunc = DefaultValidationProcessor
	}

	return CallFunc(api.Handler, request, server.BindingNamingStrategy, validationFailedFunc)
}

func DefaultValidationProcessor(failedFields []string) interface{} {
	message := "[validation][failed][Fields]" + strings.Join(failedFields, ",")
	return ApiException{Message: message, ErrorCode: InvalidParameters}
}
