package market

import (
	"context"
	"portfolio/server/model"
)

type Market interface {
	Register(context.Context, model.Portfolio) (chan model.Stream, error)
}
