package fw

// The WEGO command catalog - a catalog of all commands in the system
// Every command,strategy, service, singleton (or whatever else you want to call it) is registered as
// a command in this catalog. It does not do anything more than contain all the commands
// with a name. The consumers will use the name to accomplish all wiring
// Wiring helpers can exist independent of the catalog
type CommandCatalog interface{
	Command(name string)interface{}
	RegisterCommand(name string, command interface{})
}

var commandCatalog CommandCatalog
func init(){
	commandCatalog =  &commandCatalogImpl{
		make(map[string]interface{}),
	}
}

//MakeCommandCatalog - Returns an example of a Command Catalog
func MakeCommandCatalog()CommandCatalog{
	return commandCatalog
}

type commandCatalogImpl struct{
	commandsMap map[string]interface{}
}

func (cc *commandCatalogImpl)Command(name string)interface{}{
	return cc.commandsMap[name]
}

func (cc *commandCatalogImpl)RegisterCommand(name string,command interface{}){
	cc.commandsMap[name] = command
}