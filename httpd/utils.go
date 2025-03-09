package httpd

import (
	"encoding/json"
	"net/http"
	"strings"
	"strconv"

	"mycfgrest/types"
)

func parsingBody(body []byte, m map[string]any) *types.AppError {
	if err := json.Unmarshal(body, m); err != nil {
		return types.NewAppError(err, "parsing body is error")
	}

	return nil
}

func getRequestUrlData(h *httpHandle, r *http.Request, pVal *types.ParsingMap) *types.AppError {
	qUrl := r.URL.Query()

	for k, val := range h.meta.Data.Request.QueryString {
		data := qUrl.Get(k)

		if data != "" {
			pVal.Set(0, strings.Join([]string{"request","query_string",val.Symbol},"."), data, types.ParsingValueDataType(val.Type))
		}
	}

	return nil
}

func getRequestBodyData(h *httpHandle, r *http.Request, pVal *types.ParsingMap) *types.AppError {
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
			pVal.Set(0, strings.Join([]string{"request","body",val.Symbol},"."), data, types.ParsingValueDataType(val.Type))
		}
	}

	return nil
}