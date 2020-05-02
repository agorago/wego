package mw_test

import (
	"context"
	"fmt"
	bplusc "gitlab.intelligentb.com/devops/bplus/context"
	bpluserr "gitlab.intelligentb.com/devops/bplus/err"
	"gitlab.intelligentb.com/devops/bplus/fw"
	"gitlab.intelligentb.com/devops/bplus/internal/mw"
	"gitlab.intelligentb.com/devops/bplus/internal/testutils"
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
	ctx = bplusc.SetPayload(ctx, &input)
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

	o := bplusc.GetResponsePayload(ctx)
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

	er := bplusc.GetError(ctx)
	er1, ok := er.(bpluserr.BPlusError)
	if !ok {
		log.Printf("Error in casting the error to type BPlusError. Type = %#v\n", er)
		t.Fail()
		return
	}
	fmt.Printf("%#v\n", er1)
	assert.Equal(t, er1.HTTPErrorCode, http.StatusBadRequest)
}
