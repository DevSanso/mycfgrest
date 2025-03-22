package conn

import (
	"context"
	"mycfgrest/types"
)

type SQLConn interface {
	Run(ctx context.Context, cmd string, param *types.ParsingMap) (*types.ParsingMap, error)
}
