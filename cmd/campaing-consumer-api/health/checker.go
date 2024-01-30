package health

import (
	"context"
	"sync"
	"time"

	"github.com/joomcode/errorx"
)

// Actual healthcheck implemenetation
type HealthMonitor struct {
	Timeout time.Duration
	Checked map[string]HealthCheck
}

// Internal data sent through validation channels
type outcomeResponse struct {
	monitored string
	status    *errorx.Error
}

const defaultTimeout = 5 * time.Second

func (h HealthMonitor) timeoutOrDefault() time.Duration {
	if h.Timeout != 0 {
		return h.Timeout
	}

	return defaultTimeout
}

// Readiness check implementation
func (h HealthMonitor) Check() healthCheckData {
	health := healthCheckData{
		Status:     StatusUp,
		Components: map[string]MonitoredStatus{},
	}

	ctx, cancelFunc := context.WithTimeout(context.Background(), h.timeoutOrDefault())
	defer cancelFunc()

	monitoredResponse := make(chan outcomeResponse, len(h.Checked))
	wg := sync.WaitGroup{}
	wg.Add(len(h.Checked))
	go func() {
		wg.Wait()
		close(monitoredResponse)
	}()

	for key, registered := range h.Checked {
		// output init
		health.Components[key] = StatusUnknown

		// schedule validation check
		go func(ctx context.Context, key string, checker HealthCheck, outcome chan outcomeResponse) {
			outcome <- outcomeResponse{monitored: key, status: checker.Ping(ctx)}
			wg.Done()
		}(ctx, key, registered, monitoredResponse)
	}

	for {
		done := false
		select {
		case <-ctx.Done():
			// some check took too long
			health.Status = StatusDown
			return health

		case outcome, ok := <-monitoredResponse:
			if !ok {
				done = true
				break
			}

			status := StatusUp
			if outcome.status != nil {
				health.Status = StatusDown
				status = StatusDown
			}

			health.Components[outcome.monitored] = status

		}

		if done {
			break
		}
	}

	return health
}
