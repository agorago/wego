package mw

import (
	"context"
	wegocontext "github.com/agorago/wego/context"
	"github.com/agorago/wego/fw"
	"github.com/agorago/wego/util"
)

type ProxyEntrypoint interface {
	Invoke(ctx context.Context, od fw.OperationDescriptor, params ...interface{}) (interface{}, error)
}
type ProxyEntrypointImpl struct{
	PreMiddlewares  []fw.Middleware
	PostMiddlewares []fw.Middleware
}
func MakeProxyEntrypoint() ProxyEntrypoint{
	httpInvoker := HTTPInvoker{}
	return ProxyEntrypointImpl{
		PreMiddlewares: []fw.Middleware{},
		PostMiddlewares: []fw.Middleware{httpInvoker},
	}
}
// ProxyEntryPoint - Provides an entry point for all the proxies.
// A variation of the server side entry point but handles proxy concerns
// This requires client side registration of the service in WeGO
func (pep ProxyEntrypointImpl)Invoke(ctx context.Context, od fw.OperationDescriptor, params ...interface{}) (interface{}, error) {
	// set up the middleware
	ctx = fw.SetOperationDescriptor(ctx, od)
	// construct payload and request headers
	payload := findPayload(od, params...)
	ctx = wegocontext.SetPayload(ctx, payload)
	ctx = constructHeaders(ctx, od, params...)
	chain := pep.constructChain(od)
	return startChain(ctx, chain)
}

func findPayload(od fw.OperationDescriptor, params ...interface{}) interface{} {
	// loop thru the params and return the one that is of type payload
	for index, param := range params {
		if od.Params[index+1].ParamOrigin == fw.PAYLOAD {
			return param
		}
	}
	return nil
}

func constructHeaders(ctx context.Context, od fw.OperationDescriptor, params ...interface{}) context.Context {
	for index, param := range params {
		pd := od.Params[index+1]
		if pd.ParamOrigin != fw.PAYLOAD && pd.ParamOrigin != fw.CONTEXT {
			ctx = wegocontext.Add(ctx, pd.Name, util.ConvertToString(param, pd.ParamKind))
		}
	}
	return ctx
}

func (pep ProxyEntrypointImpl)constructChain(od fw.OperationDescriptor) fw.MiddlewareChain {
	chain := fw.MakeChain()
	chain.Add(pep.PreMiddlewares...)
	if od.ProxyMiddleware != nil {
		chain.Add(od.ProxyMiddleware...)
	}
	chain.Add(pep.PostMiddlewares...)
	return chain
}
