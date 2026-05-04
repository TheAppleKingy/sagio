package transports

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"sagio/internal/infra/configs"
	"sagio/internal/interfaces"
)

type SagioCore struct {
	semaphore   chan struct{}
	wg          sync.WaitGroup
	ctx         context.Context
	cancel      context.CancelFunc
	resource    string
	controllers map[string]interfaces.Controller
}

func NewSagioCore(resources string) *SagioCore {
	//nolint:gosec
	ctx, cancel := context.WithCancel(context.Background())
	return &SagioCore{
		semaphore:   make(chan struct{}),
		ctx:         ctx,
		cancel:      cancel,
		resource:    resources,
		controllers: make(map[string]interfaces.Controller),
	}
}

func (c *SagioCore) StepProcessing(data []byte, key string, transportWg *sync.WaitGroup) error {
	controller, exists := c.controllers[key]
	if !exists {
		return fmt.Errorf("controller for handling '%s' was not registered", key)
	}
	select {
	case c.semaphore <- struct{}{}:
	case <-c.ctx.Done():
		return fmt.Errorf("handling process was cancelled")
	}
	c.wg.Add(1)
	transportWg.Add(1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				slog.Error("Panic occurred when handling message", "type", key)
			}
			transportWg.Done()
			c.wg.Done()
			<-c.semaphore
		}()
		if err := controller(c.ctx, data); err != nil {
			slog.Error("Unable to process saga step step", "step-type", key, "error", err)
		}
	}()
	return nil
}

func (c *SagioCore) Stop() {
	done := make(chan struct{})
	c.cancel()
	go func() {
		c.wg.Wait()
		close(done)
	}()
	select {
	case <-done:
		return
	case <-time.After(time.Second * time.Duration(configs.AppConfig.AppShutdownTimeout)):
		return
	}
}
