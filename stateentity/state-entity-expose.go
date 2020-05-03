package stateentity

import (
	"context"
	"reflect"

	wegocontext "github.com/agorago/wego/context"
	fw "github.com/agorago/wego/fw"
	e "github.com/agorago/wego/internal/err"
)

func (str SubTypeRegistration) setupStateEntityService() {
	var subtypemaker = func(ctx context.Context) (interface{}, error) {
		return str.StateEntitySubTypeMaker(ctx)
	}

	var processParamMaker = func(ctx context.Context) (interface{}, error) {
		var err error

		stm, err := str.OrderSTMChooser(ctx)
		if err != nil {
			return nil, e.MakeBplusError(ctx, e.ErrorInObtainingSTM, map[string]interface{}{})
		}
		eventID := wegocontext.Value(ctx, "Eventid")
		if eventID == nil {
			return nil, e.MakeBplusError(ctx, e.EventNotFoundInRequest, map[string]interface{}{})
		}

		eID, ok := eventID.(string)
		if !ok {
			return nil, e.MakeBplusError(ctx, e.EventNotFoundInRequest, map[string]interface{}{})
		}

		var ptm = stm.ParamTypeMaker(eID)
		if ptm == nil {
			return nil, e.MakeBplusError(ctx, e.EventNotFoundInRequest, map[string]interface{}{})
		}

		return ptm.MakeParam(ctx)
	}
	var pdcreate = []fw.ParamDescriptor{
		fw.ParamDescriptor{
			Name:        "ctx",
			ParamOrigin: fw.CONTEXT,
		},
		fw.ParamDescriptor{
			Name:        "request",
			ParamOrigin: fw.PAYLOAD,
		},
	}
	var pdprocess = []fw.ParamDescriptor{
		fw.ParamDescriptor{
			Name:        "ctx",
			ParamOrigin: fw.CONTEXT,
		},
		fw.ParamDescriptor{
			Name:        "Orderid",
			ParamOrigin: fw.HEADER,
			ParamKind:   reflect.String,
		},
		fw.ParamDescriptor{
			Name:        "Eventid",
			ParamOrigin: fw.HEADER,
			ParamKind:   reflect.String,
		},
		fw.ParamDescriptor{
			Name:        "param",
			ParamOrigin: fw.PAYLOAD,
		},
	}

	var ods = []fw.OperationDescriptor{
		fw.OperationDescriptor{
			Name:            "Create",
			URL:             "/create",
			HTTPMethod:      "POST",
			OpRequestMaker:  subtypemaker,
			OpResponseMaker: subtypemaker,
			Params:          pdcreate,
		},
		fw.OperationDescriptor{
			Name:            "Process",
			URL:             "/process",
			HTTPMethod:      "POST",
			OpRequestMaker:  processParamMaker,
			OpResponseMaker: subtypemaker,
			Params:          pdprocess,
		},
	}

	var sd = fw.ServiceDescriptor{
		ServiceToInvoke: str,
		Name:            str.Name,
		Operations:      ods,
	}
	fw.RegisterService(str.Name, sd)
}
