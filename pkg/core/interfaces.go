package core

import (
	"context"
)

// ShutDowner represents anything that can be shutdown.
type ShutDowner interface {
	ShutDown(ctx context.Context) error
}
