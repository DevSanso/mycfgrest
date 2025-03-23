package main

import (
	"fmt"
	"log/slog"
	"net/http"

	"mycfgrest/global"
	"mycfgrest/httph"
	"mycfgrest/types"
)


func httpRootHandleFunc(w http.ResponseWriter, r *http.Request) {
	reqUrl := r.URL.Path

	meta := global.GetHandleMetaRerf(reqUrl)
	if meta == nil {
		slog.Error(types.NewAppError(types.ErrorAppHttpBadRequest,"not exists handle url : %s", reqUrl).Error())
		w.WriteHeader(404)
		return
	}
	handle := httph.NewHttpHandle(meta)
	
	err := handle.HandleRun(w, r)
	if err != nil {
		switch err {
		case types.ErrorAppHttpBadRequest:
			w.WriteHeader(400)
		case types.ErrorAppDuplicate:
			w.WriteHeader(400)
		default:
			w.WriteHeader(500)
		}

		slog.Error(err.Error())
	} else {
		w.WriteHeader(200)
		slog.Debug(fmt.Sprintf("send data, [addr:%s, url:%s, agent:%s] ", r.RemoteAddr, r.RequestURI, r.UserAgent()))
	}
}