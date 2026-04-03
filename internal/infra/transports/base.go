package transports

import "context"

type Transport interface {
	Start(ctx context.Context)
}
