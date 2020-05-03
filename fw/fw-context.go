package fw

// An extension to the context file to take care of operation descriptors which are
// framework specific. These functions should not be in the context.go file within the
// context package since that introduces a dependency on the fw package from there
import (
	"context"
	wegoc "github.com/agorago/wego/context"
)

// GetOperationDescriptor - gets the operation descriptor from the context
func GetOperationDescriptor(ctx context.Context) OperationDescriptor {
	return wegoc.Value(ctx, "OD").(OperationDescriptor)
}

// SetOperationDescriptor - sets the operation descriptor and returns the enhance context
func SetOperationDescriptor(ctx context.Context, od OperationDescriptor) context.Context {
	return wegoc.Add(ctx, "OD", od)
}

// GetProxyOperationDescriptor - gets the operation descriptor from the context
func GetProxyOperationDescriptor(ctx context.Context) OperationDescriptor {
	return wegoc.Value(ctx, "PROXY_OD").(OperationDescriptor)
}

// SetProxyOperationDescriptor - sets the operation descriptor and returns the enhance context
func SetProxyOperationDescriptor(ctx context.Context, od OperationDescriptor) context.Context {
	return wegoc.Add(ctx, "PROXY_OD", od)
}
