package conn

import (
	"context"
	"mycfgrest/types"
)

type Conn interface {
	Run(ctx context.Context, cmd string, param *types.ParsingMap) (*types.ParsingMap, error)
}