/**
 * @Author raven
 * @Description
 * @Date 2022/6/28
 **/
package binding

import (
	"net/http"
)

type bindingUri struct{}

func (u bindingUri) Name() BindType {
	return BindUri
}

func (u bindingUri) Bind(r *http.Request, model interface{}) error {
	dErr := uriDecoder.Decode(model, r.URL.Query())
	if dErr != nil {
		return dErr
	}
	vErr := valid.Struct(model)
	if vErr != nil {
		return vErr
	}
	return nil
}
