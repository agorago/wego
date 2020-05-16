package mw

import (
	"context"

	wegoc "github.com/agorago/wego/context"
	fw "github.com/agorago/wego/fw"
)

// Entrypoint - The server-side entry point for invoking any service registered in WeGO
func Entrypoint(ctx context.Context) (interface{}, error) {
	// set up the middleware
	od := fw.GetOperationDescriptor(ctx)
	chain := fw.MakeChain()
	chain.Add(decoder)
	chain.Add(v10validator)
	if od.OpMiddleware != nil {
		for _, mid := range od.OpMiddleware {
			chain.Add(mid)
		}
	}
	chain.Add(ServiceInvoker)
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
