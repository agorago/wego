package stateentity_test

import (
	"context"
	"fmt"
	"github.com/agorago/wego"
	"github.com/agorago/wego/config"
	"github.com/agorago/wego/internal/testutils"
	"github.com/agorago/wego/stateentity"
	"github.com/magiconair/properties/assert"
	"log"
	"net/http"
	"os"
	"testing"
)

func TestRegisterSubType(t *testing.T) {
	os.Setenv("BPLUS.PORT", "5000")
	commandCatalog := testutils.SetupOrder()
	httpHandler := commandCatalog.Command(wego.WegoHTTPHandler).(http.Handler)
	go func() {
		a := ":" + config.Value("bplus.port")
		log.Printf("Starting server at address %s\n", a)
		http.ListenAndServe(a, httpHandler)
	}()
	sep := commandCatalog.Command(wego.StateEntityProxy).(stateentity.StateEntityProxy)
	ctx := context.TODO()
	order := &testutils.Order{}
	o, e := sep.Create(ctx, "order", order)
	order, e = checkIfValid(t, o, e)
	if e != nil {
		return
	}
	assert.Equal(t, order.ID, "first")
	assert.Equal(t, order.State, "created")

	var expectedActions, expectedMessages []string
	ap := testutils.ActionParam{Message: "start-order"}
	expectedActions = append(expectedActions, "startPaymentProcess")
	expectedMessages = append(expectedMessages, "start-order")
	o, e = sep.Process(ctx, "order", order.ID, "startPaymentProcess", ap)
	order, e = checkIfValid(t, o, e)
	if e != nil {
		return
	}
	assert.Equal(t, order.State, "pending-payment")
	assert.Equal(t, order.VisitedActions, expectedActions)
	assert.Equal(t, order.AllMessages, expectedMessages)

	expectedActions = append(expectedActions, "receivePayment")
	expectedMessages = append(expectedMessages, "receive-payment")
	ap = testutils.ActionParam{Message: "receive-payment"}
	o, e = sep.Process(ctx, "order", order.ID, "receivePayment", ap)
	order, e = checkIfValid(t, o, e)
	if e != nil {
		return
	}
	assert.Equal(t, order.State, "completed")
	assert.Equal(t, order.VisitedActions, expectedActions)
	assert.Equal(t, order.AllMessages, expectedMessages)

	expectedActions = append(expectedActions, "closeOrder")
	expectedMessages = append(expectedMessages, "close-order")
	ap = testutils.ActionParam{Message: "close-order"}
	o, e = sep.Process(ctx, "order", order.ID, "closeOrder", ap)
	order, e = checkIfValid(t, o, e)
	if e != nil {
		return
	}
	assert.Equal(t, order.State, "closed")
	assert.Equal(t, order.VisitedActions, expectedActions)
	assert.Equal(t, order.AllMessages, expectedMessages)
}

func checkIfValid(t *testing.T, o interface{}, e error) (*testutils.Order, error) {
	if e != nil {
		log.Printf("Error in service invocation. Error = %s\n", e.Error())
		t.Fail()
		return nil, e
	}
	retOrder, ok := o.(*testutils.Order)
	if !ok {
		log.Printf("Error in casting the output to type order. output =  %#v\n", 0)
		t.Fail()
		return nil, fmt.Errorf("error in casting")
	}
	return retOrder, nil
}
