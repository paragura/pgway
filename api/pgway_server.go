package api

import (
	"net/http"
	"pgway/model"
	"pgway/util"
	"sort"
	"strings"
)

type PgwayServer struct {
	Apis                      PgwayApis
	ValidationFailedProcessor func([]string) interface{} // called when validation failed
	BindingNamingStrategy     PgwayBindingNamingStrategy
}

func (server *PgwayServer) handle(request *PgwayRequest) PgwayResponse {
	sort.Sort(server.Apis)

	if err := request.initRequestData(); err != nil {
		exception := model.ApiException{Message: "unsupported type post data.:" + request.Body, Error: err, ErrorCode: model.InternalServerError}
		return server.createResponse(exception)
	}

	sanitizedPath := util.UrlSanitize(request.Path)

	for _, api := range server.Apis {
		if api.IsSamePath(sanitizedPath) && api.HTTPMethod == request.HTTPMethod {
			//
			// found handler
			resp := server.Exec(&api, request)
			return server.createResponse(resp)
		}
	}

	return server.createResponse(model.ApiException{ErrorCode: model.ApiNotFound})
}

func (server *PgwayServer) createResponse(baseResponse interface{}) PgwayResponse {

	response := PgwayResponse{
		StatusCode: http.StatusOK,
	}

	if exception, ok := baseResponse.(model.ApiException); ok {
		//
		// error response
		response.StatusCode = exception.ErrorCode.HttpStatus
		body, err := util.CreateJsonString(exception.ErrorCode)
		if err != nil {
			panic(err)
		}
		response.Body = body
	} else {
		//
		// normal response
		body, err := util.CreateJsonString(baseResponse)
		if err != nil {
			return server.createResponse(model.ApiException{Error: err, ErrorCode: model.InternalServerError})
		}
		response.Body = body
	}
	return response
}

func (server *PgwayServer) Exec(api *PgwayApi, request *PgwayRequest) interface{} { //

	validationFailedFunc := server.ValidationFailedProcessor
	if validationFailedFunc == nil {
		validationFailedFunc = DefaultValidationProcessor
	}

	return CallFunc(api.Handler, request.RequestData, server.BindingNamingStrategy, validationFailedFunc)
}

func DefaultValidationProcessor(failedFields []string) interface{} {
	message := "[validation][failed][Fields]" + strings.Join(failedFields, ",")
	return model.ApiException{Message: message, ErrorCode: model.InvalidParameters}
}
