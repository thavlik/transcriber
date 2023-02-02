package history

import (
	"context"
)

type History interface {
	Push(ctx context.Context, search *Search) error
}
