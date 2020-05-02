package stm

import (
	"context"
	"encoding/json"
	e "gitlab.intelligentb.com/devops/bplus/internal/err"
	"io/ioutil"
)

// Register a set of States and restrict the transitions only to valid events for each of them

// ParamTypeMaker - a function that allows to make a parameter type
type ParamTypeMaker interface {
	MakeParam(Context context.Context) (interface{}, error)
}

type Event struct {
	NewStateID string `json:"newState"`
	meta       map[string]interface{}
}

type State struct {
	Initial   bool // is this the initial State?
	Automatic bool // does this State compute the Event itself?
	Events    map[string]Event
	meta      map[string]interface{}
}

// Action - all actions must implement this interface
type Action interface {
	Process(context.Context, StateTransitionInfo) error
}

// AutomaticState - all auto States must implement this interface
// the return value for the Process method is a string that denotes
// the Event
type AutomaticState interface {
	Process(context.Context, StateEntity) (string, error)
}

// StateTransitionInfo - the structure that will be passed to all the actions.
// It gives useful information about the State change that is currently happening
type StateTransitionInfo struct {
	OldState            string
	NewState            string
	Event               string
	Param               interface{}
	AffectedStateEntity StateEntity
}

// Stm - The definition of the Stm structure
type Stm struct {
	initialState  string
	States        map[string]State
	actionCatalog map[string]interface{}
}

// StateEntity - The State Entity against which the State machine operates
type StateEntity interface {
	GetState() string
	SetState(newstate string)
}

func (stm *Stm) populate(bytes []byte) error {
	stm.States = make(map[string]State)
	if err := json.Unmarshal(bytes, &stm.States); err != nil {
		return e.MakeBplusError(context.TODO(), e.UnparseableFile, nil)
	}
	return nil
}

// MakeStm - Make a new State machine with the filename that contains the State Transition Diagram
// and the action catalog
func MakeStm(filename string, actionCatalog map[string]interface{}) (*Stm, error) {
	stm := Stm{
		States:        make(map[string]State),
		actionCatalog: make(map[string]interface{}),
	}
	dat, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, e.MakeBplusError(context.TODO(), e.CannotReadFile,
			map[string]interface{}{
				"File": filename, "Error": err.Error()})
	}
	err = stm.populate(dat)
	if err != nil {
		return nil, err
	}
	stm.actionCatalog = actionCatalog
	return &stm, nil
}

// Process - called during run time to transition from one State to the next
func (stm Stm) Process(ctx context.Context, stateEntity StateEntity, event string, param interface{}) (StateEntity, error) {
	stateTransitionInfo := StateTransitionInfo{
		Param:               param,
		AffectedStateEntity: stateEntity,
		Event:               event,
		OldState:            stateEntity.GetState(),
	}

	if stateTransitionInfo.OldState == "" || event == InitialEvent {
		// we need to initialize the entity for the first time.
		stateTransitionInfo.Event = InitialEvent
		stateTransitionInfo.NewState = stm.InitialState()
		return stm.doProcess(ctx, stateTransitionInfo)
	}
	stateInfo, exists := stm.States[stateTransitionInfo.OldState]
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
	return stm.checkIfAutomaticState(context, stateTransitionInfo)
}

func (stm Stm) checkIfAutomaticState(ctx context.Context, info StateTransitionInfo) (StateEntity, error) {
	stateEntity := info.AffectedStateEntity
	currentState := stateEntity.GetState()
	stateInfo, exists := stm.States[currentState]
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
	return stm.Process(ctx, stateEntity, event, info.Param)
}

// STM Constants
const (
	Generic                = "generic"
	PreProcessorSuffix     = "PreProcessor"
	PostProcessorSuffix    = "PostProcessor"
	TransitionActionSuffix = "TransitionAction"
	ParamTypeMakerSuffix   = "ParamTypeMaker"
	AutomaticStateSuffix   = "AutoState"
	InitialEvent           = "InitialEvent" // Used for initializing the State entity. Call it the first Event used to create the entity
)

// ParamTypeMaker - returns a param type maker for an Event.
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

// InitialState - gives the initial State for the STM
func (stm *Stm) InitialState() string {
	if stm.initialState == "" {
		for stateID, stateInfo := range stm.States {
			if stateInfo.Initial {
				stm.initialState = stateID
				break
			}
		}
	}
	return stm.initialState
}
