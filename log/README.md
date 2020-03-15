# BPlus Logging

The BPlus logger uses zap to perform logging at the following levels:
* Debug - Mostly used in development mode.
* Info - Can be useful to trace the flow of logic. 
* Warn - For Warning messages. Useful to find problems within code paths
* Error - For more severe errors.

All messages above Info level will be written to New Relic as well. So please use Info level with caution to 
avoid clutter.

## Fields

All logging is performed in JSON. The logging level, message and time stamp are available using the following
JSON fields:
* level - in capital letters for each level such as DEBUG, ERROR etc.
* message
* time in ISO 8601 format

Additionally, the log system also creates additional fields from the context. Hence context is mandatory.
The list of fields created using context include the following:
* Trace ID - this keeps a trace of the HTTP request. Trace IDs are passed from one service to the other 
in an execution chain. They can be very helpful in tracing a request. If no trace ID is available (such as
during startup) NO_TRACE_ID is printed.

## Sinks

Currently, logging is done to stdout and all errors are printed to stderr. These logs need to be collected 
and sent to one consolidated log. This must be achieved during deployment. In development mode logs get 
printed to the console.

# Additional Log fields

The following fields are also enabled or disabled as per configuration:
* caller - this gives additional information about the caller of the log. (goland code source, line # etc.). 
This is optional and can be configured using the configuration value "bplus.disable_caller". 
* stacktrace - this gives additional information about the entire call stack trace to the logger. This is
configured using the configuration value "bplus.disable_stacktrace'

# Signatures

To optimize log performance and to enhance usability, the following signatures are supported for each logging
level. This is shown for Info level below but is applicable to other levels as well.

```golang
    func Infof(ctx context.Context, message string, args ...interface{})
    func Info(ctx context.Context, message string)
    func InfoWithFieldsf(ctx context.Context, fields map[string]string,message string, args ...interface{})
    func InfoWithFields(ctx context.Context,fields map[string]string,  message string)
``` 

A few observations about the signatures above:
1. Context is mandatory for logging. Context must be passed around. If context is not available (such as during)
startup, please use context.Background()
2. Message is the logging message 
3. args are printf style args supported in the signatures that end with f.
4. The WithFields signature allows custom fields to be logged using a string map. This can be useful in machine
parsing of logs and is strongly encouraged.

# Configurations

The configuration system allows the logging to be configured. Some of them are mentioned above (see the section
on additional fields above). Here is a full list of applicable logger configurations:

There are some more configurations such as the ones below:
* bplus.development_mode - true/false puts zap in development mode
* bplus.log_level - info/debug/error/warn the level used for logging. 
* bplus.disable_stacktrace - true/false
* bplus.disable_caller - true/false

