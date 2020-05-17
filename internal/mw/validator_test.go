package mw

import (
	"context"
	"encoding/json"
	wegoe "github.com/agorago/wego/err"
	"github.com/agorago/wego/log"
	"net/http"
	"testing"

	wegocontext "github.com/agorago/wego/context"
	"github.com/agorago/wego/fw"
	"gopkg.in/go-playground/assert.v1"
)

type A struct {
	B int64  `json:"b"`
	C string `json:"c" validate:"required"`
}

func TestNoPayloadv10validator(t *testing.T) {
	chain := fw.MakeChain()
	chain.Add(V10Validator{})

	ctx := context.Background()

	ctx = chain.DoContinue(ctx)
	err := wegocontext.GetError(ctx)
	assert.Equal(t, err, nil)
}

func TestValidv10validator(t *testing.T) {
	var request A

	chain := fw.MakeChain()
	chain.Add(V10Validator{})

	ctx := context.Background()
	r := []byte(`{
    "b": 1,
    "c": "abc"
  }`)
	json.Unmarshal(r, &request)
	ctx = wegocontext.SetPayload(ctx, request)
	ctx = chain.DoContinue(ctx)
	err := wegocontext.GetError(ctx)
	assert.Equal(t, err, nil)
}

func TestInvalidv10validator(t *testing.T) {
	var request A
	chain := fw.MakeChain()
	chain.Add(V10Validator{})

	ctx := context.Background()

	r := []byte(`{
    "b": 1
  }`)
	json.Unmarshal(r, &request)
	ctx = wegocontext.SetPayload(ctx, request)
	ctx = chain.DoContinue(ctx)
	err := wegocontext.GetError(ctx)
	e1, ok := err.(wegoe.WeGOError)
	if !ok {
		log.Errorf(ctx, "Cannot cast the error into WeGO error. Err = %#v", err)
		t.Fail()
		return
	}
	assert.Equal(t, e1.HTTPErrorCode, http.StatusBadRequest)
}

func terminator(ctx context.Context, _ *fw.MiddlewareChain) context.Context {
	return ctx
}
