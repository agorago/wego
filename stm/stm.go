package stm

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"

	e "github.com/MenaEnergyVentures/bplus/internal/err"
)

// Register a set of states and restrict the transitions only to valid events for each of them

// ParamTypeMaker - a function that allows to make a parameter type
type ParamTypeMaker interface {
	MakeParam(Context context.Context) (interface{}, error)
}

type event struct {
	NewStateID string `json:"newState"`
	meta       map[string]interface{}
}

type state struct {
	Initial   bool // is this the initial state?
	Automatic bool // does this state compute the event itself?
	Events    map[string]event
	meta      map[string]interface{}
}

// Action - all actions must implement this interface
type Action interface {
	Process(context.Context, StateTransitionInfo) error
}

// AutomaticState - all auto states must implement this interface
// the return value for the Process method is a string that denotes
// the event
type AutomaticState interface {
	Process(context.Context, StateEntity) (string, error)
}

// StateTransitionInfo - the structure that will be passed to all the actions.
// It gives useful information about the state change that is currently happening
type StateTransitionInfo struct {
	OldState            string
	NewState            string
	Event               string
	Param               interface{}
	AffectedStateEntity StateEntity
	// AffectedStateEntity interface{}
}

// Stm - The definition of the Stm structure
type Stm struct {
	initialState  string
	states        map[string]state
	actionCatalog map[string]interface{}
}

// StateEntity - The State Entity against which the State machine operates
type StateEntity interface {
	GetState() string
	SetState(newstate string)
}

func (stm *Stm) populate(bytes []byte) {
	stm.states = make(map[string]state)
	if err := json.Unmarshal(bytes, &stm.states); err != nil {
		log.Fatal(err)
	}
}

// MakeStm - Make a new State machine with the filename that contains the State Transition Diagram
// and the action catalog
func MakeStm(filename string, actionCatalog map[string]interface{}) (*Stm, error) {
	stm := Stm{
		states:        make(map[string]state),
		actionCatalog: make(map[string]interface{}),
	}
	dat, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, e.MakeBplusError(context.TODO(), e.CannotReadFile,
			map[string]interface{}{
				"File": filename, "Error": err.Error()})
	}
	stm.populate(dat)
	stm.actionCatalog = actionCatalog
	return &stm, nil
}

// Process - called during run time to transition from one state to the next
func (stm Stm) Process(ctx context.Context, stateEntity StateEntity, event string, param interface{}) (StateEntity, error) {
	stateTransitionInfo := StateTransitionInfo{
		Param:               param,
		AffectedStateEntity: stateEntity,
		Event:               event,
		OldState:            stateEntity.GetState(),
	}

	if stateTransitionInfo.OldState == "" {
		// we need to initialize the entity for the first time.
		stateTransitionInfo.Event = InitialEvent
		stateTransitionInfo.NewState = stm.InitialState()
		return stm.doProcess(ctx, stateTransitionInfo)
	}
	stateInfo, exists := stm.states[stateTransitionInfo.OldState]
	if !exists {
		return stateEntity, e.MakeBplusError(ctx, e.InvalidState, map[string]interface{}{
			"State": stateTransitionInfo.OldState})
	}

	eventInfo, exists := stateInfo.Events[event]

	if !exists {
		return nil, e.MakeBplusError(ctx, e.InvalidEvent, map[string]interface{}{
			"Event": event, "State": stateTransitionInfo.OldState})
	}
	stateTransitionInfo.NewState = eventInfo.NewStateID
	return stm.doProcess(ctx, stateTransitionInfo)
}

func (stm Stm) doProcess(context context.Context, stateTransitionInfo StateTransitionInfo) (StateEntity, error) {
	preprocessor := stm.lookupAction(stateTransitionInfo.NewState, PreProcessorSuffix)
	postprocessor := stm.lookupAction(stateTransitionInfo.OldState, PostProcessorSuffix)
	transitionAction := stm.lookupAction(stateTransitionInfo.Event, TransitionActionSuffix)

	if postprocessor != nil {
		postprocessor.Process(context, stateTransitionInfo)
	}
	if transitionAction != nil {
		transitionAction.Process(context, stateTransitionInfo)
	}
	stateTransitionInfo.AffectedStateEntity.SetState(stateTransitionInfo.NewState)
	if preprocessor != nil {
		preprocessor.Process(context, stateTransitionInfo)
	}
	return stm.checkIfAutomaticState(context, stateTransitionInfo.AffectedStateEntity)
}

func (stm Stm) checkIfAutomaticState(ctx context.Context, stateEntity StateEntity) (StateEntity, error) {
	currentState := stateEntity.GetState()
	stateInfo, exists := stm.states[currentState]
	if !exists {
		return stateEntity, e.MakeBplusError(ctx, e.InvalidState, map[string]interface{}{
			"State": currentState})
	}
	if !stateInfo.Automatic {
		return stateEntity, nil
	}
	autoStateProcessor := stm.lookupAutoState(currentState, AutomaticStateSuffix)
	if autoStateProcessor == nil {
		return stateEntity, e.MakeBplusError(ctx, e.AutoStateNotConfigured, map[string]interface{}{
			"State": currentState})
	}
	event, err := autoStateProcessor.Process(ctx, stateEntity)
	if err != nil {
		return stateEntity, e.MakeBplusError(ctx, e.ErrorInAutoState, map[string]interface{}{
			"State": currentState, "Error": err.Error()})
	}
	return stm.Process(ctx, stateEntity, event, nil)
}

// STM Constants
const (
	Generic                = "generic"
	PreProcessorSuffix     = "PreProcessor"
	PostProcessorSuffix    = "PostProcessor"
	TransitionActionSuffix = "TransitionAction"
	ParamTypeMakerSuffix   = "ParamTypeMaker"
	AutomaticStateSuffix   = "AutoState"
	InitialEvent           = "InitialEvent" // Used for initializing the state entity. Call it the first event used to create the entity
)

// ParamTypeMaker - returns a param type maker for an event.
func (stm Stm) ParamTypeMaker(event string) ParamTypeMaker {
	f := stm.lookup(event, ParamTypeMakerSuffix)
	if f != nil {
		return f.(ParamTypeMaker)
	}
	return nil
}

func (stm Stm) lookupAction(prefix string, suffix string) Action {
	f := stm.lookup(prefix, suffix)
	if f != nil {
		return f.(Action)
	}
	return nil
}

func (stm Stm) lookupAutoState(prefix string, suffix string) AutomaticState {
	f := stm.lookup(prefix, suffix)
	if f != nil {
		return f.(AutomaticState)
	}
	return nil
}

func (stm Stm) lookup(prefix string, suffix string) interface{} {
	name := prefix + suffix
	action := stm.actionCatalog[name]
	if action == nil {
		name = Generic + suffix
		action = stm.actionCatalog[name]
	}
	return action
}

// InitialState - gives the initial state for the STM
func (stm *Stm) InitialState() string {
	if stm.initialState == "" {
		for stateID, stateInfo := range stm.states {
			if stateInfo.Initial {
				stm.initialState = stateID
				break
			}
		}
	}
	return stm.initialState
}
