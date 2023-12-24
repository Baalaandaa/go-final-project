package kafka

import (
	"context"
)

type Adapter interface {
	Consume(ctx context.Context) error
	Shutdown(ctx context.Context)
}
