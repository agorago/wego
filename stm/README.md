# State Transition Machine (STM)

Transactional entities exhibit a life cycle. Their lifecycle can be modeled using events. Examples of such 
entities are orders, payments, customer relationship issues etc. These entities have a workflow associated
with them. The workflow dictates the kind of events that these entities undergo. In most organizations, these
workflows are extremely dynamic. They can evolve over a period of time. 

Hence it is imperative to model the workflow as a first class entity. This allows the workflow to evolve 
as the organizations evolve. Besides,  an organization such as Intelligent B enters into contracts with 
myriad organizations across multiple geographies. This means that it is hard to evolve if the lifecycle 
of key entities are hard-coded to conform to specific workflows. 

The state transition machine is an artifact that separates the workflow from the core entity. Workflows 
can be modeled graphically and fed into a state machine using a JSON (or similar). The state machine ensures
that the entity changes state in conformance to the workflow. At various points of the workflow code can get 
called. This code can do specific actions.

This gives a huge deal of flexibility to the service that is responsible for the entity. The BPlus state 
transition machine is a simple but effective way of managing entity workflows. 

# Typical Life Cycle of a State Managed Transactional Entity 

These entities get created by an authorized user. 
The entities mutate their state when events happen to them.
If these events can be stored and replayed then the transactional entity can be replayed in full.
In that way, state entities offer a natural way of event sourcing
Event modeling of such entities makes them auditable. 

# Modeling State
State entities are modeled using a _state transition diagram_. (STD) State entity is a graph of states
as nodes connected by events which are edges connecting the nodes. An entity can traverse from one 
state to the next only by handling events that depict the edges. 

# The BPlus State Machine (STM) - code 
The BPlus state machine provides a way to implement a state machine by capturing the entire STD using a
JSON. 

In code, this is done by invoking the _makeStm_ method. 
```golang
  actionCatalog := map[string]interface{} 
  // populate action catalog.. see notes below
  fsm, err := stm.MakeStm("test.json",actionCatalog)
  // instantiate a state entity i.e. an entity that returns  GetState and SetState methods
  // Example: if there an exists a type called Order which implements both GetState and SetState
  type Order struct{ State string}
  func (ord *Order)GetState() string{
     return ord.State
   }
   func (ord *Order) SetState(s string){
       ord.State = s
   }
  order := &Order{} 
   // pass eventId and eventParam for 
  o = fsm.Process(ctx, &order, eventId, eventParam).(*Order)
  // o is a pointer to the mutated version of the Order
```
As seen above, the first step is to create an STM and then use it with an entity such as __order__ that
implements the methods in the __StateEntity__ interface defined in the _stm_ package. 

Let us look at the test case for more details. 

## stm_test

_stm_test_ uses a test.json that defines four states 
created -> confirmed -> fulfilled -> closed 

for a testOrder. The testOrder supports the following events to progress from one state to the next

*confirm* - to progress from __created__ to __confirmed__ state
*fulfil* - to progress from __confirmed__ to __fulfilled__ state
*close* - to progress from __fulfilled__ to __closed__ state

Additionally there is a state called __cancelled__ to which the testOrder can progress from either the
__created__ or __confirmed__ states. However, progressing from __confirmed__ to __cancelled__ state 
incurs a penalty of 25 currency units. This makes the order non closeable from the cancelled state. 

This process is represented in the _std_ using the state __check-if-closeable__ which checks if the
charges are > 0 to decide to allow the transition.

### testOrder struct
*testOrder* is the name of the state entity in __stm_test__. This entity needs to implement the method
in the StateEntity interface (obviously!) to qualify as a state entity. Besides that, it has specific
fields with various order attributes. The field __CalledFunctions__ keeps track of all the functions
that have been called in an array. This is useful for testing and allows assertions to be made. 

# Actions

There are four types of actions supported by the Bplus STM:
* Transition Actions - an action that is called for every event
* Preprocessor Action - which is called _before_ a state is entered
* Postprocessor Action - which is called _after_ a state is exited 
* AutomaticState Action - which is called to determine the event in an automatic state

## stm_test - actions

In __stm_test__ there exists a preprocessor and a postprocessor.
Additionally all actions are mapped to a struct called _action_ which merely adds an entry to the 
_CalledFunctions_ array.
cancelledConfirmedAction is special since it levies charges on the testOrder. So this action is mapped
only to the cancelledConfirmedAction event. 

### Automatic States

The __check-if-closeable__ state is an automatic state. This means that the __stm__ requires a special
action to allow the computation of events given the state of the entity. 

The __checkIfCloseable__ type takes care of this requirement. It uses the charges as the basis to determine
if the __testOrder__ is eligible to be closed or not. 

## The test cases
The test cases cover the gamut of functionality as the entity transitions from one state to the other. 


 