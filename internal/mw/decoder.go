package mw

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	bplusc "github.com/MenaEnergyVentures/bplus/context"
	fw "github.com/MenaEnergyVentures/bplus/fw"
)

// decodes a readcloser and makes it into a payload using the operation descriptor's payload maker

func decoder(ctx context.Context, chain *fw.MiddlewareChain) context.Context {

	od := fw.GetOperationDescriptor(ctx)
	if od.OpPayloadMaker == nil {
		return chain.DoContinue(ctx)
	}
	var request = od.OpPayloadMaker()
	r := bplusc.GetPayload(ctx).(io.ReadCloser)
	if err := json.NewDecoder(r).Decode(&request); err != nil {
		return bplusc.SetError(ctx, fmt.Errorf("Error in decoding the request. error = %s", err.Error()))
	}

	ctx = bplusc.SetPayload(ctx, request)
	ctx = chain.DoContinue(ctx)
	return ctx
}
