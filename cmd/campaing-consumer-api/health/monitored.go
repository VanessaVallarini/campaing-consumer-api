package health

import (
	"context"

	"github.com/joomcode/errorx"
)

type HealthCheck interface {
	Ping(ctx context.Context) *errorx.Error
}
