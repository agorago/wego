package testutils

import (
	"context"
	"fmt"
	wego "github.com/agorago/wego"
	"github.com/agorago/wego/config"
	e "github.com/agorago/wego/err"
	"github.com/agorago/wego/fw"
	"log"
	"net/http"
	"reflect"
)

// Set up a test echo service to facilitate testing of various components in WeGO
type Service struct{}
type Input struct {
	In string
}

func MakeInput(_ context.Context) (interface{}, error) {
	return &Input{}, nil
}
func MakeOutput(_ context.Context) (interface{}, error) {
	return &Output{}, nil
}

type Output struct {
	Out string
}

func (Service) Echo(ctx context.Context, input *Input) (Output, error) {
	if input.In == "xxx" {
		return Output{}, e.MakeErrWithHTTPCode(ctx, e.Error, 12000, "Invalid string xxx",
			400, nil)
	}
	return Output{input.In}, nil
}

func (Service) EchoString(_ context.Context, s string) (Output, error) {
	return Output{s}, nil
}

// the testing is internal (i.e. using the same http package) since the package does not expose public methods
func CreateEcho() fw.ServiceDescriptor {
	od := fw.OperationDescriptor{
		Name:            "Echo",
		URL:             "/echo",
		OpRequestMaker:  MakeInput,
		OpResponseMaker: MakeOutput,
		HTTPMethod:      "GET",
		Params: []fw.ParamDescriptor{
			{
				Name:        "ctx",
				ParamOrigin: fw.CONTEXT,
				Description: "context",
			},
			{
				Name:        "input",
				ParamOrigin: fw.PAYLOAD,
				Description: "Request",
			},
		},
	}
	odstring := fw.OperationDescriptor{
		Name:            "EchoString",
		URL:             "/echo-string",
		OpResponseMaker: MakeOutput,
		HTTPMethod:      "GET",
		Params: []fw.ParamDescriptor{
			{
				Name:        "ctx",
				ParamOrigin: fw.CONTEXT,
				Description: "context",
			},
			{
				Name:        "Input",
				ParamOrigin: fw.HEADER,
				Description: "Request",
				ParamKind:   reflect.String,
			},
		},
	}
	service := fw.ServiceDescriptor{
		Name:            "EchoService",
		Description:     "Echoes input to output",
		ServiceToInvoke: Service{},
		Operations:      []fw.OperationDescriptor{od, odstring},
	}
	return service
}

func StartServer() (fw.RegistrationService,fw.ServiceDescriptor){
	commandCatalog,err := fw.MakeInitializedCommandCatalog(wego.MakeWegoInitializer())
	if err != nil {
		panic(fmt.Sprintf("Application cannot start. Wego has not been able to start. Error = %v",err))
	}
	sd := CreateEcho()
	rs := commandCatalog.Command(wego.WEGO).(fw.RegistrationService)
	rs.RegisterService("EchoService", sd)

	httphandler := commandCatalog.Command(wego.WegoHTTPHandler).(http.Handler)

	go func() {
		a := ":" + config.Value("bplus.port")
		log.Printf("Starting server at address %s\n", a)
		http.ListenAndServe(a, httphandler)
	}()
	return rs,sd
}
