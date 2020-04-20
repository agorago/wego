package http

import (
	"context"
	"gitlab.intelligentb.com/devops/bplus/config"
	e "gitlab.intelligentb.com/devops/bplus/err"
	"gitlab.intelligentb.com/devops/bplus/fw"
	"gitlab.intelligentb.com/devops/bplus/internal/testutils"
	"gopkg.in/go-playground/assert.v1"
	"log"
	"net/http"
	"os"
	"testing"
)



func TestOperationSetup(t *testing.T) {
	os.Setenv("BPLUS.PORT","5000")
	fw.RegisterService("EchoService",testutils.CreateEcho())

	go func (){
		a := ":" + config.Value("bplus.port")
		log.Printf("Starting server at address %s\n",a)
		http.ListenAndServe(a, HTTPHandler)
	}()

	// access the exposed service via proxy
	ret,err := ProxyRequest(context.TODO(),"EchoService","Echo",&testutils.Input{In:"hello"})
	if err != nil {
		log.Printf("Error in issuing an Http request. Error = %s",err.Error())
		t.Fail()
		return
	}
	out,ok := ret.(*testutils.Output)
	if !ok {
		log.Printf("Error in casting the output to type Output. Error = %s",err.Error())
		t.Fail()
		return
	}

	assert.Equal(t,out.Out,"hello")

	// give an error input
	ret,err = ProxyRequest(context.TODO(),"EchoService","Echo",&testutils.Input{In:"xxx"})
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

	assert.Equal(t,bpluserr.HTTPErrorCode,400)

}
