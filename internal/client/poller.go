package client

import (
	"context"
	"fmt"
	"math"
	"time"
)

type PollConfig struct {
	InitialInterval time.Duration
	MaxInterval     time.Duration
	Timeout         time.Duration
}

func DefaultPollConfig() PollConfig {
	return PollConfig{
		InitialInterval: 1 * time.Second,
		MaxInterval:     30 * time.Second,
		Timeout:         10 * time.Minute,
	}
}

func (c *IonQClient) PollJobUntilComplete(ctx context.Context, jobID string, config PollConfig) (*JobResponse, error) {
	deadline := time.Now().Add(config.Timeout)
	attempt := 0

	for {
		if time.Now().After(deadline) {
			return nil, fmt.Errorf("job %s did not complete within %s", jobID, config.Timeout)
		}

		jobResp, err := c.GetJob(ctx, jobID)
		if err != nil {
			return nil, fmt.Errorf("polling job %s: %w", jobID, err)
		}

		switch jobResp.Status {
		case "completed":
			return jobResp, nil
		case "failed":
			return jobResp, fmt.Errorf("job %s failed", jobID)
		case "canceled":
			return jobResp, fmt.Errorf("job %s was canceled", jobID)
		}

		// Exponential backoff: min(initial * 2^attempt, max)
		backoff := time.Duration(float64(config.InitialInterval) * math.Pow(2, float64(attempt)))
		if backoff > config.MaxInterval {
			backoff = config.MaxInterval
		}
		attempt++

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(backoff):
		}
	}
}
