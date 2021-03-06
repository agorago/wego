package mw_test

import (
	"context"
	"github.com/magiconair/properties/assert"
	wegocontext "github.com/agorago/wego/context"
	wegoe "github.com/agorago/wego/err"
	"github.com/agorago/wego/fw"
	"github.com/agorago/wego/internal/mw"
	"github.com/agorago/wego/internal/testutils"
	"log"
	"testing"
)

func setupContext(testString string) context.Context {
	ctx := context.TODO()
	sd := testutils.CreateEcho()
	od := sd.Operations[0]
	od.Service = sd
	ctx = fw.SetOperationDescriptor(ctx, od)

	input := testutils.Input{In: testString}
	ctx = wegocontext.SetPayload(ctx, &input)
	return ctx
}

func setupContextWithInvalidPayload(testString string) context.Context {
	ctx := context.TODO()
	sd := testutils.CreateEcho()
	od := sd.Operations[0]
	od.Service = sd
	ctx = fw.SetOperationDescriptor(ctx, od)

	input := map[string]interface{}{"In": testString}
	ctx = wegocontext.SetPayload(ctx, &input)
	return ctx
}

func setupContextWithBadMethodName(testString string) context.Context {
	ctx := context.TODO()
	sd := testutils.CreateEcho()
	od := sd.Operations[0]
	od.Service = sd
	od.Name = "BadMethod"
	ctx = fw.SetOperationDescriptor(ctx, od)

	input := testutils.Input{In: testString}
	ctx = wegocontext.SetPayload(ctx, &input)
	return ctx
}

func setupContextHeader(testString string) context.Context {
	ctx := setupContextHeaderNoValueSet(testString)
	ctx = wegocontext.Add(ctx, "Input", testString)
	return ctx
}

func setupContextHeaderNoValueSet(testString string) context.Context {
	ctx := context.TODO()
	sd := testutils.CreateEcho()
	od := sd.Operations[1]
	od.Service = sd
	ctx = fw.SetOperationDescriptor(ctx, od)
	return ctx
}

func TestServiceInvoker(t *testing.T) {
	const testing = "testing"
	ctx := setupContext(testing)
	si := mw.ServiceInvoker{}
	ctx = si.Intercept(ctx, nil)

	o := wegocontext.GetResponsePayload(ctx)
	output, ok := o.(testutils.Output)
	if !ok {
		log.Printf("Error in casting the output to type Output. Type = %#v\n", o)
		t.Fail()
		return
	}
	assert.Equal(t, output.Out, testing)
}

func TestServiceInvokerWithErrors(t *testing.T) {
	const testing = "xxx"
	ctx := setupContext(testing)
	si := mw.ServiceInvoker{}
	ctx = si.Intercept(ctx, nil)

	e := wegocontext.GetError(ctx)

	err1, ok := e.(wegoe.WeGOError)
	if !ok {
		log.Printf("Error in casting the error to type WeGOError. Type = %#v\n", e)
		t.Fail()
		return
	}
	assert.Equal(t, err1.HTTPErrorCode, 400)
}

func TestServiceInvokerHeader(t *testing.T) {
	const testing = "testing"
	ctx := setupContextHeader(testing)
	si := mw.ServiceInvoker{}
	ctx = si.Intercept(ctx, nil)

	o := wegocontext.GetResponsePayload(ctx)
	output, ok := o.(testutils.Output)
	if !ok {
		log.Printf("Error in casting the output to type Output. Type = %#v\n", o)
		t.Fail()
		return
	}
	assert.Equal(t, output.Out, testing)
}

func TestServiceInvokerHeaderNoValueSet(t *testing.T) {
	const testing = "testing"
	ctx := setupContextHeaderNoValueSet(testing)
	si := mw.ServiceInvoker{}
	ctx = si.Intercept(ctx, nil)

	e := wegocontext.GetError(ctx)

	err1, ok := e.(wegoe.WeGOError)
	if !ok {
		log.Printf("Error in casting the error to type WeGOError. Type = %#v\n", e)
		t.Fail()
		return
	}
	assert.Equal(t, err1.HTTPErrorCode, 400)
}

func TestServiceInvokerWithBadMethodName(t *testing.T) {
	const testing = "testing"
	ctx := setupContextWithBadMethodName(testing)
	si := mw.ServiceInvoker{}
	ctx = si.Intercept(ctx, nil)

	e := wegocontext.GetError(ctx)

	err1, ok := e.(wegoe.WeGOError)
	if !ok {
		log.Printf("Error in casting the error to type WeGOError. Type = %#v\n", e)
		t.Fail()
		return
	}
	log.Printf("error returned is %#v\n", err1)
	assert.Equal(t, err1.HTTPErrorCode, 500)
}

func TestServiceInvokerWithInvalidPayload(t *testing.T) {
	const testing = "testing"
	ctx := setupContextWithInvalidPayload(testing)
	si := mw.ServiceInvoker{}
	ctx = si.Intercept(ctx, nil)

	e := wegocontext.GetError(ctx)

	err1, ok := e.(wegoe.WeGOError)
	if !ok {
		log.Printf("Error in casting the error to type WeGOError. Type = %#v\n", e)
		t.Fail()
		return
	}
	log.Printf("error returned is %#v\n", err1)
	assert.Equal(t, err1.HTTPErrorCode, 500)
}
