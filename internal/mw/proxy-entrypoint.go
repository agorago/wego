package mw

import (
	"context"
	bplusc "gitlab.intelligentb.com/devops/bplus/context"
	"gitlab.intelligentb.com/devops/bplus/fw"
	"gitlab.intelligentb.com/devops/bplus/util"
)

// ProxyEntryPoint - Provides an entry point for all the proxies.
// A variation of the server side entry point but handles proxy concerns
// This requires client side registration of the service in BPlus
func ProxyEntrypoint(ctx context.Context, od fw.OperationDescriptor,params ...interface{}) (interface{}, error) {
	// set up the middleware
	ctx = fw.SetOperationDescriptor(ctx,od)
	// construct payload and request headers
	payload := findPayload(od,params...)
	ctx = bplusc.SetPayload(ctx,payload)
	ctx = constructHeaders(ctx,od,params...)
	chain := constructChain(od)
	return startChain(ctx,chain)
}

func findPayload(od fw.OperationDescriptor,params ...interface{})interface{}{
	// loop thru the params and return the one that is of type payload
	for index, param := range params {
		if od.Params[index+1].ParamOrigin == fw.PAYLOAD {
			return param
		}
	}
	return nil
}

func constructHeaders(ctx  context.Context,od fw.OperationDescriptor,params ...interface{}) context.Context{
	for index, param := range params {
		pd := od.Params[index+1]
		if pd.ParamOrigin != fw.PAYLOAD && pd.ParamOrigin != fw.CONTEXT{
			ctx = bplusc.Add(ctx,pd.Name, util.ConvertToString(param, pd.ParamKind))
		}
	}
	return ctx
}

func constructChain(od fw.OperationDescriptor)fw.MiddlewareChain{
	chain := fw.MakeChain()

	if od.ProxyMiddleware != nil {
		for _, mid := range od.ProxyMiddleware {
			chain.Add(mid)
		}
	}
	chain.Add(HTTPInvoker)
	return chain
}