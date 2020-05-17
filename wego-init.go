package wego

import (
	_ "github.com/agorago/wego/config" // init config functions
	"github.com/agorago/wego/fw"
	wegohttp "github.com/agorago/wego/http"
	_ "github.com/agorago/wego/i18n" // init i18n functions
	"github.com/agorago/wego/internal/mw"
	"github.com/agorago/wego/stateentity"
	"net/http"
)

const (
	WegoHTTPHandler                = "WegoHTTPHandler"
	WEGO                           = "WEGO"
	ProxyService                   = "ProxyService"
	StateEntityProxy               = "StateEntityProxy"
	StateEntityRegistrationService = "StateEntityRegistrationService"
)
func MakeWegoInitializer()fw.Initializer{
	return wegoinit{}
}

type wegoinit struct{}
func (wegoinit)ModuleName()string{
	return WEGO
}
// The WEGO initializer
func (wegoinit)Initialize(commandCatalog fw.CommandCatalog)(fw.CommandCatalog,error){
	wego := fw.MakeRegistrationService()
	ep := mw.MakeEntrypoint()
	httphandler := wegohttp.InitializeHTTP(wego,ep)
	commandCatalog.RegisterCommand(WegoHTTPHandler,httphandler)
	commandCatalog.RegisterCommand(WEGO,wego)

	// create proxy - basic proxy and state entity proxy
	pep := mw.MakeProxyEntrypoint()
	proxyService := wegohttp.MakeProxyService(wego,pep)
	stateEntityProxy := stateentity.MakeStateEntityProxy(proxyService)
	commandCatalog.RegisterCommand(StateEntityProxy,stateEntityProxy)
	commandCatalog.RegisterCommand(ProxyService,proxyService)

	// create subentityregistrationservice
	sers := stateentity.MakeStateEntityRegistrationService(wego)
	commandCatalog.RegisterCommand(StateEntityRegistrationService,sers)
	return commandCatalog,nil
}

// Have Getters for every one of the commands that have been initialized above
func GetWego(commandCatalog fw.CommandCatalog)(fw.RegistrationService,error){
	return commandCatalog.Command(WEGO).(fw.RegistrationService),nil
}

func GetHTTPHandler(commandCatalog fw.CommandCatalog)(http.Handler,error){
	return commandCatalog.Command(WegoHTTPHandler).(http.Handler),nil
}

func GetProxyService(commandCatalog fw.CommandCatalog)(wegohttp.ProxyService,error){
	return commandCatalog.Command(ProxyService).(wegohttp.ProxyService),nil
}

func GetStateEntityProxy(commandCatalog fw.CommandCatalog)(stateentity.StateEntityProxy,error){
	return commandCatalog.Command(StateEntityProxy).(stateentity.StateEntityProxy),nil
}

func GetStateEntityRegistrationService(commandCatalog fw.CommandCatalog)(stateentity.StateEntityRegistrationService,error){
	return commandCatalog.Command(WegoHTTPHandler).(stateentity.StateEntityRegistrationService),nil
}



