package fw

import (
	"context"
	"reflect"

	e "github.com/agorago/wego/internal/err"
)
// The chief framework for WEGO starts here.
// This provides the registration functionality for new commands.

// PayloadMaker - makes a payload of a certain type
type PayloadMaker func(context.Context) (interface{}, error)

// Origin - defines a param type as int with a defined set of values
type Origin int

// Param type enumeration
const (
	HEADER  Origin = iota + 1 // when the param is populated from headers in the context object
	PAYLOAD                   // when the param is populated from the payload of the request
	CONTEXT                   // when the context is the object required to be populated
)

// ParamDescriptor - which describes the param and denotes its type
type ParamDescriptor struct {
	Name        string
	ParamOrigin Origin
	Description string
	ParamKind   reflect.Kind
}

// OperationDescriptor - to describe the properties of the operation
type OperationDescriptor struct {
	Name                string
	Description         string
	Service             ServiceDescriptor
	RequestDescription  string // the description of the request - useful for swagger
	ResponseDescription string // the description of the response
	URL                 string
	OpRequestMaker      PayloadMaker
	OpResponseMaker     PayloadMaker
	RequestType         interface{}
	ResponseType        interface{}
	OpMiddleware        []Middleware // specific middleware required by this operation. These will be
	// invoked before the operation is invoked by BPlus
	HTTPMethod      string
	ProxyMiddleware []Middleware // this is only used on the proxy side and not on the server side.
	Params          []ParamDescriptor
}

// ServiceDescriptor - the root of the service registration.
type ServiceDescriptor struct {
	Name            string      // The name of the service.
	Description     string      // description of the service
	ServiceToInvoke interface{} // a pointer to the actual service to invoke
	Operations      []OperationDescriptor // All transports exposed by the service
}

type RegistrationService interface{
	RegisterService(ID string, sd ServiceDescriptor)
	FindOperationDescriptor(serviceName string, opName string) (OperationDescriptor, error)
	FindServiceDescriptor(serviceName string) (*ServiceDescriptor, error)
	RegisterTransport(id string,transport Transport)
}

type RegistrationServiceImpl struct{
	AllServices map[string]*ServiceDescriptor
	AllTransports map[string]Transport
}

func MakeRegistrationService() RegistrationService{
	return &RegistrationServiceImpl{
		 make(map[string]*ServiceDescriptor),
		 make(map[string]Transport),
	}
}

// RegisterService - call this function to register a service with an ID
func (rsimpl *RegistrationServiceImpl)RegisterService(ID string, sd ServiceDescriptor) {
	rsimpl.AllServices[ID] = &sd
	rsimpl.setupService(&sd)
}

func (rs RegistrationServiceImpl) setupService(sd *ServiceDescriptor) {
	for _, od := range sd.Operations {
		od.Service = *sd
		rs.setupOperation(od,(sd.ServiceToInvoke == nil))
		// if the service to invoke is null then
		// all transports are disabled for this service registration. It is only used to
		// update the local registry of available services. (useful to invoke using proxy etc.)
	}
}

// Transport - any function that exposes a transport at the operation level
type Transport func(OperationDescriptor)

// RegisterTransport - Register a transport as an extension to WEGO
func (rs RegistrationServiceImpl) RegisterTransport(id string, transport Transport) {
	rs.AllTransports[id] = transport
}

func (rs RegistrationServiceImpl) setupOperation(od OperationDescriptor,disableTransport bool) {
	if disableTransport {
		return
	}
	for _, transport := range rs.AllTransports {
		transport(od)
	}
}

// FindOperationDescriptor - find an op within a service
func (rsimpl *RegistrationServiceImpl)FindOperationDescriptor(serviceName string, opName string) (OperationDescriptor, error) {
	sd := rsimpl.AllServices[serviceName]
	if sd == nil {
		return OperationDescriptor{}, e.MakeBplusError(context.TODO(), e.ServiceNotFound, map[string]interface{}{
			"Service": serviceName})
	}
	for _, od := range sd.Operations {
		if od.Name == opName {
			od.Service = *sd
			return od, nil
		}
	}
	return OperationDescriptor{}, e.MakeBplusError(context.TODO(), e.OperationNotFound, map[string]interface{}{
		"Operation": opName, "Service": serviceName})
}

// FindServiceDescriptor - find a service descriptor that has been registered with the name
func (rsimpl *RegistrationServiceImpl)FindServiceDescriptor(serviceName string) (*ServiceDescriptor, error) {
	sd := rsimpl.AllServices[serviceName]
	if sd == nil {
		return nil, e.MakeBplusError(context.TODO(), e.ServiceNotFound, map[string]interface{}{
			"Service": serviceName})
	}
	return sd, nil
}
