package fw_test

import (
	"context"
	"github.com/magiconair/properties/assert"
	bplusc "gitlab.intelligentb.com/devops/bplus/context"
	"gitlab.intelligentb.com/devops/bplus/fw"
	"testing"
)

func generatMiddleware(s string) fw.Middleware {
	return func(ctx context.Context, chain *fw.MiddlewareChain) context.Context {
		ctx = addString(ctx, "pre_"+s)
		ctx = chain.DoContinue(ctx)
		ctx = addString(ctx, "post_"+s)
		return ctx
	}
}

func addString(ctx context.Context, s string) context.Context {
	arr, _ := bplusc.Value(ctx, "STACK").([]string)
	arr = append(arr, s)
	return bplusc.Add(ctx, "STACK", arr)
}

func TestChain(t *testing.T) {
	chain := fw.MakeChain()
	chain.Add(generatMiddleware("mw1"))
	chain.Add(generatMiddleware("mw2"))
	chain.Add(func(ctx context.Context, _ *fw.MiddlewareChain) context.Context {
		return addString(ctx, "dest")
	})
	ctx := chain.DoContinue(context.TODO())
	expected := []string{"pre_mw1", "pre_mw2", "dest", "post_mw2", "post_mw1"}
	assert.Equal(t, bplusc.Value(ctx, "STACK"), expected)
}
