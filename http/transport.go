package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	bplusc "github.com/MenaEnergyVentures/bplus/context"
	fw "github.com/MenaEnergyVentures/bplus/fw"
	mw "github.com/MenaEnergyVentures/bplus/internal/mw"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
)

// http package allows the hooks to HTTP. It creates a transport from
// server configuration

// HTTPHandler - Grab hold of this to set up HTTP routes
var HTTPHandler *mux.Router

func init() {
	HTTPHandler = mux.NewRouter()
	// register the operation with BPlus as an extension
	fw.RegisterOperations(setupOperation)
}

type httpod struct {
	Od fw.OperationDescriptor
}

func setupOperation(od fw.OperationDescriptor) {
	hod := httpod{Od: od}
	if od.URL == "" {
		return
	}
	handler := httptransport.NewServer(
		hod.makeEndpoint(),
		hod.decodeRequest,
		encodeGenericResponse,
	)
	fmt.Printf("setting up the service for %s\n", od.URL)
	HTTPHandler.Methods(od.HTTPMethod).Path(od.URL).Handler(handler)
}

func (hod httpod) makeEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		r := request.(*http.Request)
		ctx = fw.SetOperationDescriptor(ctx, hod.Od)
		ctx = bplusc.Enhance(ctx, r)
		ctx = bplusc.SetPayload(ctx, r.Body)

		return mw.Entrypoint(ctx)
	}
}

func (hod httpod) decodeRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	return r, nil
}

func encodeGenericResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
