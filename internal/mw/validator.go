package mw

import (
	"context"
	"net/http"

	"github.com/go-playground/validator/v10"
	bplusc "gitlab.intelligentb.com/devops/bplus/context"
	fw "gitlab.intelligentb.com/devops/bplus/fw"
	e "gitlab.intelligentb.com/devops/bplus/internal/err"
)

var VALIDATE = validator.New()

func v10validator(ctx context.Context, chain *fw.MiddlewareChain) context.Context {
	request := bplusc.GetPayload(ctx)
	if request == nil {
		return chain.DoContinue(ctx)
	}

	if errs := validateReq(request); errs != nil {
		return bplusc.SetError(ctx, e.MakeBplusErrorWithErrorCode(ctx, http.StatusBadRequest, e.ValidationError, map[string]interface{}{
			"Error": encodeV10Errors(errs.(validator.ValidationErrors))}))
	}

	ctx = chain.DoContinue(ctx)
	return ctx
}

// validate the request with the validation struct defined
func validateReq(req interface{}) error {
	if err := VALIDATE.Struct(req); err != nil {
		return err
	}
	return nil
}

// encodes the validation error to string
func encodeV10Errors(errs validator.ValidationErrors) []string {
	var errorsSlice []string
	for _, err := range errs {
		errorsSlice = append(errorsSlice, toField(err.Field())+":"+err.Tag())
	}
	return errorsSlice
}

// converts the err.Field() to Field Name (trims the quotes)
func toField(s string) string {
	field := []byte(s)
	field[0] = field[0] | ('a' - 'A')
	return string(field)
}
