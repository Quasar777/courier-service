package worker

import (
	"context"
	"fmt"
	"time"
)

type Worker struct {
	courier  courierRepo
}

func NewWorker(c courierRepo) *Worker {
	return &Worker{
		courier:  c,
	}
}

func (w *Worker) RunDeliveryChecker(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			fmt.Println("stopping ticker")
			return
		case <-ticker.C:
			err := w.ReleaseExpiredDeliveries(ctx)
			if err != nil {
				fmt.Println("failed to release expired deliveries:", err)
				return
			}
		}
	}
}

func (w *Worker) ReleaseExpiredDeliveries(ctx context.Context) error {
    return w.courier.ReleaseCouriers(ctx)
}