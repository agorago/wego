# State Entity package

The stateentity package in BPlus helps to expose (and also consume) state entities via HTTP.

The state entity in question is exposed via HTTP when it is registered. The http methods are hooked up to
the state machine and it then becomes possible to control the lifecycle of the state entity via HTTP.

The package provides the ability  for state entities to register themselves along with a set of actions
that need to act on them. At the time of registration the state entity methods are exposed via Bplus.
This allows BPlus to treat the state entity service like any other Bplus service. Hence all the 
facilities of interception (both client side and server side) as well as other standard Bplus services
are available.