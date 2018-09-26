package pgway

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

type B struct {
	Test string
}

type A struct {
	UserId string `pgway_v:"true"`
	Name   string
	A      []string
}

func sampleFunc(a A, request *Request) interface{} {
	println("name:" + a.Name)
	println("request:" + request.RequestData["name"])
	return a.UserId
}
func TestCallFunc_Normal(t *testing.T) {
	userId := "1"

	data := map[string]string{
		"user_id": userId,
		"name":    "paragura",
		"a":       "[{\"Test\" : \"fe\"},{\"Test\" : \"fef\"}]",
	}

	// test
	request := &Request{RequestData: data}
	ret := CallFunc(sampleFunc, request, BindingStrategyCamelCaseToSnakeCase, nil)

	// assert
	assert.Equal(t, userId, ret)
}

func TestCallFunc_FailedValidation(t *testing.T) {

	data := map[string]string{
		"Name": "paragura",
	}

	validationFailedResult := "validation"
	//
	// test
	request := &Request{RequestData: data}
	ret := CallFunc(sampleFunc, request, BindingStrategyKeep, func(a []string) interface{} {
		println("validationFailed:" + strings.Join(a, ","))
		return validationFailedResult
	})
	assert.Equal(t, validationFailedResult, ret)
}
