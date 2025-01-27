package httpd

import (
	"encoding/json"

	"mycfgrest/app_error"
)


func parsingBody(body []byte, m map[string]any) (*app_error.AppError) {
	if err := json.Unmarshal(body, m); err != nil {
		return app_error.NewError(err, "parsing body is error")
	}

	return nil
}