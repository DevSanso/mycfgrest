package httpd

import (
	"net/http"

	"mycfgrest/loader/handle"
	"mycfgrest/app_error"
)

type HttpRoot struct {
	mux http.ServeMux
}

func NewHttpRoot(li []*handle.HandleMeta) (HttpRoot, *app_error.AppError){
	ret := HttpRoot{}

	for _, m := range li {
		ret.mux.Handle(m.Data.Url,newHttpHandle(m))
	}

	return ret, nil
}