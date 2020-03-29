# Swagger documentation

BPlus supports a service driven spec i.e. it is possible to generate an Open API spec from the service.

# Enable Swagger documentation

When you register a service ensure that the description fields in both the ServiceDescriptor and 
OperationDescriptor are properly populated. 
The following fields are of interest:
* ServiceDescriptor.Description
* OperationDescriptor.Description
* OperationDescriptor.RequestDescription
* OperationDescriptor.ResponseDescription

# Can we do this automatically when we generate the service using Code Gen?
Yes. It is possible. Just make sure that your interface file is properly commented. *code-gen* uses the 
following comments to populate the fields mentioned above:

* Docs and comments for the interface will be copied to ServiceDescriptor.Description
* Docs and comments for the Operations within the interface are used to populate OperationDescriptor.Description
* Docs and comments for the Request Payload struct to populate OperationDescriptor.RequestDescription
* Docs and comments for the Request Payload struct to populate OperationDescriptor.ResponseDescription

# What if my service is already generated?
Look at register.go and enhance the fields mentioned above to whatever you want.

# How to generate swagger
The following steps need to be followed:
1. Make sure that you generate swagger-gen executable for your service. Run **Make swagger-gen-build**
2. Run **make swagger-docs** to use the swagger-gen to generate the swagger docs in *internal/docs*
3. swagger-gen is not perfect yet. Go to internal/docs and make sure that there are no syntax errors.
4. Run **make swagger-build** to generate the swagger.yaml file in the root folder of the project.
5. Run **make swagger-serve** to actually serve your swagger documentation.

