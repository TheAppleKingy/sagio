package transports

import (
	"fmt"
	"log/slog"
	"runtime"
	"time"

	"sagio/internal/infra/configs"

	"github.com/caarlos0/env/v11"
	"github.com/google/uuid"
	"github.com/rabbitmq/amqp091-go"
)

// AMQPTransport manages RabbitMQ communication for the saga.
// It handles connection and channel lifecycle and processes messages with concurrency control using WaitGroup.
type AMQPTransport struct {
	*BaseTransport
	conf        *configs.AMQPTransportConfig
	conn        *amqp091.Connection
	channel     *amqp091.Channel
	consumerTag string
}

// NewAMQPTransport creates a new AMQP transport instance.
// Parses configuration from environment variables.
// PrefetchCount defaults to the number of CPUs if not set.
// Returns an error if configuration parsing fails.
func NewAMQPTransport(core *SagioCore) (Transport, error) {
	conf := &configs.AMQPTransportConfig{}
	if err := env.Parse(conf); err != nil {
		return nil, err
	}
	if conf.PrefetchCount <= 0 {
		conf.PrefetchCount = runtime.NumCPU()
	}
	res := &AMQPTransport{
		conf: conf,
	}
	res.BaseTransport = &BaseTransport{
		core:            core,
		Type:            AMQP,
		shutdownTimeout: conf.ShutdownTimeout,
		impl:            res,
		ctx:             core.ctx,
	}
	return res, nil
}

// connect establishes or re-establishes a connection and channel to RabbitMQ.
// The function is idempotent and safe to call multiple times.
// Returns an error if connection or channel creation fails.
func (t *AMQPTransport) connect() error {
	if t.conn == nil || t.conn.IsClosed() {
		conn, err := amqp091.Dial(t.conf.ConnURL())
		if err != nil {
			return err
		}
		t.conn = conn
	}
	if t.channel == nil || t.channel.IsClosed() {
		channel, err := t.conn.Channel()
		if err != nil {
			return err
		}
		t.channel = channel
	}
	t.consumerTag = fmt.Sprintf("sagio-%s", uuid.NewString())
	msgs, err := t.channel.ConsumeWithContext(t.ctx, t.core.resource, t.consumerTag, false, false, false, false, nil)
	if err != nil {
		t.disconnect()
		return err
	}
	go t.handleMessageChan(msgs)
	slog.Info("AMQP transport started consuming messages", "tag", t.consumerTag, "queue", t.core.resource)
	return nil
}

// disconnect gracefully closes the channel and connection to RabbitMQ.
// Errors are intentionally ignored, as this is a cleanup operation.
// Safe to call multiple times.
//
//nolint:errcheck,gosec
func (t *AMQPTransport) disconnect() {
	if t.channel != nil && !t.channel.IsClosed() {
		if err := t.channel.Cancel(t.consumerTag, false); err != nil {
			slog.Error(fmt.Sprintf("Error occurred when AMQP transport tried to cancel consuming: %v", err))
		}
		t.channel.Close()
		t.channel = nil
	}
	if t.conn != nil && !t.conn.IsClosed() {
		t.conn.Close()
		t.conn = nil
	}
}

// handleMessageChan processes messages received from RabbitMQ.
// Checks for context cancellation and triggers shutdown.
// For each message, extracts the message type from Headers and calls core.StepProcessing.
// Sends Nack for errors or missing controller, Ack on success.
func (t *AMQPTransport) handleMessageChan(msgs <-chan amqp091.Delivery) {
	for {
		select {
		case <-t.ctx.Done():
			slog.Error("AMQP transport got cancel command. Unable to continue messages handling")
			return
		case msg, ok := <-msgs:
			if !ok {
				slog.Warn("Message chan was closed", "tag", t.consumerTag, "queue", t.core.resource)
				return
			}
			value, ok := msg.Headers[configs.AppConfig.StepKeyName]
			if !ok {
				slog.Error("Unable to get step type")
				//nolint:errcheck
				msg.Nack(false, false)
				continue
			}
			keyString, ok := value.(string)
			if !ok {
				//nolint:errcheck
				msg.Nack(false, false)
				continue
			}
			if err := t.core.StepProcessing(msg.Body, keyString, &t.wg); err != nil {
				slog.Error(fmt.Sprintf("AMQP transport cannot process message: %v", err))
				//nolint:errcheck
				msg.Nack(false, false)
			} else {
				//nolint:errcheck
				msg.Ack(false)
			}
		default:
			time.Sleep(5 * time.Millisecond)
		}
	}
}
