package util

import (
	"context"
	"os"
	"os/signal"
	"sync/atomic"

	log "github.com/sirupsen/logrus"
)

// ContextWithCancelOnSignals returns a context that is canceled when one of the specified signals
// are received
func ContextWithCancelOnSignals(ctx context.Context, sigs ...os.Signal) context.Context {
	ctx, cancel := context.WithCancel(ctx)
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, sigs...)

	go func() {
		cancels := 0
		for {
			s := <-sig
			cancels++
			if cancels == 1 {
				log.Warnf("got signal: %s", s)
				cancel()
			} else {
				log.Warn("forcing exit")
				os.Exit(1)
			}
		}
	}()

	return ctx
}

// ContextWithCancelOnDone returns a context that depends on the done channel
// of another context.
func ContextWithCancelOnDone(done <-chan struct{}) context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-done
		cancel()
	}()

	return ctx
}

// Must calls log.Fatal in case err is not nil
func Must(err error) {
	if err != nil {
		log.WithError(err).Fatal("fatal error")
	}
}

// AtomicBool a nice wrapper around the atomit package.
type AtomicBool uint32

func (ab *AtomicBool) IsSet() bool {
	return atomic.LoadUint32((*uint32)(ab)) == 1
}

func (ab *AtomicBool) Set() {
	atomic.SwapUint32((*uint32)(ab), 1)
}

func (ab *AtomicBool) Unset() {
	atomic.SwapUint32((*uint32)(ab), 0)
}
