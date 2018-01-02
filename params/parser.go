// Package params unifies access to httprouter, POST, PUT, GET, etc., parameters
// and automatically converts parameters from string to the destination type.
package params

import (
	"errors"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
)

// Parser parses httprouter, POST, PUT, GET, etc., parameters.
type Parser struct {
	request      *http.Request
	routerParams httprouter.Params
}

// NewParser returns a new Parser.
func NewParser(request *http.Request, params httprouter.Params) (*Parser, error) {

	if request.Form == nil {
		if err := request.ParseForm(); err != nil {
			return nil, err
		}
	}

	return &Parser{
		request:      request,
		routerParams: params,
	}, nil
}

// Parse takes a pointer to a struct, and for each struct field it tries to find
// a corresponding parameter, converts the parameter from string to the struct
// field’s type and writes it to the struct field. A struct field and parameter
// correspond when the parameter name matches the lowercased struct field name.
func (p *Parser) Parse(dest interface{}) error {
	v := reflect.ValueOf(dest)

	if v.Kind() != reflect.Ptr || reflect.Indirect(v).Kind() != reflect.Struct {
		return errors.New("argument must be a pointer to a struct")
	}

	v = reflect.Indirect(v)
	t := reflect.TypeOf(v.Interface())

	for i, j := 0, v.NumField(); i < j; i++ {
		paramName := strings.ToLower(t.Field(i).Name)
		if name := t.Field(i).Tag.Get("param"); name != "" {
			paramName = name
		}

		paramValues := p.param(paramName)

		if len(paramValues) == 0 {
			continue
		}

		field := v.Field(i)

		switch field.Type().String() {
		case "string":
			field.SetString(paramValues[0])
		case "bool":
			if s := strings.ToLower(paramValues[0]); s == "yes" {
				field.SetBool(true)
				break
			} else if s == "no" {
				field.SetBool(false)
				break
			}

			x, err := strconv.ParseBool(paramValues[0])
			if err != nil {
				return err
			}
			field.SetBool(x)
		case "float32":
			x, err := strconv.ParseFloat(paramValues[0], 32)
			if err != nil {
				return err
			}
			field.SetFloat(x)
		case "float64":
			x, err := strconv.ParseFloat(paramValues[0], 64)
			if err != nil {
				return err
			}
			field.SetFloat(x)
		case "int":
			x, err := strconv.ParseInt(paramValues[0], 10, 0)
			if err != nil {
				return err
			}
			field.SetInt(x)
		case "int8":
			x, err := strconv.ParseInt(paramValues[0], 10, 8)
			if err != nil {
				return err
			}
			field.SetInt(x)
		case "int16":
			x, err := strconv.ParseInt(paramValues[0], 10, 16)
			if err != nil {
				return err
			}
			field.SetInt(x)
		case "int32":
			x, err := strconv.ParseInt(paramValues[0], 10, 32)
			if err != nil {
				return err
			}
			field.SetInt(x)
		case "int64":
			x, err := strconv.ParseInt(paramValues[0], 10, 64)
			if err != nil {
				return err
			}
			field.SetInt(x)
		case "uint":
			x, err := strconv.ParseUint(paramValues[0], 10, 0)
			if err != nil {
				return err
			}
			field.SetUint(x)
		case "uint8":
			x, err := strconv.ParseUint(paramValues[0], 10, 8)
			if err != nil {
				return err
			}
			field.SetUint(x)
		case "uint16":
			x, err := strconv.ParseUint(paramValues[0], 10, 16)
			if err != nil {
				return err
			}
			field.SetUint(x)
		case "uint32":
			x, err := strconv.ParseUint(paramValues[0], 10, 32)
			if err != nil {
				return err
			}
			field.SetUint(x)
		case "uint64":
			x, err := strconv.ParseUint(paramValues[0], 10, 64)
			if err != nil {
				return err
			}
			field.SetUint(x)
		case "[]string":
			field.Set(reflect.ValueOf(paramValues))
		case "[]bool":
			s := make([]bool, 0, len(paramValues))
			for _, value := range paramValues {
				if str := strings.ToLower(value); str == "yes" {
					s = append(s, true)
					continue
				} else if str == "no" {
					s = append(s, false)
					continue
				}

				x, err := strconv.ParseBool(value)
				if err != nil {
					return err
				}
				s = append(s, x)
			}
			field.Set(reflect.ValueOf(s))
		case "[]float32":
			s := make([]float32, 0, len(paramValues))
			for _, value := range paramValues {
				x, err := strconv.ParseFloat(value, 32)
				if err != nil {
					return err
				}
				s = append(s, float32(x))
			}
			field.Set(reflect.ValueOf(s))
		case "[]float64":
			s := make([]float64, 0, len(paramValues))
			for _, value := range paramValues {
				x, err := strconv.ParseFloat(value, 64)
				if err != nil {
					return err
				}
				s = append(s, x)
			}
			field.Set(reflect.ValueOf(s))
		case "[]int":
			s := make([]int, 0, len(paramValues))
			for _, value := range paramValues {
				x, err := strconv.ParseFloat(value, 0)
				if err != nil {
					return err
				}
				s = append(s, int(x))
			}
			field.Set(reflect.ValueOf(s))
		case "[]int8":
			s := make([]int8, 0, len(paramValues))
			for _, value := range paramValues {
				x, err := strconv.ParseFloat(value, 8)
				if err != nil {
					return err
				}
				s = append(s, int8(x))
			}
			field.Set(reflect.ValueOf(s))
		case "[]int16":
			s := make([]int16, 0, len(paramValues))
			for _, value := range paramValues {
				x, err := strconv.ParseFloat(value, 16)
				if err != nil {
					return err
				}
				s = append(s, int16(x))
			}
			field.Set(reflect.ValueOf(s))
		case "[]int32":
			s := make([]int32, 0, len(paramValues))
			for _, value := range paramValues {
				x, err := strconv.ParseFloat(value, 32)
				if err != nil {
					return err
				}
				s = append(s, int32(x))
			}
			field.Set(reflect.ValueOf(s))
		case "[]int64":
			s := make([]int64, 0, len(paramValues))
			for _, value := range paramValues {
				x, err := strconv.ParseFloat(value, 64)
				if err != nil {
					return err
				}
				s = append(s, int64(x))
			}
			field.Set(reflect.ValueOf(s))
		case "[]uint":
			s := make([]uint, 0, len(paramValues))
			for _, value := range paramValues {
				x, err := strconv.ParseFloat(value, 0)
				if err != nil {
					return err
				}
				s = append(s, uint(x))
			}
			field.Set(reflect.ValueOf(s))
		case "[]uint8":
			s := make([]uint8, 0, len(paramValues))
			for _, value := range paramValues {
				x, err := strconv.ParseFloat(value, 8)
				if err != nil {
					return err
				}
				s = append(s, uint8(x))
			}
			field.Set(reflect.ValueOf(s))
		case "[]uint16":
			s := make([]uint16, 0, len(paramValues))
			for _, value := range paramValues {
				x, err := strconv.ParseFloat(value, 16)
				if err != nil {
					return err
				}
				s = append(s, uint16(x))
			}
			field.Set(reflect.ValueOf(s))
		case "[]uint32":
			s := make([]uint32, 0, len(paramValues))
			for _, value := range paramValues {
				x, err := strconv.ParseFloat(value, 32)
				if err != nil {
					return err
				}
				s = append(s, uint32(x))
			}
			field.Set(reflect.ValueOf(s))
		case "[]uint64":
			s := make([]uint64, 0, len(paramValues))
			for _, value := range paramValues {
				x, err := strconv.ParseFloat(value, 64)
				if err != nil {
					return err
				}
				s = append(s, uint64(x))
			}
			field.Set(reflect.ValueOf(s))
		default:
			return errors.New("unsupported field type " + field.Type().String())
		}
	}

	return nil
}

// param returns the parameter that matches the provided name. It checks
// httprouter, POST, PUT, GET, etc., parameters for a match.
func (p *Parser) param(name string) []string {
	if len(p.routerParams) > 0 {
		for _, routeParam := range p.routerParams {
			if routeParam.Key == name {
				return []string{routeParam.Value}
			}
		}
	}

	if p.request != nil {
		if values, ok := p.request.Form[name]; ok {
			return values
		}
	}

	return nil
}
