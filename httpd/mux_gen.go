package httpd

import (
	"net/http"

	"mycfgrest/loader/handle"
	"mycfgrest/types"
)

type HttpRoot struct {
	mux http.ServeMux
}

func NewHttpRoot(li []*handle.HandleMeta) (HttpRoot, *types.AppError) {
	ret := HttpRoot{}

	for _, m := range li {
		ret.mux.Handle(m.Data.Url, newHttpHandle(m))
	}

	return ret, nil
}
