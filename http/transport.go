package http

import (
	"context"
	"encoding/json"
	"net/http"

	"gitlab.intelligentb.com/devops/bplus/log"

	"github.com/gorilla/mux"

	bplusc "gitlab.intelligentb.com/devops/bplus/context"
	bpluserr "gitlab.intelligentb.com/devops/bplus/err"
	fw "gitlab.intelligentb.com/devops/bplus/fw"
	mw "gitlab.intelligentb.com/devops/bplus/internal/mw"
	nr "gitlab.intelligentb.com/devops/bplus/internal/nr"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	nrgorilla "github.com/newrelic/go-agent/v3/integrations/nrgorilla"
)

// http package allows the hooks to HTTP. It creates a transport from
// server configuration

// HTTPHandler - Grab hold of this to set up HTTP routes
var HTTPHandler *mux.Router

func init() {
	HTTPHandler = mux.NewRouter()
	HTTPHandler.Use(nrgorilla.Middleware(nr.NRApp))
	// register the HTTP registration with BPlus as an extension
	fw.RegisterOperations(setupOperation)
}

type httpGenericResponse struct {
	resp interface{}
	err  error
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
	log.Infof(context.Background(), "setting up the service for %s/%s\n", od.Service.Name, od.URL)
	HTTPHandler.Methods(od.HTTPMethod).PathPrefix("/" + od.Service.Name).Path(od.URL).Handler(handler)
}

func (hod httpod) makeEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		r := request.(*http.Request)
		ctx = fw.SetOperationDescriptor(ctx, hod.Od)
		ctx = bplusc.Enhance(ctx, r)
		ctx = bplusc.SetPayload(ctx, r.Body)

		resp, err := mw.Entrypoint(ctx)
		return httpGenericResponse{resp, err}, nil
	}
}

func (hod httpod) decodeRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	return r, nil
}

func encodeGenericResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json")

	gresp := response.(httpGenericResponse)
	if gresp.err != nil {
		return handleError(w, gresp.err)
	}

	if gresp.resp != nil {
		return json.NewEncoder(w).Encode(gresp.resp)
	} else {
		w.WriteHeader(http.StatusNoContent)
		return nil
	}
}

func handleError(w http.ResponseWriter, err error) error {
	e, ok := err.(bpluserr.HttpCodeProvider)
	if ok {
		w.WriteHeader(e.GetHttpCode())
		w.Write([]byte(err.Error()))
		return nil
	}
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(err.Error()))
	return nil
}
