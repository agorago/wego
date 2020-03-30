package mw

import (
	"context"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	bplusc "gitlab.intelligentb.com/devops/bplus/context"
	"gitlab.intelligentb.com/devops/bplus/fw"
	e "gitlab.intelligentb.com/devops/bplus/internal/err"
	"gitlab.intelligentb.com/devops/bplus/log"
)

var VALIDATE = validator.New()

func init() {
	// RegisterTagNameFunc registers a function to get alternate
	// names for StructFields.
	VALIDATE.RegisterTagNameFunc(
		func(fld reflect.StructField) string {
			// Use the names which have been specified for JSON representations of structs,
			// rather than normal Go field names
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}

			return name
		})
}

func v10validator(ctx context.Context, chain *fw.MiddlewareChain) context.Context {
	request := bplusc.GetPayload(ctx)
	if request == nil {
		return chain.DoContinue(ctx)
	}

	if errs := validateReq(request); errs != nil {
		er, ok := errs.(validator.ValidationErrors)
		if !ok {
			// if it is not validator errors then just log them and move on.
			log.Infof(ctx,
				"Error in validation - which are not of types validator.ValidationErrors. err = %#v",
				errs)
			ctx = chain.DoContinue(ctx)
			return ctx
		}
		return bplusc.SetError(ctx, e.MakeBplusErrorWithErrorCode(ctx, http.StatusBadRequest, e.ValidationError,
			map[string]interface{}{
				"Error": encodeV10Errors(er)}))
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
		errorsSlice = append(errorsSlice, err.Field()+":"+err.Tag())
	}
	return errorsSlice
}
