/**
 * @Author raven
 * @Description
 * @Date 2022/6/28
 **/
package binding

import (
	"io/ioutil"
	"net/http"

	"github.com/RavenHuo/go-pkg/encode/json"
)

type bindingJson struct{}

func (u bindingJson) Name() BindType {
	return BindJson
}

func (u bindingJson) Bind(r *http.Request, model interface{}) error {
	bodyBytes := make([]byte, 0)
	if r.Body != nil {
		bodyBytes, _ = ioutil.ReadAll(r.Body)
	}

	dErr := json.NewEncoder().Decode(bodyBytes, &model)
	if dErr != nil {
		return dErr
	}
	vErr := valid.Struct(model)
	if vErr != nil {
		return vErr
	}
	return nil
}
