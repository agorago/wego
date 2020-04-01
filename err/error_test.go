package err_test

import (
	"context"
	"github.com/magiconair/properties/assert"
	"gitlab.intelligentb.com/devops/bplus/err"
	"net/http"
	"testing"
)

func TestMakeErrWithHTTPCode(t *testing.T) {
	e := err.MakeErrWithHTTPCode(context.TODO(),err.Error,123,"Some_Message",http.StatusBadRequest,
		nil)
	assert.Equal(t,123,e.ErrorCode)
	assert.Equal(t,"Some_Message",e.ErrorMessage)
	assert.Equal(t,400,e.HTTPErrorCode)
	assert.Equal(t,err.Error,e.LogLevel)
}
