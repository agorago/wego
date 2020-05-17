package fw_test

import (
	"context"
	"github.com/magiconair/properties/assert"
	wegoc "github.com/agorago/wego/context"
	"github.com/agorago/wego/fw"
	"testing"
)

type TestMiddleware struct{
	S string
}

func (tm TestMiddleware)Intercept(ctx context.Context, chain *fw.MiddlewareChain) context.Context {
	ctx = addString(ctx, "pre_"+tm.S)
	ctx = chain.DoContinue(ctx)
	ctx = addString(ctx, "post_"+tm.S)
	return ctx
}

type Terminator struct{}
func (Terminator)Intercept(ctx context.Context, chain *fw.MiddlewareChain) context.Context {
	return addString(ctx, "dest")
}


func addString(ctx context.Context, s string) context.Context {
	arr, _ := wegoc.Value(ctx, "STACK").([]string)
	arr = append(arr, s)
	return wegoc.Add(ctx, "STACK", arr)
}

func TestChain(t *testing.T) {
	chain := fw.MakeChain()
	chain.Add(TestMiddleware{"mw1",})
	chain.Add(TestMiddleware{"mw2",})
	chain.Add(Terminator{})
	ctx := chain.DoContinue(context.TODO())
	expected := []string{"pre_mw1", "pre_mw2", "dest", "post_mw2", "post_mw1"}
	assert.Equal(t, wegoc.Value(ctx, "STACK"), expected)
}
