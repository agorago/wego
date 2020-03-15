# BPlus
The framework for Intelligent B

# Introduction

BPlus is a framework that is suitable for rapid development of micro services. It has been tuned towards the 
requirements of Intelligent B but has nothing specific hard coded inside it. It can easily be used as a framework for creating
developing micro services anywhere.

The chief benefit of BPlus is that it comes integrated with different frameworks that have been standardized across
Intelligent B such as 
- gorilla MUX (for HTTP)
- Uber zap (for logging)
- New Relic (for monitoring)
- GoDog (for BDD)
- Viper (for configuration management)
- Nick Snyder's go-i18n 

Besides integration with frameworks, BPlus accomplishes the following:
* Provide a standardized folder structure for development.
* Allows developers to focus on business logic without being bothered about how services are exposed via HTTP etc.
* Facilitates interface service separation
* Integrates with Make for builds. See [approach to builds](build.md)
* Integrates with Docker for deployment. See [docker integration](docker.md)
* Exposes HTTP for the services. For details see the [http module](http/README.md)
* Enforces a common set of middlewares. For details see the [discussion on middleware](internal/mw/README.md)
* Facilitates rapid development using a code generator. For details see the 
[discussion on code-gen](https://gitlab.intelligentb.com/devops/code-gen/README.md)
* Implements a state machine. See [the BPlus State Machine](stm/README.md)
* Provides a consistent way of handling Error Handling. See [Error Handling](err/README.md)
* Provides the following functionalities as horizontal services. (a growing set)
    # Logging. See [BPlus Logging](log/README.md)
    # Auditing (not implemented yet)
    # Monitoring. See [BPlus integration with New Relic](new-relic.md)
    # Configuration Management. See [BPlus approach to Configuration Management](config.md)

# Modularized code base

BPlus is a modularized code base with a base framework and several set of plug-ins. All the plug-ins
are shipped together in one GO Module. However, it is easy to dissociate them if the need so arises.

# Folder Structure
BPlus is best implemented using a folder structure as follows:
$ mkdir src
$ cd src
$ git clone intelligentb.gitlab.com/bplus (get the actual git repo url)
$ git clone intelligentb.gitlab.com/code-gen (get the actual git repo url)
$ mkdir configs

## Project Types
We would have three sets of projects that need to be developed. All these project templates are generated 
automatically using code-gen. These include the following:
1. A service project. This is the project that a developer uses to develop a micro service.
2. A deploy project. This gives a devops person the flexibility to package multiple micro services into a 
single deployable.
3. A state service project. This is a project that uses a state machine for developing services. The 
chief difference between this and the service project is that this project is driven by a state machine.
The code structuring is around the state machine. 




