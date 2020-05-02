package stateentity

import (
	"context"
	bplusHTTP "gitlab.intelligentb.com/devops/bplus/http"
	"gitlab.intelligentb.com/devops/bplus/stm"
)

type StateEntityProxy struct{}

func MakeStateEntityProxy() StateEntityProxy {
	return StateEntityProxy{}
}
func doProcess(ctx context.Context, entityName string, op string, args ...interface{}) (stm.StateEntity, error) {
	ent, err1 := bplusHTTP.ProxyRequest(ctx, entityName, op, args...)
	if err1 != nil {
		return nil, err1
	}
	en := ent.(stm.StateEntity)
	return en, nil
}

func (StateEntityProxy) Create(ctx context.Context, entityName string, entity stm.StateEntity) (stm.StateEntity, error) {
	return doProcess(ctx, entityName, "Create", entity)
}

func (StateEntityProxy) Process(ctx context.Context, entityName, id, eventID string, param interface{}) (stm.StateEntity, error) {
	return doProcess(ctx, entityName, "Process", id, eventID, param)
}
