# The WeGO Middlewares & Entry Points

This package contains the middlewares and entry points for WeGO. 


# WeGO Entry Point

An entry point is the point of entry to invoke WeGO services or proxies. There are two entry points - one for 
the server side and the other one for invoking a proxy. 

The entry point is responsible to set up the other middlewares and start the [middleware chain.](../../fw/README.md)
The entry point then returns the result of the invocation along with any errors encountered. 

There are two entry points here. One is the server side entrypoint.  Simply called as entry point. The 
other one is the entry point  in client side called proxy entry point. 

# The terminators

There are two terminating middlewares in this package. These are set up to be the last of the middlewares by
the entry points. Terminators will not continue the chain. Instead, they invoke the service either on the 
client or the server side as discussed below.

The server side entry point sets up service-invoker as the last 
of the middlewares. Service Invoker's responsibility is to invoke the actual service with the correct 
parameters as set up in OperationDescriptor and ParamDescriptor. The arguments are gathered from the 
context object and are passed to the service. The appropriate method of the service is called with the
parameters. Since it internally uses GO reflection, this can panic if the service name or the params
are mis-configured. Hence the service invoker resolves this panic and spits out the appropriate error
message

The proxy entry point sets up http-invoker as the terminator. This terminator will utilize the information
in the ServiceDescriptor to construct the HTTP request. It then uses a configuration property called
serviceName + "host_port" to derive the host name for the service. It then issues a request to the host
name. The response is then converted to a GO response. The errors are also appropriately re-created.

# Other Middlewares

In addition to the entry points and the terminators, this package also contains other middlewares which are 
discussed below. 

## Decoder

Decoder decodes the JSON request and converts it into a GOLANG struct of the appropriate type. 
__OperationDescriptor.OpRequestMaker__ is used to construct the request payload object. JSON parser is 
used to fill it up from the JSON string. The resulting struct is put as a payload into the context 
object using: 
```go
wegocontext.SetPayload()
```

The payload then becomes available for other middlewares and terminators. This is the first of the 
server side middlewares since the payload is needed by other middlewares

## V10 Validator

The V10 validator uses v10 to validate the payload. If the payload does not validate correctly, then 
an error is returned and the chain is discontinued. 
