package ports

import "context"

type WakeListener interface {
	Listen(ctx context.Context) error
}
