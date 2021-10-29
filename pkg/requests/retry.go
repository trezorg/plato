package requests

import (
	"context"
	"math/rand"
	"time"

	"github.com/Rican7/retry"
	"github.com/Rican7/retry/backoff"
	"github.com/Rican7/retry/jitter"
	"github.com/Rican7/retry/strategy"
	"github.com/trezorg/plato/pkg/logger"
)

func minDuration(t1, t2 time.Duration) time.Duration {
	if t2 > t1 {
		return t1
	}
	return t2
}

func backoffWithJitterContextMaxTimeout(ctx context.Context, maxTimeout time.Duration, algorithm backoff.Algorithm, transformation jitter.Transformation) strategy.Strategy {
	return func(attempt uint) bool {
		if attempt > 0 {
			timeout := transformation(minDuration(maxTimeout, algorithm(attempt)))
			logger.Infof("Get timeout to wait for next request: %v", timeout)
			select {
			case <-ctx.Done():
				return false
			case <-time.After(timeout):
				return true
			}
		}
		return true
	}
}

func makeRetry(ctx context.Context, maxTimeout time.Duration, factor time.Duration, maxJitter float64, attempts uint, f retry.Action) error {
	seed := time.Now().UnixNano()
	random := rand.New(rand.NewSource(seed))
	return retry.Retry(
		f,
		strategy.Limit(attempts),
		backoffWithJitterContextMaxTimeout(
			ctx,
			maxTimeout,
			backoff.BinaryExponential(factor),
			jitter.Deviation(random, maxJitter),
		),
	)
}
