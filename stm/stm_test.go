package stm_test

import (
	"context"
	"fmt"
	"github.com/magiconair/properties/assert"
	"gitlab.intelligentb.com/devops/bplus/err"
	"gitlab.intelligentb.com/devops/bplus/i18n"
	bpluse "gitlab.intelligentb.com/devops/bplus/internal/err"
	"gitlab.intelligentb.com/devops/bplus/stm"
	"os"
	"testing"
)
// A test case where simple order details are captured.
// The order progresses from one State to the next sequentially
// Created -> Confirmed -> Fulfilled -> Closed
// Order can be cancelled in Created and Confirmed States
// Cancellation in confirmed State incurs penalties
type testOrder struct{
	ID string
	State string
	Items []string // items of the order
	Price float32 // order price
	Charges float32 // charges incurred if the order is cancelled
	CalledFunctions []string // a list of all the functions called. Used to
	// establish if a function is indeed called by the STM. It
	// is also useful to establish the order of calling
}
func (ts *testOrder) GetState()string{
	return ts.State
}

func (ts *testOrder) SetState(s string){
	ts.State = s
}

func createMachine(filename string, m map[string]interface{}) *stm.Stm{
	stm,err := stm.MakeStm(filename,m)
	if err != nil {
		fmt.Fprintf(os.Stderr,"Cannot create stm. Error = %s\n",err.Error())
		return nil
	}
	return stm
}

type checkIfCloseable struct{}
func (checkIfCloseable)Process(ctx context.Context, se stm.StateEntity) (string, error){
	to,ok := se.(*testOrder)
	if !ok{
		return "",err.MakeErr(ctx,err.Error,123,"STM-TEST: Cannot cast to the correct type",nil)
	}
	if to.Charges > 0 {
		return "no",nil
	}
	return "yes",nil
}

type cancelConfirmed struct{}
func (cancelConfirmed) Process(ctx context.Context, info stm.StateTransitionInfo) error {
	to,ok := info.AffectedStateEntity.(*testOrder)
	if !ok{
		return err.MakeErr(ctx,err.Error,123,"STM-TEST: Cannot cast to the correct type",nil)
	}
	// levy a charge of 25 units on confirmed orders
	to.Charges = 25.0
	return add(info,"pre"+info.NewState)
}
type preProcessor struct{}
func (preProcessor) Process(_ context.Context, info stm.StateTransitionInfo) error {
	return add(info,"pre"+info.NewState)
}
type postProcessor struct{}
func (postProcessor) Process(_ context.Context, info stm.StateTransitionInfo) error {
	return add(info,"post"+info.OldState)
}
type action struct{}
func (action) Process(_ context.Context, info stm.StateTransitionInfo) error {
	return add(info,info.Event)
}

func add(info stm.StateTransitionInfo,s string)error{
	to,ok := info.AffectedStateEntity.(*testOrder)
	if !ok{
		return err.MakeErr(context.TODO(),err.Error,123,"STM-TEST: Cannot cast to the correct type",nil)
	}
	to.CalledFunctions = append(to.CalledFunctions,s)
	return nil
}

var fsm *stm.Stm
func initialize(t *testing.T) {
	act := action{}
	i18n.InitConfig("../configs")
	m := map[string]interface{}{
		stm.Generic + stm.PreProcessorSuffix:          preProcessor{},
		stm.Generic + stm.PostProcessorSuffix:         postProcessor{},
		stm.InitialEvent + stm.TransitionActionSuffix: act,
		"confirm" + stm.TransitionActionSuffix:        act,
		"fulfil" + stm.TransitionActionSuffix:         act,
		"close" + stm.TransitionActionSuffix:          act,
		"cancelCreatedOrder" + stm.TransitionActionSuffix:          act,
		"cancelConfirmedOrder" + stm.TransitionActionSuffix:          cancelConfirmed{},
		"check-if-closeable" + stm.AutomaticStateSuffix:          checkIfCloseable{},
	}
	fsm = createMachine("test.json",m)

	if fsm == nil {
		t.Fail()
	}
}

func TestStateProgression(t *testing.T) {
	initialize(t)
	to := testOrder{State: ""}
	ts := invoke(t,&to,"",nil).(*testOrder)
	assert.Equal(t,ts.GetState(),"created","STM-TEST: State change not successful")
	ts = invoke(t,ts,"confirm",nil).(*testOrder)
	assert.Equal(t,ts.GetState(),"confirmed","STM-TEST: State change not successful")
	ts = invoke(t,ts,"fulfil",nil).(*testOrder)
	assert.Equal(t,ts.GetState(),"fulfilled","STM-TEST: State change not successful")
	ts = invoke(t,ts,"close",nil).(*testOrder)
	assert.Equal(t,ts.GetState(),"closed","STM-TEST: State change not successful")

	assert.Equal(t,[]string{
		"post","InitialEvent","precreated", // for the initial create there wont be an old State
		"postcreated","confirm","preconfirmed",
		"postconfirmed","fulfil","prefulfilled",
		"postfulfilled","close","preclosed",
		},

		ts.CalledFunctions,"STM_TEST:functions not called in expected order")
}

func invoke(t *testing.T,to *testOrder,event string, param interface{}) stm.StateEntity{
	ts,err := fsm.Process(context.TODO(),to,event,param)
	if err != nil {
		fmt.Printf("Failed to invoke. Error = %s\n",err.Error())
		t.Fail()
	}
	return ts
}


func TestCloseable(t *testing.T){
	initialize(t)
	to := testOrder{State: ""}
	ts := invoke(t,&to,"",nil).(*testOrder)
	assert.Equal(t,ts.GetState(),"created","STM-TEST: State change not successful")
	ts = invoke(t,ts,"cancelCreatedOrder",nil).(*testOrder)
	assert.Equal(t,ts.GetState(),"cancelled","STM-TEST: State change not successful")
	ts = invoke(t,ts,"closeCancelledOrder",nil).(*testOrder)
	assert.Equal(t,ts.GetState(),"closed","STM-TEST: State change not successful")
}

func TestNonCloseable(t *testing.T){
	initialize(t)
	to := testOrder{State: ""}
	ts := invoke(t,&to,"",nil).(*testOrder)
	assert.Equal(t,ts.GetState(),"created","STM-TEST: State change not successful")
	ts = invoke(t,ts,"confirm",nil).(*testOrder)
	assert.Equal(t,ts.GetState(),"confirmed","STM-TEST: State change not successful")
	ts = invoke(t,ts,"cancelConfirmedOrder",nil).(*testOrder)
	assert.Equal(t,ts.GetState(),"cancelled","STM-TEST: State change not successful")
	ts = invoke(t,ts,"closeCancelledOrder",nil).(*testOrder)
	// verify that the order remains in cancelled State since there are pending charges
	assert.Equal(t,ts.GetState(),"cancelled","STM-TEST: State change not successful")
}

func TestInvalidEvent(t *testing.T){
	initialize(t)
	to := testOrder{State: ""}
	ts := invoke(t,&to,"",nil).(*testOrder)
	assert.Equal(t,ts.GetState(),"created","STM-TEST: State change not successful")
	_,er := fsm.Process(context.TODO(),ts,"close",nil)
	e,ok := er.(err.BPlusError)
	if !ok {
		fmt.Fprintf(os.Stderr,"Error not castable to BPlusError!\n")
		t.Fail()
	}
	code := int(bpluse.InvalidEvent)
	assert.Equal(t,code,e.ErrorCode)
}

func TestInvalidState(t *testing.T){
	initialize(t)
	to := testOrder{State: ""}
	ts := invoke(t,&to,"",nil).(*testOrder)
	assert.Equal(t,ts.GetState(),"created","STM-TEST: State change not successful")
	ts.SetState("INVALID") // set it to an invalid State
	_,er := fsm.Process(context.TODO(),ts,"confirm",nil)
	e,ok := er.(err.BPlusError)
	if !ok {
		fmt.Fprintf(os.Stderr,"Error not castable to BPlusError!\n")
		t.Fail()
	}
	code := int(bpluse.InvalidState)
	assert.Equal(t,code,e.ErrorCode)
}

func TestInvalidFilename(t *testing.T){
	_,er := stm.MakeStm("test1.json",nil)
	e,ok := er.(err.BPlusError)
	if !ok {
		fmt.Fprintf(os.Stderr,"Error not castable to BPlusError!\n")
		t.Fail()
	}
	code :=  int(bpluse.CannotReadFile)
	assert.Equal(t,code,e.ErrorCode)
}

func TestUnparseableFilename(t *testing.T){
	_,er := stm.MakeStm("test-unparseable.json",nil)
	e,ok := er.(err.BPlusError)
	if !ok {
		fmt.Fprintf(os.Stderr,"Error not castable to BPlusError!\n")
		t.Fail()
	}
	code :=  int(bpluse.UnparseableFile)
	assert.Equal(t,code,e.ErrorCode)
}



