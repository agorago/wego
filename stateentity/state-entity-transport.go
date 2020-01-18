package stateentity

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	bplusHTTP "github.com/MenaEnergyVentures/bplus/http"
	e "github.com/MenaEnergyVentures/bplus/internal/err"
	"github.com/MenaEnergyVentures/bplus/stm"
	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
)

// This file contains DEAD CODE - Not used anymore. BPLUS exposes the transport
// see state-entity-expose.go
func (str SubTypeRegistration) makeCreateEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(stm.StateEntity)
		v, err := str.Create(ctx, req)

		if err != nil {
			return v, err
		}
		return v, nil
	}
}

// ProcessRequest - the request exposed by the Process() method
type ProcessRequest struct {
	EntityID string
	EventID  string
	Param    interface{}
}

func (str SubTypeRegistration) makeProcessEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(ProcessRequest)

		order, err := str.Process(ctx, req.EntityID, req.EventID, req.Param)
		if err != nil {
			return order, err
		}
		return order, nil
	}
}

func (str SubTypeRegistration) decodeCreateRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	request, err := str.StateEntitySubTypeMaker(ctx)
	if err != nil {
		return nil, e.MakeBplusError(ctx, e.ErrorInDecoding)
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func (str SubTypeRegistration) decodeProcessRequest(context context.Context, r *http.Request) (interface{}, error) {

	var request ProcessRequest

	request.EntityID = r.Header.Get("orderId")
	request.EventID = r.Header.Get("eventId")
	var err error
	// request.Param, err = ioutil.ReadAll(r.Body)
	stm, err := str.OrderSTMChooser(context)
	if err != nil {
		fmt.Printf("Problem with obtaining stm\n")
	}
	var ptm = stm.ParamTypeMaker(request.EventID)
	if ptm == nil {
		return request, nil
	}

	obj, err := ptm.MakeParam(context)
	//fmt.Printf("Got the following from ParamTypeMaker: %v. Type = %v\n", obj,
	//	reflect.TypeOf(obj))
	if err != nil {
		return nil, err
	}

	if err := json.NewDecoder(r.Body).Decode(&obj); err != nil {
		return nil, err
	}

	//fmt.Printf("After decoding the parameter. The parameter value is %v\n", obj)
	request.Param = obj
	return request, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

// setupOrderService - exposes the state entity service at the  end point
func (str SubTypeRegistration) setupStateEntityService1() {

	createHandler := httptransport.NewServer(
		str.makeCreateEndpoint(),
		str.decodeCreateRequest,
		encodeResponse,
	)

	processHandler := httptransport.NewServer(
		str.makeProcessEndpoint(),
		str.decodeProcessRequest,
		encodeResponse,
	)

	bplusHTTP.HTTPHandler.Methods("POST").Path(fmt.Sprintf("/%s/create", str.URLPrefix)).
		Handler(createHandler)
	bplusHTTP.HTTPHandler.Methods("POST").Path(fmt.Sprintf("/%s/process", str.URLPrefix)).
		Handler(processHandler)

}
