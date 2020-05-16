package mw

import (
	"context"
	"github.com/agorago/wego/log"
	"net/http"
	"reflect"

	wegocontext "github.com/agorago/wego/context"
	fw "github.com/agorago/wego/fw"
	e "github.com/agorago/wego/internal/err"
	util "github.com/agorago/wego/util"
)

// since this is the last middleware we would not invoke the chain anymore
func ServiceInvoker(ctx context.Context, _ *fw.MiddlewareChain) context.Context {
	var hasResponse = false

	od := fw.GetOperationDescriptor(ctx)
	args, err := makeArgs(ctx, od)
	if err != nil {
		ctx = wegocontext.SetError(ctx, err)
		return ctx
	}

	if od.OpResponseMaker != nil {
		hasResponse = true
	}

	v, err := invoke(ctx, od.Service.ServiceToInvoke, od.Name, args, hasResponse)
	ctx = wegocontext.SetResponsePayload(ctx, v)

	if err != nil {
		ctx = wegocontext.SetError(ctx, err)
	}
	return ctx
}

func makeArgs(ctx context.Context, od fw.OperationDescriptor) ([]interface{}, error) {
	var args []interface{}
	for _, pd := range od.Params {
		arg, err := makeArg(ctx, pd)
		if err != nil {
			return nil, err
		}
		args = append(args, arg)
	}
	return args, nil
}

func makeArg(ctx context.Context, param fw.ParamDescriptor) (interface{}, error) {
	switch param.ParamOrigin {
	case fw.HEADER:
		s, ok := wegocontext.Value(ctx, param.Name).(string)
		if !ok {
			return nil, e.HTTPError(ctx, http.StatusBadRequest,
				e.ParameterMissingInRequest, map[string]interface{}{
					"Param": param.Name})
		}
		return util.ConvertFromString(s, param.ParamKind), nil
	case fw.PAYLOAD:
		return wegocontext.GetPayload(ctx), nil
	case fw.CONTEXT:
		return ctx, nil
	}
	return nil, nil
}

// We use named return values to ensure that the defer() function is able to recover from
// panic and return sensible values
func invoke(ctx context.Context, any interface{}, name string, args []interface{}, hasResponse bool) (ret interface{}, retErr error) {
	inputs := make([]reflect.Value, len(args))
	for i := range args {
		inputs[i] = reflect.ValueOf(args[i])
	}
	defer func() {
		if r := recover(); r != nil {
			log.Errorf(ctx, "Service Invocation Exception. Panic'ed during reflection.  Error = %#v\n", r)
			retErr = e.HTTPError(ctx, http.StatusInternalServerError,
				e.ErrorInInvokingService, map[string]interface{}{
					"OperationName": name, "Error": r})
		}
	}()
	x := reflect.ValueOf(any).MethodByName(name).Call(inputs)
	return serviceResponse(x, hasResponse)
}

func serviceResponse(x []reflect.Value, hasResponse bool) (interface{}, error) {
	var retVal, retErr interface{}

	if hasResponse {
		retVal = x[0].Interface()
		retErr = x[1].Interface()
	} else {
		retVal = nil
		retErr = x[0].Interface()
	}

	if retErr == nil {
		return retVal, nil
	}
	return retVal, retErr.(error)
}
