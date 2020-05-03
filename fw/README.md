# The WEGO Framework package
WEGO is at its core, a framework that invokes a service (more precisely it invokes an operation in a
service). For this to happen, it needs the following:
1. It needs to know the services that exist.
2. It needs to know their dependencies as defined in a makeService()
3. It needs to know about the interfaces exposed by the service (so that it can potentially mock the
service for testing purposes).
4. It needs to know the transports that the service decides to expose so that WEGO can provide the 
transport.
5. It needs to know what middlewares that the service wants to have before it. The middlewares can be
at the client side or the server side. 


The operation is invoked in 
response to a request at a transport layer. The core framework package provides the following features:
0. Instantiate the framework along with a DI bean factory. DI in GO has thus far been provided by 
wiring frameworks that There is no core bean factory within GO. Hence
it becomes difficult to do things like Auto-wiring in the style of th
1. The ability for a service to register itself with the framework.
2. The ability for a transport such as HTTP to register an operation in order to expose it to the 
outside world.
3. Definition and registration of middlewares. These middlewares are invoked whenever the service is 
invoked. 
4. Base interfaces are defined in this package. These interfaces facilitate service registration.
5. This package also implements the notion of a Middleware Chain which allows the set up of chains of
middleware.These chains invoke a set of middlewares and finally culminate in the service. Please see
the internal/mw package for more information about middleware chains.
 
# Service Registration

A service can have one or more multiple methods called Operations. Each operation can in turn 
support parameters. 

All BPlus services register here. They can register in two modes:
1. Client Mode - the actual service information is registered. However the service is not exposed
in this mode using any transport. e
2. Server Mode -  the service is not only registered. It is also exposed via a transport layer such as HTTP. 

Client mode is useful if it is intended to invoke the service using the Proxy framework that BPlus supports.
However in this mode,  transports are not supported. Consequently, the service will not be available for 
remote invocation via HTTP or other means.  Services can only be registered in server mode in one executable.
In other executables, the service is registered only as in client mode.

# How do services register themselves?

A service is registered using a model object called _Service Descriptor_ that describes the service. It contains
Name and a description. Besides that it contains an array of operations that describe the operations that 
the service supports. Each operation is described using OperationDescriptor.

If the service is registered in server mode, then a reference to the actual service is expected to be 
populated in the  the service descriptor. If service reference is null then it is assumed that the service 
is registered in client mode. Otherwise, the service reference is required since the transport needs to 
call the actual service when a request for service is received by the transport end point.

Transports are discussed a little later.

# Operation Descriptor

Operation descriptors describe the operation that belongs to the service. The operation would be typically
exposed by the transport. Since an operation is accessible from outside, it will expose a URL. Operations
are invoked by BPlus when the transport is accessed. (for example when someone invokes the HTTP URL)

Operation descriptor gives details about the operation such as the following:
name - used to actually invoke the operation from the service reference in ServiceDescriptor
description, Request Description, Response Description - used to describe the service in generated swagger 
documentations
URL - used to expose the operation via HTTP
OpRequestMaker , OpResponseMaker - These functions generate empty request and response objects. useful to 
populate the objects from serialized streams of data (such as JSON streams)
OpMiddleware, ProxyMiddleware - the middleware functions that will be invoked when the operation is invoked

# Transports 

Transport provides the end point that initiates all service invocations. Example: an HTTP end point which 
when invoked, invokes the operation. There are properties in Operation Descriptor which can be used to 
configure the transport end points. For example, the HTTP end point is configured using the URL property.

Transports are extensions to the core framework and are implemented using an OperationRegistration. 
Transports register themselves with the core framework by using RegisterOperation. The transports will then
be provided an opportunity to register every operation that is registered with BPlus. This happens only 
when the service is registed in server-mode. 

# Param Descriptor

ParamDescriptor allows the registration of each parameter that needs to be passed to the operation on 
invocation.
An operation can have the following types of parameters as determined by ParamOrigin:
Context - a mandatory first parameter for every operation. It should always be with name ctx. This
contains the entire context of the request.
HEADER - a parameter that is derived from a header in the request. In HTTP, HEADER parameters are 
passed as a HTTP HEADER attribute. Alternately, HTTP query params and PATH params are also 
available as header params.
PAYLOAD - the payload parameter is passed as part of the body of the request. (for example HTTP 
request body) Payload will be de-serialized from a resource stream such as JSON.

# Middleware Chain

The framework package defines a Middleware Chain. This implements the chain of responsibility pattern.
Please see middleware_test.go to get an idea about the usage of this class.
