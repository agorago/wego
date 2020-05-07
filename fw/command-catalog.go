package fw

// The WEGO command catalog - a catalog of all commands in the system
// Every command,strategy, service, singleton (or whatever else you want to call it) is registered as
// a command in this catalog. It does not do anything more than contain all the commands
// with a name. The consumers will use the name to accomplish all wiring
// Wiring helpers can exist independent of the catalog
type CommandCatalog interface{
	Command(name string)interface{}
	RegisterCommand(name string, command interface{})
	IsModuleRegistered(module string) bool
	SetModuleRegistered(module string)
}

var commandCatalog CommandCatalog
func init(){
	commandCatalog =  &commandCatalogImpl{
		commandsMap: make(map[string]interface{}),
		initializedMap: make(map[string]bool),
	}
}

//MakeCommandCatalog - Returns an example of a Command Catalog
func MakeCommandCatalog()CommandCatalog{
	return commandCatalog
}

type commandCatalogImpl struct{
	commandsMap map[string]interface{}
	initializedMap map[string]bool
}

func (cc *commandCatalogImpl)Command(name string)interface{}{
	return cc.commandsMap[name]
}

func (cc *commandCatalogImpl)RegisterCommand(name string,command interface{}){
	cc.commandsMap[name] = command
}

func (cc *commandCatalogImpl)IsModuleRegistered(module string) bool{
	return cc.initializedMap[module]
}

func (cc *commandCatalogImpl)SetModuleRegistered(module string) {
	cc.initializedMap[module] = true
}

// Initializers  take the burden out of initializing a bunch of commands in the correct order
// by calling their "Make" functions. They also handle injection. Initializers can in theory be
// augmented with dependency injection frameworks. (such as wire)
// However, initializers facilitate modularized dependency injection by linking commands across
// modules using a command catalog.
// Initialize functions might rely on other commands already being present in the commandCatalog.
// This provides a light weight version of DI without tying it to "magic" DI containers that contain
// too much hidden logic.
// 1. Prerequisites - this contains all the pre-requisites for initialization. This might include
// commands in the Command Catalog that are expected to be populated, some configs that need
// to be specified in the config and some environment variables that need to be set.
// 2. Results of the initialize()- These are the commands that the initialize() creates.
// These commands are populated in the commandCatalog which is returned
type Initializer interface{
	ModuleName() string
	Initialize(CommandCatalog)(CommandCatalog,error)
}

func MakeInitializedCommandCatalog(initializers ...Initializer)(CommandCatalog,error){
	commandCatalog := MakeCommandCatalog()
	var err error
	for _,initializer := range initializers{
		commandCatalog,err = initializer.Initialize(commandCatalog)
		if err != nil{
			return commandCatalog,err
		}
		commandCatalog.SetModuleRegistered(initializer.ModuleName())
	}
	return commandCatalog,nil
}
