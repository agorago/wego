package mw

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	bplusc "github.com/agorago/wego/context"
	fw "github.com/agorago/wego/fw"
	e "github.com/agorago/wego/internal/err"
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
		return bplusc.SetError(ctx, e.MakeBplusErrorWithErrorCode(ctx, http.StatusBadRequest, e.DecodingError, map[string]interface{}{
			"Error": err.Error()}))
	}

	ctx = bplusc.SetPayload(ctx, request)
	ctx = chain.DoContinue(ctx)
	return ctx
}
