package http

import (
	"context"
	fw "github.com/agorago/wego/fw"
	"github.com/agorago/wego/internal/mw"
)

// ProxyService - it proxies a request for a service and an operation and returns a result.
// Calls the HTTP service and passes all the parameters to it.
// Converts the return values from the request back with error handling included.
type ProxyService interface{
	ProxyRequest(ctx context.Context, service string, operation string, params ...interface{}) (interface{}, error)
}

func MakeProxyService(registrationService fw.RegistrationService,
	proxyEntrypoint mw.ProxyEntrypoint)ProxyService{
	return proxyServiceImpl{
		registrationService: registrationService,
		proxyEntrypoint: proxyEntrypoint,
	}
}
type proxyServiceImpl struct{
	registrationService fw.RegistrationService
	proxyEntrypoint mw.ProxyEntrypoint
}

// ProxyRequest - create a proxy that invokes the service and operation remotely.
func (p proxyServiceImpl)ProxyRequest(ctx context.Context, service string, operation string, params ...interface{}) (interface{}, error) {
	od, err := p.registrationService.FindOperationDescriptor(service, operation)
	if err != nil {
		return nil, err
	}
	return p.proxyEntrypoint.Invoke(ctx, od, params...)
}
