# Context 

This package is idiomatically used and referred to as bplusc. In code, it is imported as: 
```
    import bplusc "gitlab.intelligentb.com/devops/bplus/context"
```
The bplusc package has the following uses:

1. It provides an abstraction to use the context.Context object provided by Go Lang. The methods in bplusc
provide a key that can be used to get and set the values in context. The context.Context documentation 
recommends the use of a custom key to store and retrieve attributes in context. bplusc provides such a key.

2. This package provides convenience methods that allows the getting and setting of specific keys and values. Eg:
__SetPayload()__ allows the payload to be set in the context. Get methods obviate the use of casting.

3. Every attribute added to the context is remembered and can be retrieved using: 
```
    bplusc.GetAllKeys(ctx)
```
4. bplusc provides convenience methods to copy headers from and to http request to the context object.

5. It allows the generation of a trace ID for every request. 
