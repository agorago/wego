package stateentity

import (
	"context"

	e "github.com/MenaEnergyVentures/bplus/internal/err"
	"github.com/MenaEnergyVentures/bplus/stm"
)

// Create - the implementation of the Create method. This creates the state entity
func (str SubTypeRegistration) Create(context context.Context, stateEntity stm.StateEntity) (stm.StateEntity, error) {
	//fmt.Printf("State Entity received is %v\n", stateEntity)
	stateEntity.SetState("")
	return str.doProcess(context, stateEntity, "", nil)
}

// Process - process the state entity with an event
func (str SubTypeRegistration) Process(context context.Context, stateEntityId string, event string, param interface{}) (stm.StateEntity, error) {
	//fmt.Printf("Received stateEntityId = %s , event = %s, param = %v\n ", stateEntityId, event, param)

	stateEntity, err := str.StateEntityRepo.Retrieve(stateEntityId)
	if err != nil {
		return nil, e.MakeBplusError(context, e.CannotMakeStateEntity, map[string]interface{}{
			"Error": err.Error()})
	}
	return str.doProcess(context, stateEntity, event, param)
}

func (str SubTypeRegistration) doProcess(context context.Context, stateEntity stm.StateEntity,
	event string, param interface{}) (stm.StateEntity, error) {
	machine, err := str.obtainSTM(context)
	if err != nil {
		return nil, err
	}

	return machine.Process(context, stateEntity, event, param)
}

func (str SubTypeRegistration) obtainSTM(context context.Context) (*stm.Stm, error) {
	return str.OrderSTMChooser(context)
}
