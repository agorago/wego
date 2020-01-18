package fw

import (
	"context"
	"reflect"

	e "github.com/MenaEnergyVentures/bplus/internal/err"
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
	URL             string
	OpRequestMaker  PayloadMaker
	OpResponseMaker PayloadMaker
	HTTPMethod      string
	Params          []ParamDescriptor
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
		od.setupOperation()
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

func (od OperationDescriptor) setupOperation() {
	for _, op := range operations.Operations {
		op(od)
	}
}

// FindOperationDescriptor - find an op within a service
func FindOperationDescriptor(serviceName string, opName string) (OperationDescriptor, error) {
	sd := allServices[serviceName]
	if sd == nil {
		return OperationDescriptor{}, e.MakeBplusError(context.TODO(), e.ServiceNotFound, serviceName)
	}
	for _, od := range sd.Operations {
		if od.Name == opName {
			return od, nil
		}
	}
	return OperationDescriptor{}, e.MakeBplusError(context.TODO(), e.OperationNotFound, opName, serviceName)
}
