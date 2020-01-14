package http

import (
	"context"
	"fmt"
	"reflect"
	"runtime"
	"time"

	bplusc "github.com/MenaEnergyVentures/bplus/context"
	fw "github.com/MenaEnergyVentures/bplus/fw"
	mw "github.com/MenaEnergyVentures/bplus/internal/mw"
)

// MakeProxy - Constructs a generic remote proxy based on the service descriptor
func proxyFunc(f interface{}) interface{} {
	rf := reflect.TypeOf(f)
	if rf.Kind() != reflect.Func {
		panic("expects a function")
	}
	vf := reflect.ValueOf(f)
	wrapperF := reflect.MakeFunc(rf, func(in []reflect.Value) []reflect.Value {
		start := time.Now()
		out := vf.Call(in)
		end := time.Now()
		fmt.Printf("calling %s took %v\n", runtime.FuncForPC(vf.Pointer()).Name(), end.Sub(start))
		return out
	})
	return wrapperF.Interface()
}

// doHttp - Given an operation descriptor invokes the
// operation remotely using HTTP
func doHTTP(ctx context.Context, od fw.OperationDescriptor, params []interface{}) (interface{}, error) {
	// instantiates a list of proxy middlewares
	// the last of which invokes the HTTP request
	ctx = fw.SetProxyOperationDescriptor(ctx, od)
	ctx = bplusc.SetProxyRequestParams(ctx, params)
	chain := fw.MakeChain()
	chain.Add(mw.HTTPInvoker)
	ctx = chain.DoContinue(ctx)
	// fmt.Printf("doHTTP: od is %#v\n", od)

	resp := bplusc.GetProxyResponsePayload(ctx)
	err := bplusc.GetProxyResponseError(ctx)
	return resp, err

}

// ProxyRequest - create a proxy that invokes the service and operation remotely.
func ProxyRequest(ctx context.Context, service string, operation string, params ...interface{}) (interface{}, error) {
	od, err := fw.FindOperationDescriptor(service, operation)
	if err != nil {
		return nil, err
	}
	return doHTTP(ctx, od, params)
}
