package stateentity

import (
	"context"
	wegohttp "github.com/agorago/wego/http"
	"github.com/agorago/wego/stm"
)

type StateEntityProxy struct{
	proxy wegohttp.ProxyService
}

func MakeStateEntityProxy(service wegohttp.ProxyService) StateEntityProxy {
	return StateEntityProxy{
		proxy: service,
	}
}
func (proxy StateEntityProxy)doProcess(ctx context.Context, entityName string, op string, args ...interface{}) (stm.StateEntity, error) {
	ent, err1 := proxy.proxy.ProxyRequest(ctx, entityName, op, args...)
	if err1 != nil {
		return nil, err1
	}
	en := ent.(stm.StateEntity)
	return en, nil
}

func (proxy StateEntityProxy) Create(ctx context.Context, entityName string, entity stm.StateEntity) (stm.StateEntity, error) {
	return proxy.doProcess(ctx, entityName, "Create", entity)
}

func (proxy StateEntityProxy) Process(ctx context.Context, entityName, id, eventID string, param interface{}) (stm.StateEntity, error) {
	return proxy.doProcess(ctx, entityName, "Process", id, eventID, param)
}
