package msshutdown

import (
	"context"
	"errors"
	"log"
	"os/signal"
	"syscall"
	"time"
)

type Callback struct {
	Name  string
	FnCtx func(ctx context.Context) error
}

var (
	h *handler
)

var (
	ErrTimeoutExceeded = errors.New("failed to perform graceful shutdown: timeout exceeded")
	ErrForceShutdown   = errors.New("failed to perform graceful shutdown: force shutdown occurred")
)

func init() {
	setupHandler()
}

func setupHandler() {
	h = newHandler()
	signal.Notify(h.stop, syscall.SIGINT, syscall.SIGTERM)
}

// AddCallback registers a callback for execution before shutdown.
func AddCallback(cb *Callback) {
	h.add(cb)
}

// Wait waits for application shutdown.
func Wait(config *Config) error {
	cfg := config
	if cfg == nil {
		cfg = NewDefaultConfig()
	}

	h.wait()

	done := make(chan bool)
	stop := h.stop
	timer := time.NewTimer(cfg.WaitTimeout).C

	go func() {
		time.Sleep(cfg.Delay)

		for i := len(h.callbacks) - 1; i >= 0; i-- {
			log.Printf("> '%v' shuting down...", h.callbacks[i].Name)

			if err := handleCallback(h.callbacks[i], cfg.CallbackTimeout); err != nil {
				log.Printf("> '%v' shutdown error: %v", h.callbacks[i].Name, err)
			} else {
				log.Printf("> '%v' gracefully stopped\n", h.callbacks[i].Name)
			}
		}

		done <- true
	}()

	select {
	case <-done:
		return nil
	case <-stop:
		return ErrForceShutdown
	case <-timer:
		return ErrTimeoutExceeded
	}
}

// Now sends event to initiate graceful shutdown.
func Now() {
	h.stop <- syscall.SIGINT
}

func handleCallback(callback *Callback, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if callback.FnCtx != nil {
		return callback.FnCtx(ctx)
	}

	select {
	case <-ctx.Done():
		return ErrTimeoutExceeded
	}
}
