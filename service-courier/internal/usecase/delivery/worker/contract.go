package worker

import "context"

type courierRepo interface {
	ReleaseCouriers(ctx context.Context) error
}