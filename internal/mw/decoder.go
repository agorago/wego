package mw

import (
	"context"
	"encoding/json"
	"io"

	bplusc "github.com/MenaEnergyVentures/bplus/context"
	fw "github.com/MenaEnergyVentures/bplus/fw"
	e "github.com/MenaEnergyVentures/bplus/internal/err"
)

// decodes a readcloser and makes it into a payload using the operation descriptor's payload maker

func decoder(ctx context.Context, chain *fw.MiddlewareChain) context.Context {

	od := fw.GetOperationDescriptor(ctx)
	if od.OpRequestMaker == nil {
		return chain.DoContinue(ctx)
	}
	var request, _ = od.OpRequestMaker(ctx)
	r := bplusc.GetPayload(ctx).(io.ReadCloser)
	if err := json.NewDecoder(r).Decode(&request); err != nil {
		return bplusc.SetError(ctx, e.MakeBplusError(e.DecodingError, err.Error()))
	}

	ctx = bplusc.SetPayload(ctx, request)
	ctx = chain.DoContinue(ctx)
	return ctx
}
