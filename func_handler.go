package pgway

import (
	"encoding/json"
	"errors"
	"github.com/huandu/xstrings"
	r "reflect"
	"strconv"
	"strings"
)

type BindingNamingStrategy int

const (
	BindingStrategyKeep = iota
	BindingStrategyCamelCaseToSnakeCase
	BindingStrategySnakeCaseToCamelCase
)

const validationTagName = "pgway_v"
const bindingTagName = "pgway_binding"

func CallFunc(handler interface{}, request *Request, bindingNamingStrategy BindingNamingStrategy, validationFailedProcessor func([]string) interface{}) interface{} {
	method := r.ValueOf(handler)
	if method.Kind() != r.Func {
		panic("[pgway]definition error. Handler must be a func")
	}

	methodType := method.Type()

	in := make([]r.Value, methodType.NumIn())

	if methodType.NumIn() > 0 {
		for i := 0; i < methodType.NumIn(); i++ {
			p := methodType.In(i)

			if p.Kind() == r.Ptr && p.String() == "*pgway.Request" {
				in[i] = r.ValueOf(request)
			} else if p.Kind() != r.Struct {
				in[i] = r.New(p).Elem()
			} else {
				obj, err := CreateInstance(p, request.RequestData, bindingNamingStrategy)
				if err != nil {
					//
					// error occur
					return ApiException{Message: "[PGWay][binding][failed] + " + err.Error(), ErrorCode: InvalidParameters}
				}

				//
				// validation with tag
				// validation runs when  pgway_v: true written on struct tags
				failedFields := ValidateInstance(obj)
				if len(failedFields) > 0 {
					//
					// failed validation.
					return validationFailedProcessor(failedFields)
				}

				in[i] = r.ValueOf(obj)
			}
		}
	}

	out := method.Call(in)

	switch len(out) {
	case 0:
		return nil
	case 1:
		return out[0].Interface()
	default:
		panic("invalid definition")
	}
}

func CreateInstance(p r.Type, data map[string]string, bindingNamingStrategy BindingNamingStrategy) (interface{}, error) {
	obj := r.New(p).Elem()

	for i := 0; i < p.NumField(); i++ {
		structField := p.Field(i)
		objField := obj.Field(i)
		var keyName string
		if tagValue, ok := structField.Tag.Lookup(bindingTagName); ok {
			keyName = tagValue
		} else {
			switch bindingNamingStrategy {
			case BindingStrategyCamelCaseToSnakeCase:
				keyName = xstrings.ToSnakeCase(structField.Name)
			case BindingStrategySnakeCaseToCamelCase:
				keyName = xstrings.ToCamelCase(structField.Name)
			case BindingStrategyKeep:
				keyName = structField.Name
			default:
				panic("unknown Binding Strategy Type")
			}
			//
			// naming conversion type (camel, snake, keep)
		}

		baseValue := data[keyName]
		i, err := fetchValue(baseValue, objField.Type())
		if err != nil {
			return nil, err
		} else {
			objField.Set(r.ValueOf(i))
		}
	}
	return obj.Interface(), nil
}

func isAllowedKind(kind r.Kind) bool {
	return kind == r.String || kind == r.Int || kind == r.Float64 || kind == r.Bool || kind == r.Struct
}

func fetchValue(baseStr string, t r.Type) (interface{}, error) {

	valueFunc := func(defaultValue interface{}, value interface{}, err error) (interface{}, error) {
		if len(baseStr) == 0 {
			return defaultValue, nil
		} else {
			return value, err
		}
	}

	switch kind := t.Kind(); kind {
	case r.String:
		return baseStr, nil
	case r.Int:
		fixedValue, err := strconv.Atoi(baseStr)
		return valueFunc(0, fixedValue, err)
	case r.Float64:
		fixedValue, err := strconv.ParseFloat(baseStr, 64)
		return valueFunc(0.0, fixedValue, err)
	case r.Bool:
		fixedValue, err := strconv.ParseBool(baseStr)
		return valueFunc(false, fixedValue, err)
	case r.Struct:
		//
		// TODO: post -> json support
		i := r.New(t).Interface()
		err := json.Unmarshal([]byte(baseStr), i)
		if val := r.ValueOf(i); val.Kind() == r.Ptr {
			return val.Elem().Interface(), err
		} else {
			return i, err
		}
		return i, err
	case r.Slice:
		arrayObjType := t.Elem()
		if isAllowedKind(arrayObjType.Kind()) {
			sli := r.MakeSlice(r.SliceOf(arrayObjType), 0, 0)
			/*
							slice := r.MakeSlice(r.SliceOf(arrayObjType), 0,0)
							// Create a pointer to a slice value and set it to the slice
							slicePtr := r.New(t)
							slicePtr.Elem().Set(slice)
							json.Unmarshal([]byte(baseStr),slicePtr.Interface())
							return slicePtr.Elem().Interface(),nil
							:WHY?
				panic: reflect.Set: value of type []interface {} is not assignable to type []pgway.B [recovered]
					panic: reflect.Set: value of type []interface {} is not assignable to type []pgway.B
							slice := r.MakeSlice(r.SliceOf(arrayObjType), 0,0)
							i := slice.Interface()
							err := json.Unmarshal([]byte(baseStr),&i)
							return i,err
			*/
			//_ := json.Unmarshal([]byte(baseStr),&arr)
			values := strings.Split(baseStr, ",")
			for _, value := range values {
				i, err := fetchValue(value, arrayObjType)
				if err != nil {
					return nil, err
				} else {
					sli = r.Append(sli, r.ValueOf(i))
				}
			}
			return sli.Interface(), nil
		} else {
			return nil, errors.New("not supported array type:" + arrayObjType.Kind().String())
		}
	default:
		return nil, errors.New("not supported type:" + kind.String())
	}
}

//
// validation with tags.  ( validate if written pgway_v: true)
// return failed field name
func ValidateInstance(obj interface{}) []string {
	value := r.ValueOf(obj)
	valueType := value.Type()
	var failedFieldNames []string

	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		fieldType := valueType.Field(i)
		if tagValue, ok := fieldType.Tag.Lookup(validationTagName); ok && strings.TrimSpace(tagValue) == "true" {
			//
			// need validation
			if !validateField(field) {
				failedFieldNames = append(failedFieldNames, fieldType.Name)
			}
		}
	}
	return failedFieldNames
}

func validateField(field r.Value) bool {
	switch t := field.Type(); t.Kind() {
	case r.String:
		return len(field.String()) > 0
	default:
		return true
	}
}
