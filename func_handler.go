package pgway

import (
	"github.com/huandu/xstrings"
	r "reflect"
	"strings"
)

type PgwayBindingNamingStrategy int

const (
	BindingStrategyKeep = iota
	BindingStrategyCamelCaseToSnakeCase
	BindingStrategySnakeCaseToCamelCase
)

const PgwayValidationTagName = "pgway_v"
const PgwayBindingTagName = "pgway_binding"

func CallFunc(handler interface{}, data map[string]string, bindingNamingStrategy PgwayBindingNamingStrategy, validationFailedProcessor func([]string) interface{}) interface{} {
	method := r.ValueOf(handler)
	if method.Kind() != r.Func {
		panic("[pgway]definition error. Handler must be a func")
	}

	methodType := method.Type()

	in := make([]r.Value, methodType.NumIn())

	if methodType.NumIn() > 0 {
		for i := 0; i < methodType.NumIn(); i++ {
			p := methodType.In(i)

			if p.Kind() != r.Struct {
				in[i] = r.New(p).Elem()
				continue
			}

			obj := CreateInstance(p, data, bindingNamingStrategy)
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

func CreateInstance(p r.Type, data map[string]string, bindingNamingStrategy PgwayBindingNamingStrategy) interface{} {
	obj := r.New(p).Elem()

	for i := 0; i < p.NumField(); i++ {
		structField := p.Field(i)
		objField := obj.Field(i)
		var keyName string
		if tagValue, ok := structField.Tag.Lookup(PgwayBindingTagName); ok {
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

		//
		// TODO: string 以外
		objField.Set(r.ValueOf(data[keyName]))
	}
	return obj.Interface()
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
		if tagValue, ok := fieldType.Tag.Lookup(PgwayValidationTagName); ok && strings.TrimSpace(tagValue) == "true" {
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
	//
	// TODO: string 以外
	return len(field.String()) > 0
}
