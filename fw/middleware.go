package fw

import "context"

// Sets up the middleware that would be used throughout WeGO

// Middleware - can be used as the Middleware for any WeGO pipeline
type Middleware interface{
	Intercept(context.Context, *MiddlewareChain) context.Context
}

// MiddlewareChain - a middleware chain that
type MiddlewareChain struct {
	middlewares []Middleware
	index       int
}

// MakeChain - makes a new middleware chain with nothing in it
// this can be added to by calling the Add method below
func MakeChain() MiddlewareChain {
	return MiddlewareChain{
		index: 0,
	}
}

// DoContinue - call this function to continue the middleware execution.
// important: make sure that the call looks like this
// ctx = chain.DoContinue()
// return ctx
func (chain *MiddlewareChain) DoContinue(ctx context.Context) context.Context {
	if len(chain.middlewares) <= chain.index {
		return ctx
	}
	next := chain.middlewares[chain.index]
	chain.index++
	return next.Intercept(ctx, chain)
}

// Add - add a new middleware into the chain
func (chain *MiddlewareChain) Add(m ...Middleware) {
	chain.middlewares = append(chain.middlewares, m...)
}
