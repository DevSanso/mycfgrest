package httph

import (
	"context"
	"net/http"
	"slices"
	"strconv"
	"unsafe"

	"mycfgrest/global"
	"mycfgrest/loader/handle"
	"mycfgrest/types"
)

const (
	preAllocRequestBodySize = 1024 * 4
)

type HttpHandle struct {
	meta *handle.HandleMeta
}

func NewHttpHandle(meta *handle.HandleMeta) HttpHandle {
	return HttpHandle{meta: meta}
}

func (h *HttpHandle) parsingHttpBody(r *http.Request) (*types.ParsingMap, error) {
	ret := types.NewParsingMap()

	if err := getRequestUrlData(h, r, ret); err != nil {
		return nil, types.NewAppError(err, "")
	}
	if err := getRequestBodyData(h, r, ret); err != nil {
		return nil, types.NewAppError(err, "")
	}
	return ret,nil
}

func (h *HttpHandle)loadData(ctx  context.Context, paramAndOutput *types.ParsingMap) error {
	loadMeta := h.meta.Data.Load
	loadLen := len(loadMeta)

	for idx :=0; idx < loadLen; idx += 1 {
		lm := loadMeta[strconv.Itoa(idx)]
		
		pool := global.GetSqlPool(lm.LoadName, ctx)
		query := lm.Command
		if output,err := pool.Run(ctx, query, paramAndOutput); err != nil {
			return types.NewAppError(err, "")
		} else {
			if err = paramAndOutput.OverReadFrom(output); err != nil {
				return types.NewAppError(err, "")
			}
		}
	}

	return nil
}

func (h *HttpHandle)HandleRun(w http.ResponseWriter, r *http.Request) error {
	if err := h.preCheck(r); err != nil {
		return types.NewAppError(err, "")
	}

	var output *types.ParsingMap = nil
	if param, err := h.parsingHttpBody(r); err != nil {
		return types.NewAppError(err, "")
	} else {
		output = param		
	}

	if err := h.loadData(r.Context(), output); err != nil {
		return types.NewAppError(err, "")
	}

	response := ""
	if fetcher, err := output.Fetch(); err != nil {
		return types.NewAppError(err, "")
	} else {
		if res, createResErr := CreateResponseFromTemplate(h.meta.Data.Response.Template, fetcher); createResErr != nil {
			return types.NewAppError(createResErr, "")
		} else {
			response = res
		}
	}

	resHeader := w.Header()
	resHeader.Set("Content-Length", strconv.Itoa(len(response)))
	resHeader.Set("Content-Type", h.meta.Data.Response.ContentType)
	w.Write(unsafe.Slice(unsafe.StringData(response), len(response)))

	return nil
}

// ServeHTTP implements http.Handler.
func (h *HttpHandle) preCheck(r *http.Request) error {
	t := r.Header.Get("Content-Type")
	
	if _, exists := slices.BinarySearch(h.meta.Data.Request.ContentType, t); !exists {
		return types.NewAppError(types.ErrorAppHttpBadRequest, "content-type not support %s", t)
	}

	return nil
}


