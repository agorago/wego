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
In that way, state entities 

Event modeling of such entities makes them auditable. Besides, as workflows evolve 