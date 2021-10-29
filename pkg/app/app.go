package app

import (
	"context"
	"fmt"

	"github.com/trezorg/plato/pkg/requests"
)

type App struct {
	urls      []string
	requester *requests.Requester
}

func New(urls ...string) (*App, error) {
	if len(urls) == 0 {
		return nil, fmt.Errorf("no urls specified")
	}
	parsed, err := requests.ParseURLs(urls...)
	if err != nil {
		return nil, err
	}
	r, err := requests.Default()
	if err != nil {
		return nil, err
	}
	return &App{urls: parsed, requester: r}, nil
}

func (a *App) Start(ctx context.Context) <-chan requests.Result {
	return a.requester.Process(ctx, a.urls...)
}
