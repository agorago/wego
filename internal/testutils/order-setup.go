package testutils

import (
	"context"
	"github.com/agorago/wego/fw"
	wegohttp "github.com/agorago/wego/http"
	"github.com/agorago/wego/log"
	"github.com/agorago/wego/stateentity"
	"github.com/agorago/wego/stm"
	"path/filepath"
)

// Sets up an order object to test state entity implementation
// first define the entity
type Order struct {
	ID             string
	State          string
	VisitedActions []string
	AllMessages    []string
}

func (ord *Order) GetState() string {
	return ord.State
}

func (ord *Order) SetState(newState string) {
	ord.State = newState
}

// Next set up a dummy repo implemented as a map
type OrderRepo struct {
	entities map[string]*Order
}

func (or *OrderRepo) Create(se stm.StateEntity) (stm.StateEntity, error) {
	ent := se.(*Order)
	ent.ID = "first"
	repo.entities[ent.ID] = ent
	return ent, nil
}

func (or *OrderRepo) Retrieve(ID string) (stm.StateEntity, error) {
	return repo.entities[ID], nil
}

// repo - the repository used for persisting the data into the store
var repo = &OrderRepo{
	entities: make(map[string]*Order),
}

// Next set up some added functions for the entity
func makeModel(_ context.Context) (stm.StateEntity, error) {
	return &Order{}, nil
}

// stmChooser - makes the STM - this allows flexibility to choose the correct STD based on context
func stmChooser(_ context.Context) (*stm.Stm, error) {
	return thisStm, nil
}

// all actions associated with the order state entity
var actionCatalog = make(map[string]interface{})
var thisStm *stm.Stm

// next make an STM with a workflow for the entity.
func makeSTM() (*stm.Stm, error) {
	path, _ := filepath.Abs("test-configs" + "/workflow/" + "order-std.json")
	log.Infof(context.Background(), " make stm action catalog is %v\n", actionCatalog)
	return stm.MakeStm(path, actionCatalog)
}

// make a pre-processor to enable persistence into the repo

func registerPreprocessor() {
	actionCatalog[stm.Generic+stm.PreProcessorSuffix] = PreProcessor{}
}

type PreProcessor struct{}

func (PreProcessor) Process(ctx context.Context, info stm.StateTransitionInfo) error {
	// persist the entity into the store
	repo.Create(info.AffectedStateEntity)
	return nil
}

// finally register a bunch of actions
func registerActions() {
	registerAction("startPaymentProcess")
	registerAction("cancelOrder")
	registerAction("receivePayment")
	registerAction("yes")
	registerAction("no")
	registerAction("closeOrder")
}

func registerAction(name string) {
	actionCatalog[name+stm.TransitionActionSuffix] = genericAction{Name: name}
	actionCatalog[name+stm.ParamTypeMakerSuffix] = actionParamTypeMaker{}
}

type ActionParam struct {
	Message string
}

type genericAction struct {
	Name string
}
type actionParamTypeMaker struct{}

func (ga genericAction) Process(_ context.Context, info stm.StateTransitionInfo) error {
	p := info.Param.(*ActionParam)
	order := info.AffectedStateEntity.(*Order)
	// Add this action to the list of visited actions
	order.VisitedActions = append(order.VisitedActions, ga.Name)
	order.AllMessages = append(order.AllMessages, p.Message)
	return nil
}

func (actionParamTypeMaker) MakeParam(ctx context.Context) (interface{}, error) {
	return &ActionParam{}, nil
}

func SetupOrder() (fw.RegistrationService,stateentity.StateEntityProxy){
	var err error
	registerPreprocessor()
	registerActions()
	thisStm, err = makeSTM()
	if err != nil {
		log.Error(context.Background(), "FATAL: Cannot create STM")
	}

	var registration = stateentity.SubTypeRegistration{
		Name:                    "order",
		StateEntitySubTypeMaker: makeModel,
		ActionCatalog:           actionCatalog,
		URLPrefix:               "order",
		StateEntityRepo:         repo,
		OrderSTMChooser:         stmChooser,
	}

	rs := fw.MakeRegistrationService()
	proxyService := wegohttp.MakeProxyService(rs)
	proxy := stateentity.MakeStateEntityProxy(proxyService)
	sers := stateentity.MakeStateEntityRegistrationService(rs)
	sers.RegisterSubType(registration)
	return rs,proxy
}
