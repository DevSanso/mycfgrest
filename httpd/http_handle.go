package httpd

import (
	"net/http"

	"golang.org/x/exp/slog"

	"mycfgrest/types"
	"mycfgrest/loader/handle"
)

const (
	preAllocRequestBodySize = 1024 * 4
)

type httpHandle struct {
	meta *handle.HandleMeta
}

// ServeHTTP implements http.Handler.
func (h *httpHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := h.preCheck(r); err != nil {
		slog.Error(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))

		return
	}

	value := types.NewParsingMap()

	if err := getRequestUrlData(h, r, value); err != nil {
		slog.Error(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))

		return
	}

	if err := getRequestBodyData(h, r, value); err != nil {
		slog.Error(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))

		return
	}

}

func (h *httpHandle) preCheck(r *http.Request) *types.AppError {
	t := r.Header.Get("Content-Type")

	if t != "encoding/json" {
		return types.NewAppError(types.ErrorAppHttpBadRequest, "content-type not support %s", t)
	}

	return nil
}



func newHttpHandle(meta *handle.HandleMeta) *httpHandle {
	return &httpHandle{
		meta,
	}
}
