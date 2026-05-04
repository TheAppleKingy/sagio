package transports

type Transport interface {
	Start() error
	Stop()
	disconnect()
	connect() error
}
