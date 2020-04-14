package http

import (
	"context"
	bplusc "gitlab.intelligentb.com/devops/bplus/context"
	fw "gitlab.intelligentb.com/devops/bplus/fw"
	mw "gitlab.intelligentb.com/devops/bplus/internal/mw"
)

// doHttp - Given an operation descriptor invokes the
// operation remotely using HTTP
func doHTTP(ctx context.Context, od fw.OperationDescriptor, params []interface{}) (interface{}, error) {
	// instantiates a list of proxy middlewares
	// the last of which invokes the HTTP request
	ctx = fw.SetProxyOperationDescriptor(ctx, od)
	ctx = bplusc.SetProxyRequestParams(ctx, params)
	chain := fw.MakeChain()
	chain.Add(mw.HTTPInvoker)
	ctx = chain.DoContinue(ctx)
	// fmt.Printf("doHTTP: od is %#v\n", od)

	resp := bplusc.GetProxyResponsePayload(ctx)
	err := bplusc.GetProxyResponseError(ctx)
	return resp, err

}

// ProxyRequest - create a proxy that invokes the service and operation remotely.
func ProxyRequest(ctx context.Context, service string, operation string, params ...interface{}) (interface{}, error) {
	od, err := fw.FindOperationDescriptor(service, operation)
	if err != nil {
		return nil, err
	}
	return doHTTP(ctx, od, params)
}
