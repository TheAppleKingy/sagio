package transports

import (
	"fmt"
	"log/slog"
)

type TransportRegistry struct {
	core       *SagioCore
	transports map[ConnectionProtocol]Transport
}

func NewTransportRegistry(core *SagioCore) TransportRegistry {
	return TransportRegistry{
		core:       core,
		transports: make(map[ConnectionProtocol]Transport),
	}
}

func (r *TransportRegistry) Register(protocol ConnectionProtocol) error {
	var transport Transport
	var err error
	if _, exists := r.transports[protocol]; exists {
		return nil
	}
	transportFunc, ok := TransportFuncMap[protocol]
	if !ok {
		return fmt.Errorf("undefined transport type: %s", protocol)
	}
	transport, err = transportFunc(r.core)
	if err != nil {
		return err
	}
	if err = transport.Start(); err != nil {
		for startedType, transport := range r.transports {
			if startedType != protocol {
				transport.Stop()
				slog.Warn(fmt.Sprintf("%s transport stopped", startedType))
			}
		}
		return err
	}
	r.transports[protocol] = transport
	return nil
}

func (r *TransportRegistry) Stop() {
	r.core.Stop()
}
