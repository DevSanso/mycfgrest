package httpd

import (
	"encoding/json"
	"mycfgrest/types"
)

func parsingBody(body []byte, m map[string]any) *types.AppError {
	if err := json.Unmarshal(body, m); err != nil {
		return types.NewAppError(err, "parsing body is error")
	}

	return nil
}
