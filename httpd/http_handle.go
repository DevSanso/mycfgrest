package httpd

import (
	"net/http"
	"strconv"

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

	value := types.NewParsingValue()

	if err := h.getRequestUrlData(r, value); err != nil {
		slog.Error(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))

		return
	}

	if err := h.getRequestBodyData(r, value); err != nil {
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

func (h *httpHandle) getRequestUrlData(r *http.Request, pVal *types.ParsingValue) *types.AppError {
	qUrl := r.URL.Query()

	for k, val := range h.meta.Data.Request.QueryString {
		data := qUrl.Get(k)

		if data != "" {
			pVal.Push(val.Symbol, data, types.ParsingValueDataType(val.Type))
		}
	}

	return nil
}

func (h *httpHandle) getRequestBodyData(r *http.Request, pVal *types.ParsingValue) *types.AppError {
	cLen := r.Header.Get("Content-Length")
	if cLen == "" {
		return types.NewAppError(types.ErrorAppHttpBadRequest, "content-length is not setting")
	}

	cLenCast, castErr := strconv.Atoi(cLen)
	if castErr != nil {
		return types.NewAppError(castErr, "content-length cast error")
	}

	m := make(map[string]any)

	{
		var readErr error = nil
		var preAllocBuffer [preAllocRequestBodySize]byte
		var allocBuffer []byte = nil

		if cLenCast < preAllocRequestBodySize {
			_, readErr = r.Body.Read(preAllocBuffer[:])
		} else {
			allocBuffer = make([]byte, cLenCast)
			_, readErr = r.Body.Read(allocBuffer)
		}

		if readErr != nil {
			return types.NewAppError(readErr, "request read failed")
		}
		var parsingErr error = nil

		if cLenCast < preAllocRequestBodySize {
			parsingErr = parsingBody(preAllocBuffer[:], m)
		} else {
			parsingErr = parsingBody(allocBuffer, m)
		}

		if parsingErr != nil {
			return types.NewAppError(parsingErr, "request read failed, parsing")
		}
	}

	for k, val := range h.meta.Data.Request.Body {
		data := m[k]

		if data != nil {
			pVal.Push(val.Symbol, data, types.ParsingValueDataType(val.Type))
		}
	}

	return nil
}

func newHttpHandle(meta *handle.HandleMeta) *httpHandle {
	return &httpHandle{
		meta,
	}
}
