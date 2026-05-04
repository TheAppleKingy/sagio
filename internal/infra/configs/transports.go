package configs

import "fmt"

type HTTPTransportConfig struct {
	Port            int `env:"HTTP_PORT" envDefault:"8080"`
	ShutdownTimeout int `env:"HTTP_SHUTDOWN_TIMEOUT" envDefault:"10"`
}

type AMQPTransportConfig struct {
	User            string `env:"RABBITMQ_DEFAULT_USER" envDefault:"guest"`
	Password        string `env:"RABBITMQ_DEFAULT_PASS" envDefault:"guest"`
	Host            string `env:"RABBITMQ_HOST" envDefault:"rabbitmq"`
	Port            int    `env:"RABBITMQ_PORT" envDefault:"5672"`
	PrefetchCount   int    `env:"RABBITMQ_PREFETCH_COUNT,required"`
	ShutdownTimeout int    `env:"RABBITMQ_SHUTDOWN_TIMEOUT" envDefault:"10"`
}

func (c AMQPTransportConfig) ConnURL() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%d", c.User, c.Password, c.Host, c.Port)
}
