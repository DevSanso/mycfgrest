package main

import (
	"net/http"
	
	"mycfgrest/global"
	"mycfgrest/httph"
)


func httpRootHandleFunc(w http.ResponseWriter, r *http.Request) {
	reqUrl := r.URL.Host

	meta := global.GetHandleMetaRerf(reqUrl)
	handle := httph.NewHttpHandle(meta)
	
	handle.HandleRun(w, r)
}