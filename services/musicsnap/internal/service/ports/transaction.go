package ports

import (
	"context"
)

type Transaction interface {
	Start(ctx context.Context) error
	Abort(ctx context.Context) error
	Commit(ctx context.Context) error
	End(ctx context.Context)
}
