package stateentity

import (
	"context"
	"reflect"

	bplusc "gitlab.intelligentb.com/devops/bplus/context"
	bplus "gitlab.intelligentb.com/devops/bplus/fw"
	e "gitlab.intelligentb.com/devops/bplus/internal/err"
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
		eventID := bplusc.Value(ctx, "Eventid")
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
	var pdcreate = []bplus.ParamDescriptor{
		bplus.ParamDescriptor{
			Name:        "ctx",
			ParamOrigin: bplus.CONTEXT,
		},
		bplus.ParamDescriptor{
			Name:        "request",
			ParamOrigin: bplus.PAYLOAD,
		},
	}
	var pdprocess = []bplus.ParamDescriptor{
		bplus.ParamDescriptor{
			Name:        "ctx",
			ParamOrigin: bplus.CONTEXT,
		},
		bplus.ParamDescriptor{
			Name:        "Orderid",
			ParamOrigin: bplus.HEADER,
			ParamKind:   reflect.String,
		},
		bplus.ParamDescriptor{
			Name:        "Eventid",
			ParamOrigin: bplus.HEADER,
			ParamKind:   reflect.String,
		},
		bplus.ParamDescriptor{
			Name:        "param",
			ParamOrigin: bplus.PAYLOAD,
		},
	}

	var ods = []bplus.OperationDescriptor{
		bplus.OperationDescriptor{
			Name:            "Create",
			URL:             "/create",
			HTTPMethod:      "POST",
			OpRequestMaker:  subtypemaker,
			OpResponseMaker: subtypemaker,
			Params:          pdcreate,
		},
		bplus.OperationDescriptor{
			Name:            "Process",
			URL:             "/process",
			HTTPMethod:      "POST",
			OpRequestMaker:  processParamMaker,
			OpResponseMaker: subtypemaker,
			Params:          pdprocess,
		},
	}

	var sd = bplus.ServiceDescriptor{
		ServiceToInvoke: str,
		Name:            str.Name,
		Operations:      ods,
	}
	bplus.RegisterService(str.Name, sd)
}
