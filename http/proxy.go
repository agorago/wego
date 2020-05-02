package http

import (
	"context"
	fw "gitlab.intelligentb.com/devops/bplus/fw"
	"gitlab.intelligentb.com/devops/bplus/internal/mw"
)

// ProxyRequest - create a proxy that invokes the service and operation remotely.
func ProxyRequest(ctx context.Context, service string, operation string, params ...interface{}) (interface{}, error) {
	od, err := fw.FindOperationDescriptor(service, operation)
	if err != nil {
		return nil, err
	}
	return mw.ProxyEntrypoint(ctx, od, params...)
}
