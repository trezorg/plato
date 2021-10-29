package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/trezorg/plato/pkg/app"
	"github.com/trezorg/plato/pkg/logger"
)

const (
	version = "0.0.1"
)


func main() {
	logger.Init(logger.WithLogLevel("DEBUG"))
	urls := prepareCliArgs(version)
	requester, err := app.New(urls...)
	if err != nil {
		logger.Fatal(err)
	}

	ctx, done := context.WithCancel(context.Background())
	defer done()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		for sig := range stop {
			logger.Infof("Got OS signal: %s\n", sig)
			done()
			return
		}
	}()

	for result := range requester.Start(ctx) {
		fmt.Println(result)
	}

}
