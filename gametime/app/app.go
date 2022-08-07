package app

import (
	"context"
	"errors"
	"gametime"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type App interface {
	Name() string
	Start(context.Context) error
}

func NewRunner(log gametime.Logger) *Runner {
	return &Runner{
		log: log,
	}
}

type Runner struct {
	log gametime.Logger
}

func (r *Runner) Run(ctx context.Context, apps ...App) {
	var wg sync.WaitGroup
	wg.Add(len(apps))

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		select {
		case sig := <-sigs:
			r.log.Info(ctx, "signal intercepted: %v", sig)
			cancel()
		case <-ctx.Done():
		}
	}()

	for _, a := range apps {
		go func(a App) {
			defer wg.Done()
			defer cancel()

			for {
				err := a.Start(ctx)
				var derr gametime.Error
				if errors.As(err, &derr) && derr.Type == gametime.Recoverable {
					r.log.Info(ctx, "restarting %s after recoverable error: %s", a.Name(), err)
					continue
				}
				r.log.Error(ctx, "terminating %s after unrecoverable error: %s", a.Name(), err)
				break
			}
		}(a)
	}
	wg.Wait()
}
