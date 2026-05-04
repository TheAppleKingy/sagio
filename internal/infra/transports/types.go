package transports

type ConnectionProtocol string

const (
	HTTP ConnectionProtocol = "http"
	AMQP ConnectionProtocol = "amqp"
)

type NewTransportFunc func(*SagioCore) (Transport, error)

var TransportFuncMap = map[ConnectionProtocol]NewTransportFunc{
	HTTP: NewHTTPTransport,
	AMQP: NewAMQPTransport,
}
