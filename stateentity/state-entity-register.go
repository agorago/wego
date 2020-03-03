package stateentity

import (
	"context"

	"gitlab.intelligentb.com/devops/bplus/stm"
)

// STMChooser - this is a strategy that allows the sub type to choose an STM
// given the context
type STMChooser func(context context.Context) (*stm.Stm, error)

// SubTypeMaker - a function that generates the correct subtype
type SubTypeMaker func(context context.Context) (stm.StateEntity, error)

// SubTypeRegistration - a way to register state entity sub types to this
// so that their activity gets exposed in the STM
type SubTypeRegistration struct {
	Name                    string
	StateEntitySubTypeMaker SubTypeMaker
	ActionCatalog           map[string]interface{}
	URLPrefix               string // can be the same as the Name
	StateEntityRepo         Repo
	OrderSTMChooser         STMChooser
}

var allSubTypes = make(map[string]SubTypeRegistration)

// Service - the interface that is going to be implemented generically for all state entities
type Service interface {
	Create(context.Context, stm.StateEntity) (stm.StateEntity, error)
	Process(context.Context, string, string, interface{}) (stm.StateEntity, error)
}

// Repo - the state entity repository
type Repo interface {
	Create(stm.StateEntity) (stm.StateEntity, error)
	Retrieve(ID string) (stm.StateEntity, error)
}

// RegisterSubType - all state entity sub types register here in their init()
func RegisterSubType(str SubTypeRegistration) {
	allSubTypes[str.Name] = str
	str.setupStateEntityService()
}
