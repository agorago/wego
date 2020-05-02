package mw

import (
	"context"
	"encoding/json"
	bpluse "gitlab.intelligentb.com/devops/bplus/err"
	"gitlab.intelligentb.com/devops/bplus/log"
	"net/http"
	"testing"

	bplusc "gitlab.intelligentb.com/devops/bplus/context"
	"gitlab.intelligentb.com/devops/bplus/fw"
	"gopkg.in/go-playground/assert.v1"
)

type A struct {
	B int64  `json:"b"`
	C string `json:"c" validate:"required"`
}

func TestNoPayloadv10validator(t *testing.T) {
	chain := fw.MakeChain()
	chain.Add(v10validator)
	chain.Add(terminator)
	ctx := context.Background()

	ctx = chain.DoContinue(ctx)
	err := bplusc.GetError(ctx)
	assert.Equal(t, err, nil)
}

func TestValidv10validator(t *testing.T) {
	var request A

	chain := fw.MakeChain()
	chain.Add(v10validator)
	chain.Add(terminator)
	ctx := context.Background()
	r := []byte(`{
    "b": 1,
    "c": "abc"
  }`)
	json.Unmarshal(r, &request)
	ctx = bplusc.SetPayload(ctx, request)
	ctx = chain.DoContinue(ctx)
	err := bplusc.GetError(ctx)
	assert.Equal(t, err, nil)
}

func TestInvalidv10validator(t *testing.T) {
	var request A
	chain := fw.MakeChain()
	chain.Add(v10validator)
	chain.Add(terminator)
	ctx := context.Background()

	r := []byte(`{
    "b": 1
  }`)
	json.Unmarshal(r, &request)
	ctx = bplusc.SetPayload(ctx, request)
	ctx = chain.DoContinue(ctx)
	err := bplusc.GetError(ctx)
	e1, ok := err.(bpluse.BPlusError)
	if !ok {
		log.Errorf(ctx, "Cannot cast the error into Bplus error. Err = %#v", err)
		t.Fail()
		return
	}
	assert.Equal(t, e1.HTTPErrorCode, http.StatusBadRequest)
}

func terminator(ctx context.Context, _ *fw.MiddlewareChain) context.Context {
	return ctx
}
