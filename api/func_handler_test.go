package api

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

type A struct {
	UserId string `pgway_v:"true"`
	Name   string
}

func sampleFunc(a A) interface{} {
	println(a.Name)
	return a.UserId
}
func TestCallFunc_Normal(t *testing.T) {
	userId := "1"

	data := map[string]string{
		"user_id": userId,
		"name":    "paragura",
	}

	// test
	ret := CallFunc(sampleFunc, data, BindingStrategyCamelCaseToSnakeCase, nil)

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
	ret := CallFunc(sampleFunc, data, BindingStrategyKeep, func(a []string) interface{} {
		println("validationFailed:" + strings.Join(a, ","))
		return validationFailedResult
	})
	assert.Equal(t, validationFailedResult, ret)
}
