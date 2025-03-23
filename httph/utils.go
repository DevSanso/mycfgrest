package httph

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"strings"

	"mycfgrest/types"
)

func parsingBody(body []byte, m map[string]any) error {
	if err := json.Unmarshal(body, m); err != nil {
		return types.NewAppError(err, "parsing body is error")
	}

	return nil
}

func getRequestUrlData(h *HttpHandle, r *http.Request, pVal *types.ParsingMap) error {
	qUrl := r.URL.Query()

	for k, val := range h.meta.Data.Request.QueryString {
		data := qUrl.Get(k)

		if data != "" {
			pVal.Set(0, strings.Join([]string{"request","query_string",val.Symbol},"."), data, types.ParsingValueDataType(val.Type))
		}
	}

	return nil
}

func getRequestBodyData(h *HttpHandle, r *http.Request, pVal *types.ParsingMap) error {
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

func CreateResponseFromTemplate(template string, p *types.ParsingMapFetch) (res string, err error) {
	keys,vals, _, err := p.GetData()
	if err != nil {
		return "", types.NewAppError(err, "parsing value is fetch error")
	}

	if keys == nil || vals == nil {
		return "", types.NewAppError(types.ErrorAppSys, "no data")
	}
	
	var buffer bytes.Buffer
	lastIndex := 0

	for i := 0; i < len(template); i++ {
		if template[i] == '#' {
			if i+1 < len(template) && template[i+1] == '#' {
				// '##' 처리
				buffer.WriteString(template[lastIndex:i])
				buffer.WriteByte('#')
				i++
				lastIndex = i + 1
			} else if i+1 < len(template) && template[i+1] == '{' {
				// '#{...}' 처리
				end := strings.IndexByte(template[i:], '}')
				if end != -1 {
					end += i
					key := ""

					key = template[i+2 : end]
					
					if idx := slices.Index(keys, key); idx != -1 {
						buffer.WriteString(template[lastIndex:i])
						buffer.WriteString(fmt.Sprint(vals[idx]))

						i = end
						lastIndex = end + 1
					}
				}
			}
		}
	}
	buffer.WriteString(template[lastIndex:])

	return buffer.String(), nil
}