package context

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/magiconair/properties/assert"
)

func TestIt(t *testing.T) {
	ctx := context.Background()
	ctx = Add(ctx, "xxx", "yyy")

	val := Value(ctx, "xxx")
	assert.Equal(t, val, "yyy", "value of xxx does not match yyy")
}

func TestEiphenated(t *testing.T) {
	ctx := context.Background()
	ctx = Add(ctx, "Accept-Language", "yyy")

	val := Value(ctx, "Accept-Language")
	assert.Equal(t, val, "yyy", "value of Accept-Language does not match yyy")
}

func TestHttpRequestEnhance(t *testing.T) {
	r := httptest.NewRequest("POST", "localhost:5000/path/to/url?abc=123", nil)
	r.Header =
		http.Header{
			"x": []string{"y"},
		}
	ctx := context.Background()
	ctx = Enhance(ctx, r)
	val := Value(ctx, "x")
	assert.Equal(t, val, "y", "value of x does not match y")
	val = Value(ctx, "abc")
	assert.Equal(t, val, "123", "value of abc does not match 123")
}

func TestHttpRequestEnhanceMux(t *testing.T) {
	r := httptest.NewRequest("POST", "localhost:5000/path/to/url/123/345", nil)
	r.Header =
		http.Header{
			"x": []string{"y"},
		}
	r = mux.SetURLVars(r, map[string]string{"abc": "123"})
	ctx := context.Background()
	ctx = Enhance(ctx, r)
	val := Value(ctx, "x")
	t.Logf("Value of x is %s\n", val)
	val = Value(ctx, "abc")
	t.Logf("Value of abc is %s\n", val)

	t.Logf("All keys = %v\n", GetAllKeys(ctx))
}
