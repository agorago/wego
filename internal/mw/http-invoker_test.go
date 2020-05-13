package mw_test

import (
	"context"
	"fmt"
	wegocontext "github.com/agorago/wego/context"
	wegoe "github.com/agorago/wego/err"
	"github.com/agorago/wego/fw"
	"github.com/agorago/wego/internal/mw"
	"github.com/agorago/wego/internal/testutils"
	"gopkg.in/go-playground/assert.v1"
	"log"
	"net/http"
	"os"
	"testing"
)

func setupServer(testString string) context.Context {

	sd := testutils.CreateEcho()
	testutils.StartServer()
	ctx := context.TODO()

	od := sd.Operations[0]
	od.Service = sd
	ctx = fw.SetOperationDescriptor(ctx, od)

	input := testutils.Input{In: testString}
	ctx = wegocontext.SetPayload(ctx, &input)
	return ctx
}

func TestHTTPInvoker(t *testing.T) {
	os.Setenv("BPLUS.PORT", "5000")
	defer func() {
		os.Unsetenv("BPLUS.PORT")
	}()
	const testing = "testing"
	ctx := setupServer(testing)
	ctx = mw.HTTPInvoker(ctx, nil)

	o := wegocontext.GetResponsePayload(ctx)
	output, ok := o.(*testutils.Output)
	if !ok {
		log.Printf("Error in casting the output to type Output. Type = %#v\n", o)
		t.Fail()
		return
	}
	assert.Equal(t, output.Out, testing)
}

func TestHTTPInvokerWithError(t *testing.T) {
	os.Setenv("BPLUS.PORT", "5000")
	defer func() {
		os.Unsetenv("BPLUS.PORT")
	}()
	const testing = "xxx"
	ctx := setupServer(testing)
	ctx = mw.HTTPInvoker(ctx, nil)

	er := wegocontext.GetError(ctx)
	er1, ok := er.(wegoe.WeGOError)
	if !ok {
		log.Printf("Error in casting the error to type WeGOError. Type = %#v\n", er)
		t.Fail()
		return
	}
	fmt.Printf("%#v\n", er1)
	assert.Equal(t, er1.HTTPErrorCode, http.StatusBadRequest)
}
