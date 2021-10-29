package requests

import (
	"context"
	"fmt"
	"sync"

	"github.com/trezorg/plato/pkg/logger"
)

const maxPoolSize = 10

type job func(ctx context.Context) Result

func worker(ctx context.Context, wg *sync.WaitGroup, in <-chan job, out chan<- Result, number int) {

	defer wg.Done()
	res := make(chan Result, 1)
	logger.Infof("Started worker: %d", number)

	go func() {
		res <- (<-in)(ctx)
	}()

	select {
	case <-ctx.Done():
		out <- Result{Error: fmt.Errorf("Cancelled")}
	case result := <-res:
		out <- result
	}

}

type pool struct {
	size uint
	jobs <-chan job
}

func (p *pool) start(ctx context.Context) <-chan Result {
	out := make(chan Result)
	wg := &sync.WaitGroup{}
	go func() {
		defer close(out)
		for i := uint(0); i < p.size; i++ {
			wg.Add(1)
			go worker(ctx, wg, p.jobs, out, int(i+1))
		}
		wg.Wait()
	}()
	return out
}
