# The HTTP Module for BPlus

The BPlus http module. This module is an "extension" for the base framework of BPlus. It exposes the HTTP
transport. The module has methods to both expose a service via HTTP and to consume a HTTP service by providing
a HTTP proxy.

# Exposing a service via HTTP

The http transport registers itelf under BPlus [framework](../fw/README.md). It gets called when every 
Service registers its operations under BPlus. The URL configured for the operation is exposed to the 
outside world via HTTP. 

The HTTP module currently supports JSON as the default encoding mechanism. 