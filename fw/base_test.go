package fw_test

import (
	"context"
	"github.com/go-playground/assert/v2"
	wegoe "github.com/agorago/wego/err"
	"github.com/agorago/wego/fw"
	e "github.com/agorago/wego/internal/err"
	"log"
	"testing"
)

type myTestRequest struct {
	ID   string
	Desc string
}

func requestMaker(ctx context.Context) (interface{}, error) {
	return myTestRequest{}, nil
}

type myTestResponse struct {
	ID   string
	Desc string
}

func responseMaker(ctx context.Context) (interface{}, error) {
	return myTestResponse{}, nil
}

type myTestServiceApi interface {
	MyTestMethod(ctx context.Context, m *myTestRequest) (myTestResponse, error)
}

type myTestServiceImpl struct{}

func loggerMiddleware(ctx context.Context, chain *fw.MiddlewareChain) context.Context {
	log.Printf("pre log message")
	ctx = chain.DoContinue(ctx)
	log.Printf("post log message")
	return ctx
}

func registerOps(visitedOp *string) {
	fw.RegisterOperations(
		func(descriptor fw.OperationDescriptor) {
			*visitedOp = descriptor.Name
		})
}

func (my myTestServiceImpl) MyTestMethod(ctx context.Context, m *myTestRequest) (myTestResponse, error) {

	return myTestResponse{
		ID:   "resp" + m.ID,
		Desc: "resp" + m.Desc,
	}, nil
}

var myTestInstance = myTestServiceImpl{}

func registerService(clientSide bool) fw.ServiceDescriptor {
	pds := []fw.ParamDescriptor{
		{
			Name:        "ctx",
			ParamOrigin: fw.CONTEXT,
			Description: "context object",
		},
		{
			Name:        "m",
			ParamOrigin: fw.PAYLOAD,
			Description: "myTestRequest",
		},
	}
	ods := []fw.OperationDescriptor{
		{
			Name:                "MyTestMethod", // the actual name of the test above
			Description:         "This is what this test does. Can be the method description from above",
			Service:             fw.ServiceDescriptor{},
			RequestDescription:  "Same as comment for myTestRequest above",
			ResponseDescription: "Same as comment for myTestResponse above",
			URL:                 "/mytestmethod", // the url to invoke
			OpRequestMaker:      requestMaker,
			OpResponseMaker:     responseMaker,
			OpMiddleware:        []fw.Middleware{loggerMiddleware}, // set of middlewares that needs to be invoked at the server
			HTTPMethod:          "GET",
			ProxyMiddleware:     []fw.Middleware{loggerMiddleware}, // set of middlewares that needs to be invoked at the proxy
			Params:              pds,
		},
	}

	sd := fw.ServiceDescriptor{
		Name:        "test",
		Description: "test description",
		Operations:  ods,
	}
	if !clientSide {
		sd.ServiceToInvoke = myTestInstance
	}
	fw.RegisterService("TestService", sd)
	return sd
}

func TestRegisterService(t *testing.T) {
	var s string
	registerOps(&s)
	sd := registerService(false)
	x, err := fw.FindServiceDescriptor("TestService")
	if err != nil {
		log.Printf("Unable to find service descriptor.Error = %s\n", err.Error())
		t.Fail()
	}
	od, err := fw.FindOperationDescriptor("TestService", "MyTestMethod")
	if err != nil {
		log.Printf("Unable to find operation descriptor.Error = %s\n", err.Error())
		t.Fail()
	}
	assert.Equal(t, sd.Operations[0].Name, s)
	assert.Equal(t, sd, x)
	assert.Equal(t, sd.Operations[0].Name, od.Name)
}

func TestRegisterServiceForClientSide(t *testing.T) {
	var s string
	registerOps(&s)
	sd := registerService(true)

	assert.NotEqual(t, sd.Operations[0].Name, s)
	_, err := fw.FindOperationDescriptor("TestService", "foo")
	if err == nil {
		log.Printf("It is not expected to find the foo method\n")
		t.Fail()
	}
	bpluserr, ok := err.(wegoe.BPlusError)
	if !ok {
		log.Printf("Error returned is expected to be of type BPlusError\n")
		t.Fail()
	}
	assert.Equal(t, bpluserr.ErrorCode, int(e.OperationNotFound))
}

func TestFindServiceDescriptorService(t *testing.T) {
	registerService(false)
	_, err := fw.FindServiceDescriptor("xxx")
	if err == nil {
		log.Printf("It is not expected to find the xxx service\n")
		t.Fail()
	}
}
