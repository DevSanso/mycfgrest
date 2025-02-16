package types

import (
	"context"
	"net/http"
)

type AppHttpContext struct {
	appError error

	origin context.Context

	r http.Request
}