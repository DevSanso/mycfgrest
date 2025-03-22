package httph

import (
	"net/http"
	"slices"

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

func (h *HttpHandle)loadData(param *types.ParsingMap) error {
	loadMeta := h.meta.Data.Load
	loadLen := len(loadMeta)

	for idx :=0; idx < loadLen; idx += 1 {
		//lm := loadMeta[strconv.Itoa(idx)]
		
		//_ := global.GetSqlPool(lm.LoadName, context.Background())
		
	}

	return nil
}

func (h *HttpHandle)HandleRun(w http.ResponseWriter, r *http.Request) error {
	if err := h.preCheck(r); err != nil {
		return types.NewAppError(err, "")
	}

	/**if param, err := h.parsingHttpBody(r); err != nil {
		return types.NewAppError(err, "")
	} else {
		
	}*/


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


