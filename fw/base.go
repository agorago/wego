package fw

import (
	"context"
	"reflect"

	e "gitlab.intelligentb.com/devops/bplus/internal/err"
)

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
	Name            string
	Description     string
	Service         ServiceDescriptor
	RequestDescription string // the description of the request - useful for swagger
	ResponseDescription string // the description of the response
	URL             string
	OpRequestMaker  PayloadMaker
	OpResponseMaker PayloadMaker
	RequestType     interface{}
	ResponseType    interface{}
	OpMiddleware    []Middleware // specific middleware required by this operation. These will be
	// invoked before the operation is invoked by BPlus
	HTTPMethod string
	Params     []ParamDescriptor
}

// ServiceDescriptor - the root of the service registration.
type ServiceDescriptor struct {
	Name            string
	Description     string
	ServiceToInvoke interface{}
	Operations      []OperationDescriptor
}

var allServices = make(map[string]*ServiceDescriptor)

// RegisterService - call this function to register a service with an ID
func RegisterService(ID string, sd ServiceDescriptor) {
	allServices[ID] = &sd
	sd.setupService()
}

func (sd *ServiceDescriptor) setupService() {
	for _, od := range sd.Operations {
		od.Service = *sd
		od.setupOperation(sd.ServiceToInvoke == nil) // if the service to invoke is null then
		// all transports are disabled for this service registration. It is only used to
		// update the local registry of available services. (useful to invoke using proxy etc.)
	}
}

// OperationRegistration - any function that registers the operations and accomplishes
// something
type OperationRegistration func(OperationDescriptor)

var operations struct {
	Operations []OperationRegistration
}

// RegisterOperations - Register extension to BPlus with definied operations here
func RegisterOperations(or OperationRegistration) {
	operations.Operations = append(operations.Operations, or)
}

func (od OperationDescriptor) setupOperation(disableTransport bool) {
	if disableTransport{
		return
	}
	for _, op := range operations.Operations {
		op(od)
	}
}

// FindOperationDescriptor - find an op within a service
func FindOperationDescriptor(serviceName string, opName string) (OperationDescriptor, error) {
	sd := allServices[serviceName]
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
func FindServiceDescriptor(serviceName string)(*ServiceDescriptor,error){
	sd := allServices[serviceName]
	if sd == nil {
		return nil, e.MakeBplusError(context.TODO(), e.ServiceNotFound, map[string]interface{}{
			"Service": serviceName})
	}
	return sd,nil
}
