package http

import (
	"context"
	"gitlab.intelligentb.com/devops/bplus/config"
	e "gitlab.intelligentb.com/devops/bplus/err"
	"gitlab.intelligentb.com/devops/bplus/fw"
	"gopkg.in/go-playground/assert.v1"
	"log"
	"net/http"
	"os"
	"testing"
)

type Service struct {}
type Input struct {
	In string
}
 func MakeInput(ctx context.Context)(interface{},error){
	return &Input{},nil
}
func MakeOutput(ctx context.Context)(interface{},error){
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

// the testing is internal (i.e. using the same http package) since the package does not expose public methods
func create()fw.ServiceDescriptor{
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
  	service := fw.ServiceDescriptor{
		Name:            "EchoService",
		Description:     "Echoes input to output",
		ServiceToInvoke: Service{},
		Operations: []fw.OperationDescriptor{od,},
  	}
	return service
}

func TestOperationSetup(t *testing.T) {
	os.Setenv("BPLUS.PORT","5000")
	fw.RegisterService("EchoService",create())

	go func (){
		a := ":" + config.Value("bplus.port")
		log.Printf("Starting server at address %s\n",a)
		http.ListenAndServe(a, HTTPHandler)
	}()

	// access the exposed service via proxy
	ret,err := ProxyRequest(context.TODO(),"EchoService","Echo",&Input{In:"hello"})
	if err != nil {
		log.Printf("Error in issuing an Http request. Error = %s",err.Error())
		t.Fail()
		return
	}
	out,ok := ret.(*Output)
	if !ok {
		log.Printf("Error in casting the output to type Output. Error = %s",err.Error())
		t.Fail()
		return
	}

	assert.Equal(t,out.Out,"hello")

	// give an error input
	ret,err = ProxyRequest(context.TODO(),"EchoService","Echo",&Input{In:"xxx"})
	if err == nil {
		log.Printf("Error in issuing an Http request. Error = %s",err.Error())
		t.Fail()
		return
	}

	bpluserr,ok := err.(e.BPlusError)
	if !ok {
		log.Printf("Error in casting the error to BPlusError. Error = %s",err.Error())
		t.Fail()
		return
	}
	// assert.Equal(t,bpluserr.ErrorCode,12000)
	assert.Equal(t,bpluserr.HTTPErrorCode,400)

}
