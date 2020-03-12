package mw

import (
	"context"

	bplusc "gitlab.intelligentb.com/devops/bplus/context"
	fw "gitlab.intelligentb.com/devops/bplus/fw"
)

// Entrypoint - The entry point for invoking any service registered in BPlus
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
	chain.Add(serviceInvoker)
	// invoke it
	ctx = chain.DoContinue(ctx)

	// process responses
	response := bplusc.GetResponsePayload(ctx)

	err := bplusc.GetError(ctx)
	if err != nil {
		return response, err.(error)
	}
	return response, nil
}
