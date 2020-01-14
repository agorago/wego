package mw

import (
	"context"
	"reflect"

	bplusc "github.com/MenaEnergyVentures/bplus/context"
	fw "github.com/MenaEnergyVentures/bplus/fw"
	util "github.com/MenaEnergyVentures/bplus/util"
)

// since this is the last middleware we would not invoke the chain anymore
func serviceInvoker(ctx context.Context, _ *fw.MiddlewareChain) context.Context {
	od := fw.GetOperationDescriptor(ctx)
	args := makeArgs(ctx, od)
	v, err := invoke(od.Service.ServiceToInvoke, od.Name, args)
	ctx = bplusc.SetResponsePayload(ctx, v)
	if err != nil {
		ctx = bplusc.SetError(ctx, err)
	}
	return ctx
}

func makeArgs(ctx context.Context, od fw.OperationDescriptor) []interface{} {
	var args []interface{}
	for _, pd := range od.Params {
		args = append(args, makeArg(ctx, pd))
	}
	return args
}

func makeArg(ctx context.Context, param fw.ParamDescriptor) interface{} {
	switch param.ParamOrigin {
	case fw.HEADER:
		return util.ConvertFromString(bplusc.Value(ctx, param.Name).(string), param.ParamKind)
	case fw.PAYLOAD:
		return bplusc.GetPayload(ctx)
	case fw.CONTEXT:
		return ctx
	}
	return nil
}

func invoke(any interface{}, name string, args []interface{}) (interface{}, error) {
	inputs := make([]reflect.Value, len(args))
	for i := range args {
		inputs[i] = reflect.ValueOf(args[i])
	}
	x := reflect.ValueOf(any).MethodByName(name).Call(inputs)
	retVal := x[0].Interface()
	retErr := x[1].Interface()
	if retErr == nil {
		return retVal, nil
	}
	return retVal, retErr.(error)
}
