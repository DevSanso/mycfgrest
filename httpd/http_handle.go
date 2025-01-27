package httpd

import (
	"net/http"
	"strconv"
	"golang.org/x/exp/slog"

	"mycfgrest/app_error"
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

	value := NewParsingValue()

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

func (h *httpHandle) preCheck(r *http.Request) *app_error.AppError {
	t := r.Header.Get("Content-Type")

	if t != "encoding/json" {
		return app_error.NewError(app_error.ErrorHttpBadRequest, "content-type not support %s", t)
	}

	return nil
}

func (h *httpHandle) getRequestUrlData(r *http.Request, pVal *ParsingValue) *app_error.AppError {
	qUrl := r.URL.Query()

	for k, val := range h.meta.Data.Request.QueryString {
		data := qUrl.Get(k)
		
		if data != "" {
			pVal.Push(val.Symbol, data, ParsingValueDataType(val.Type))
		}
	}

	return nil
}

func (h *httpHandle) getRequestBodyData(r *http.Request, pVal *ParsingValue) *app_error.AppError {
	cLen := r.Header.Get("Content-Length")
	if cLen == "" {
		return app_error.NewError(app_error.ErrorHttpBadRequest, "content-length is not setting")
	}

	cLenCast, castErr := strconv.Atoi(cLen)
	if castErr != nil {
		return app_error.NewError(castErr, "content-length cast error")
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
			return app_error.NewError(readErr, "request read failed")
		}
		var parsingErr error = nil

		if cLenCast < preAllocRequestBodySize {
			parsingErr = parsingBody(preAllocBuffer[:], m)
		} else {
			parsingErr = parsingBody(allocBuffer, m)
		}

		if parsingErr != nil {
			return app_error.NewError(parsingErr, "request read failed, parsing")
		}
	}

	for k, val := range h.meta.Data.Request.Body {
		data := m[k]
		
		if data != nil {
			pVal.Push(val.Symbol, data, ParsingValueDataType(val.Type))
		}
	}
	
	return nil
}

func newHttpHandle(meta *handle.HandleMeta) *httpHandle {
	return &httpHandle{
		meta,
	}
}
