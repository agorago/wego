package mw

import (
	"context"
	bplusc "github.com/MenaEnergyVentures/bplus/context"
	fw "github.com/MenaEnergyVentures/bplus/fw"
)

// Entrypoint - The entry point for invoking any service registered in BPlus
func Entrypoint(ctx context.Context) (interface{}, error) {
	// set up the middleware

	chain := fw.MakeChain()
	chain.Add(decoder)
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
