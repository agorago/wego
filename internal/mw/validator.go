package mw

import (
"context"

bplusc "gitlab.intelligentb.com/devops/bplus/context"
fw "gitlab.intelligentb.com/devops/bplus/fw"
e "gitlab.intelligentb.com/devops/bplus/internal/err"
"github.com/go-playground/validator/v10"
)

var VALIDATE = validator.New()

func v10validator(ctx context.Context, chain *fw.MiddlewareChain) context.Context {
	request := bplusc.GetPayload(ctx)

	if errs := validateReq(request); errs != nil {
		return bplusc.SetError(ctx, e.MakeBplusError(ctx, e.ValidationError, map[string]interface{}{
			"Error": encodeV10Errors(errs.(validator.ValidationErrors))}))
	}

	ctx = chain.DoContinue(ctx)
	return ctx
}

func validateReq(req interface{}) error {
	if err := VALIDATE.Struct(req); err != nil {
		return err
	}
	return nil
}

func encodeV10Errors(errs validator.ValidationErrors) []string {
	var errorsSlice []string
	for _, err := range errs {
		errorsSlice = append(errorsSlice, toField(err.Field())+":"+err.Tag())
	}
	return errorsSlice
}

func toField(s string) string {
	field := []byte(s)
	field[0] = field[0] | ('a' - 'A')
	return string(field)
}