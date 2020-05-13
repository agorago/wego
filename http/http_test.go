package http_test

import (
	"context"
	e "github.com/agorago/wego/err"
	wegohttp "github.com/agorago/wego/http"
	"github.com/agorago/wego/internal/testutils"
	"gopkg.in/go-playground/assert.v1"
	"log"
	"os"
	"testing"
)

func TestOperationSetup(t *testing.T) {
	os.Setenv("BPLUS.PORT", "5000")
	defer func() {
		os.Unsetenv("BPLUS.PORT")
	}()
	rs,_ := testutils.StartServer()
	proxy := wegohttp.MakeProxyService(rs)

	// access the exposed service via proxy
	ret, err := proxy.ProxyRequest(context.TODO(), "EchoService", "Echo", &testutils.Input{In: "hello"})
	if err != nil {
		log.Printf("Error in issuing an Http request. Error = %s", err.Error())
		t.Fail()
		return
	}
	out, ok := ret.(*testutils.Output)
	if !ok {
		log.Printf("Error in casting the output to type Output. Error = %s", err.Error())
		t.Fail()
		return
	}

	assert.Equal(t, out.Out, "hello")

	// give an error input
	ret, err = proxy.ProxyRequest(context.TODO(), "EchoService", "Echo", &testutils.Input{In: "xxx"})
	if err == nil {
		log.Printf("Error in issuing an Http request. Error = %s", err.Error())
		t.Fail()
		return
	}

	bpluserr, ok := err.(e.WeGOError)
	if !ok {
		log.Printf("Error in casting the error to WeGOError. Error = %s", err.Error())
		t.Fail()
		return
	}

	assert.Equal(t, bpluserr.HTTPErrorCode, 400)

}
