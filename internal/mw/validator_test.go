package mw

import (
	"context"
	"encoding/json"
	"testing"

	bplusc "gitlab.intelligentb.com/devops/bplus/context"
	"gitlab.intelligentb.com/devops/bplus/fw"
	e "gitlab.intelligentb.com/devops/bplus/internal/err"
	"gopkg.in/go-playground/assert.v1"
)

type A struct {
	B int64  `json:"b"`
	C string `json:"c" validate:"required"`
}

func TestNoPayloadv10validator(t *testing.T) {
	chain := fw.MakeChain()
	chain.Add(dummy)
	ctx := context.Background()

	ctx = v10validator(ctx, &chain)
	err := bplusc.GetError(ctx)
	assert.Equal(t, err, nil)
}

func TestValidv10validator(t *testing.T) {
	var request A

	chain := fw.MakeChain()
	chain.Add(dummy)
	ctx := context.Background()
	r := []byte(`{
    "b": 1,
    "c": "abc"
  }`)
	json.Unmarshal(r, &request)
	ctx = bplusc.SetPayload(ctx, request)
	ctx = v10validator(ctx, &chain)
	err := bplusc.GetError(ctx)
	assert.Equal(t, err, nil)
}

func TestInvalidv10validator(t *testing.T) {
	var request A
	chain := fw.MakeChain()
	chain.Add(dummy)
	ctx := context.Background()

	r := []byte(`{
    "b": 1
  }`)
	json.Unmarshal(r, &request)
	ctx = bplusc.SetPayload(ctx, request)
	ctx = v10validator(ctx, &chain)
	err := bplusc.GetError(ctx)
	assert.Equal(t, err, e.MakeBplusError(ctx, e.ValidationError, map[string]interface{}{}))
}

func dummy(ctx context.Context, _ *fw.MiddlewareChain) context.Context {
	return ctx
}
