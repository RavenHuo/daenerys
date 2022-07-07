/**
 * @Author raven
 * @Description
 * @Date 2022/6/28
 **/
package binding

import (
	"net/http"

	"github.com/go-playground/validator"
	"github.com/gorilla/schema"
)

type Binding interface {
	Name() BindType
	Bind(*http.Request, interface{}) error
}

type BindType string

const (
	BindJson BindType = "bindJson"
	BindUri  BindType = "bindUri"
)

var (
	uriDecoder = schema.NewDecoder()
	valid      = validator.New()
)

func getBinding(bindType BindType) Binding {
	switch bindType {
	case BindUri:
		return bindingUri{}
	case BindJson:
		return bindingJson{}
	default:
		return bindingJson{}
	}
}

func WithType(r *http.Request, model interface{}, bindType BindType) error {
	binding := getBinding(bindType)
	return binding.Bind(r, model)
}

func Default(r *http.Request, model interface{}) error {
	binding := getBinding(BindJson)
	return binding.Bind(r, model)
}
