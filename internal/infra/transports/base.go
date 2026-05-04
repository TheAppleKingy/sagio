package transports

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"
)

type BaseTransport struct {
	core            *SagioCore
	Type            ConnectionProtocol
	shutdownTimeout int
	impl            Transport
	ctx             context.Context
	wg              sync.WaitGroup
}

// shutdown performs a graceful shutdown of message handling.
// Waits for all active StepProcessing calls to finish using WaitGroup, up to a timeout.
// Logs a warning if not all tasks complete within the configured shutdown timeout.
func (t *BaseTransport) shutdown() {
	done := make(chan struct{})
	go func() {
		t.wg.Wait()
		close(done)
	}()
	select {
	case <-done:
		slog.Info(fmt.Sprintf("%s transport successfully done all handles", t.Type))
	case <-time.After(time.Second * time.Duration(t.shutdownTimeout)):
		slog.Warn(fmt.Sprintf("%s transport did not finish all handlings in timeout. Disconnecting will be executed forcely", t.Type))
	}
}

func (t *BaseTransport) Start() error {
	if err := t.impl.connect(); err != nil {
		return err
	}
	slog.Info(fmt.Sprintf("%s transport connected", t.Type))
	t.core.wg.Add(1)
	context.AfterFunc(t.ctx, func() {
		defer t.core.wg.Done()
		t.Stop()
	})
	return nil
}

func (t *BaseTransport) Stop() {
	t.shutdown()
	t.impl.disconnect()
	slog.Info(fmt.Sprintf("%s transport disconnected", t.Type))
}
