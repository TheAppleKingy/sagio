package configs

type HTTPTransportConfig struct {
	Port int `env:"HTTP_PORT" envDefault:"8080"`
}

type AMQPTransportConfig struct {
	User     string `env:"AMQP_USER" envDefault:"guest"`
	Password string `env:"AMQP_PASSWORD" envDefault:"guest"`
	Host     string `env:"AMQP_HOST" envDefault:"rabbitmq"`
	Port     int    `env:"AMQP_PORT" envDefault:"5672"`
}
