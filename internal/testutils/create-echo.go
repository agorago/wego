package testutils

import (
	"context"
	e "gitlab.intelligentb.com/devops/bplus/err"
	"gitlab.intelligentb.com/devops/bplus/fw"
	"reflect"
)

// Set up a test echo service to facilitate testing of various components in BPlus
type Service struct {}
type Input struct {
	In string
}
func MakeInput(_ context.Context)(interface{},error){
	return &Input{},nil
}
func MakeOutput(_ context.Context)(interface{},error){
	return &Output{},nil
}
type Output struct {
	Out string
}
func (Service) Echo(ctx context.Context, input *Input)(Output,error){
	if input.In == "xxx"{
		return Output{},e.MakeErrWithHTTPCode(ctx,e.Error,12000,"Invalid string xxx",
			400,nil)
	}
	return Output{input.In},nil
}

func (Service) EchoString(_ context.Context, s string)(Output,error){
	return Output{s},nil
}

// the testing is internal (i.e. using the same http package) since the package does not expose public methods
func CreateEcho()fw.ServiceDescriptor{
	od := fw.OperationDescriptor{
		Name: "Echo",
		URL:             "/echo",
		OpRequestMaker:  MakeInput,
		OpResponseMaker: MakeOutput,
		HTTPMethod:      "GET",
		Params:          []fw.ParamDescriptor{
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
		Name: "EchoString",
		URL:             "/echo-string",
		OpResponseMaker: MakeOutput,
		HTTPMethod:      "GET",
		Params:          []fw.ParamDescriptor{
			{
				Name:        "ctx",
				ParamOrigin: fw.CONTEXT,
				Description: "context",
			},
			{
				Name:        "Input",
				ParamOrigin: fw.HEADER,
				Description: "Request",
				ParamKind: reflect.String,
			},
		},
	}
	service := fw.ServiceDescriptor{
		Name:            "EchoService",
		Description:     "Echoes input to output",
		ServiceToInvoke: Service{},
		Operations: []fw.OperationDescriptor{od,odstring},
	}
	return service
}

func StartServer(){
	fw.RegisterService("EchoService",CreateEcho())


}
