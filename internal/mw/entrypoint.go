package mw

import (
	"context"

	wegoc "github.com/agorago/wego/context"
	fw "github.com/agorago/wego/fw"
)

type Entrypoint interface {
	Invoke(ctx context.Context) (interface{}, error)
}
type EntrypointImpl struct{
	PreMiddlewares  []fw.Middleware
	PostMiddlewares []fw.Middleware
}
func MakeEntrypoint() Entrypoint{
	decoder := Decoder{}
	v10validator := V10Validator{}
	serviceInvoker := ServiceInvoker{}
	return EntrypointImpl{
		PreMiddlewares: []fw.Middleware{decoder,v10validator},
		PostMiddlewares: []fw.Middleware{serviceInvoker},
	}
}

// Invoke - The server-side entry point for invoking any service registered in WeGO
func (ep EntrypointImpl)Invoke(ctx context.Context) (interface{}, error) {
	// set up the middleware
	od := fw.GetOperationDescriptor(ctx)
	chain := fw.MakeChain()
	chain.Add(ep.PreMiddlewares...)
	if od.OpMiddleware != nil {
		chain.Add(od.OpMiddleware...)
	}
	chain.Add(ep.PostMiddlewares...)
	return startChain(ctx, chain)
}

func startChain(ctx context.Context, chain fw.MiddlewareChain) (interface{}, error) {
	// invoke it
	ctx = chain.DoContinue(ctx)
	// process responses
	response := wegoc.GetResponsePayload(ctx)

	err := wegoc.GetError(ctx)
	if err != nil {
		return response, err.(error)
	}
	return response, nil
	return response, nil
}
