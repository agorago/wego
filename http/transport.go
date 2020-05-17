package http

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/agorago/wego/log"

	"github.com/gorilla/mux"

	wegoc "github.com/agorago/wego/context"
	wegoe "github.com/agorago/wego/err"
	fw "github.com/agorago/wego/fw"
	mw "github.com/agorago/wego/internal/mw"
	nr "github.com/agorago/wego/internal/nr"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	nrgorilla "github.com/newrelic/go-agent/v3/integrations/nrgorilla"
)

// http package allows the hooks to HTTP. It enables HTTP transport for all registered
// WEGO services. To facilitate this, it registers itself as a transport in WEGO. Whenever
// a new service is registered, this module is called to register all the operations


func InitializeHTTP(rs fw.RegistrationService, ep mw.Entrypoint)http.Handler {
	HTTPHandler := mux.NewRouter()
	HTTPHandler.Use(nrgorilla.Middleware(nr.NRApp))
	// register the HTTP transport with WEGO
	rs.RegisterTransport("HTTP",func (od fw.OperationDescriptor){
		if od.URL == "" {
			return
		}
		handler := setupHTTPForOperation(od,ep)
		HTTPHandler.Methods(od.HTTPMethod).PathPrefix("/" + od.Service.Name).Path(od.URL).Handler(handler)
	})
	return HTTPHandler
}

type httpGenericResponse struct {
	resp interface{}
	err  error
}

type httpod struct {
	Od fw.OperationDescriptor
	Ep mw.Entrypoint
}

func setupHTTPForOperation(od fw.OperationDescriptor,ep mw.Entrypoint) *httptransport.Server{
	hod := httpod{Od: od,Ep: ep,}
	handler := httptransport.NewServer(
		hod.makeEndpoint(),
		hod.decodeRequest,
		encodeGenericResponse,
	)
	log.Infof(context.Background(), "setting up the service for %s%s\n", od.Service.Name, od.URL)
	return handler
}

func (hod httpod) makeEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		r := request.(*http.Request)
		ctx = fw.SetOperationDescriptor(ctx, hod.Od)
		ctx = wegoc.Enhance(ctx, r)
		ctx = wegoc.SetPayload(ctx, r.Body)

		resp, err := hod.Ep.Invoke(ctx)
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
	e, ok := err.(wegoe.HttpCodeProvider)
	if ok {
		w.WriteHeader(e.GetHttpCode())
		w.Write([]byte(err.Error()))
		return nil
	}
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(err.Error()))
	return nil
}
