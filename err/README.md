# Error Handling Framework

WeGO error handling framework has the following features:
1. Support for JSON based error messages.
2. Error handling message code along with a message. This enables consumers to take action depending on the
error that has been emitted.
3. It integrates with WeGO http transport to specify the correct HTTP error codes that need to be emitted.
4. Support for warnings and errors. Sometimes an error is not bad. It can just be a warning and can contain
valid  return values in addition to the warning.
5. Internationalization support - integrates with the i18n module of WeGO.

# Typical usage

The WeGO error handling framework is typically not used directly by the services. Instead, we recommend that 
services have their own error handling wrapper which generates WegoError types from error codes whose 
numeric ranges can be fixed for specific services. For example service 1 might have a range starting from 2000 
and service 2 might have a range starting from 3000. This allows for a clean separation of error code ranges.

Also error messages must be internationalized. To ensure that the error codes are kept consistent we recommend 
defining GO enums to denote the error codes and use i18n bundles to translate the error codes to internationalized
messages.

# Steps to create service specific errors

1. Create a codes.go file in internal/err for the service. (This is automatically generated using Code Gen)
2. Provide error enums. 
3. Use go generate to generate a file that translates the enums to corresponding "strings". Eg: an enum called
FileNotFound must be translated to the string called "FileNotFound". The stringer command with go generate allows
the generation of this translation. This command is automatically inserted into the Makefile by codegen.
4. Write the corresponding detailed error message in the i18n bundle under configs/bundles/en-US folder.


