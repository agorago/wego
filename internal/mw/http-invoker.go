package mw

import (
	"bytes"
	"fmt"
	"github.com/agorago/wego/config"
	wegoe "github.com/agorago/wego/err"
	"github.com/agorago/wego/log"

	wegocontext "github.com/agorago/wego/context"
	fw "github.com/agorago/wego/fw"
	e "github.com/agorago/wego/internal/err"

	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type HTTPInvoker struct{}
// HTTPInvoker - invokes the service via HTTP. The last middleware in the proxy
// middleware chain
func (HTTPInvoker)Intercept(ctx context.Context, _ *fw.MiddlewareChain) context.Context {
	od := fw.GetOperationDescriptor(ctx)
	resp, err := httpInvoker(ctx, od)

	ctx = wegocontext.SetResponsePayload(ctx, resp)
	ctx = wegocontext.SetError(ctx, err)
	return ctx
}

func httpInvoker(ctx context.Context, od fw.OperationDescriptor) (interface{}, error) {
	var req *http.Request
	var err error
	var URL = "http://localhost:" + config.Value("wego.port")
	if s := config.Value(od.Service.Name + "_hostname_port"); s != "" {
		URL = "http://" + s
	}
	URL += "/" + od.Service.Name + od.URL
	payload := wegocontext.GetPayload(ctx)
	req, err = createRequest(ctx, od.HTTPMethod, URL, payload)
	if err != nil {
		return nil, err
	}
	wegocontext.CopyHeadersToHTTPRequest(ctx, req)
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return nil, e.Error(ctx, e.HTTPCallFailed, map[string]interface{}{
			"Error": err.Error()})
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, e.Error(ctx, e.CannotReadResponseBody, map[string]interface{}{
			"Error": err.Error()})
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		log.Warnf(ctx, "The status code is %d.Body is %s\n", resp.StatusCode, body)
		return nil, extractErrorResponse(ctx, resp, body)
	}

	if od.OpResponseMaker != nil {
		var responsePayload, _ = od.OpResponseMaker(ctx)
		err = json.Unmarshal(body, &responsePayload)
		if err != nil {
			return nil, e.Error(ctx, e.ResponseUnmarshalException, map[string]interface{}{
				"Error": err.Error()})
		}
		return responsePayload, nil
	}
	return nil, nil
}

func extractErrorResponse(ctx context.Context, resp *http.Response, body []byte) wegoe.WeGOError {
	er := wegoe.WeGOError{}
	uerr := json.Unmarshal(body, &er)
	if uerr != nil {
		// cannot unmarshal the body into a WeGO err. So create an error.
		return e.HTTPError(ctx, resp.StatusCode, e.Non200StatusCodeReturned, map[string]interface{}{
			"StatusCode": resp.StatusCode, "Body": fmt.Sprintf("%s", body),
		})
	}
	return er
}
func createRequest(ctx context.Context, method string, URL string, payload interface{}) (*http.Request, error) {
	var buf *bytes.Buffer
	var err error

	buf, err = constructBytes(payload)
	if err != nil {
		return nil, e.Error(ctx, e.CannotGenerateHTTPRequest, map[string]interface{}{
			"Error": payload})
	}

	return http.NewRequestWithContext(ctx, method, URL, buf)
}

func constructBytes(payload interface{}) (*bytes.Buffer, error) {
	if payload == nil {
		payload = make(map[string]string)
		//return &bytes.Buffer{}, nil
	}

	buf, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	return bytes.NewBuffer(buf), nil
}
