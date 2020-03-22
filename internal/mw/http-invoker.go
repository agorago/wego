package mw

import (
	"bytes"
	"fmt"

	bplusc "gitlab.intelligentb.com/devops/bplus/context"
	fw "gitlab.intelligentb.com/devops/bplus/fw"
	e "gitlab.intelligentb.com/devops/bplus/internal/err"

	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	util "gitlab.intelligentb.com/devops/bplus/util"
)

// HTTPInvoker - invokes the service via HTTP. The last middleware in the proxy
// middleware chain
func HTTPInvoker(ctx context.Context, chain *fw.MiddlewareChain) context.Context {
	od := fw.GetProxyOperationDescriptor(ctx)
	params := bplusc.GetProxyRequestParams(ctx)
	resp, err := httpInvoker(ctx, od, params)

	ctx = bplusc.SetProxyResponsePayload(ctx, resp)
	ctx = bplusc.SetProxyResponseError(ctx, err)
	return ctx
}

func httpInvoker(ctx context.Context, od fw.OperationDescriptor, params []interface{}) (interface{}, error) {
	if len(params) != len(od.Params) {
		if od.Params[0].ParamOrigin == fw.CONTEXT {
			params = append([]interface{}{ctx}, params...)
		} else {
			return nil, e.MakeBplusError(ctx, e.ParamsNotExpected, map[string]interface{}{
				"Actual": len(params), "Expected": len(od.Params)})
		}
	}

	var req *http.Request
	var err error
	var URL = "http://localhost:5000/" + od.Service.Name + od.URL
	// We need to loop thru the params twice. Once to create the request with payload and second time
	// to enhance it with Headers.
	for index, param := range params {
		if od.Params[index].ParamOrigin == fw.PAYLOAD {
			req, err = createRequest(ctx, od.HTTPMethod, URL, param)
		}
		if err != nil {
			return nil, e.MakeBplusError(ctx, e.CannotGenerateHTTPRequest, map[string]interface{}{
				"Param": param, "Error": err.Error()})
		}
	}
	if req == nil {
		req, err = createRequest(ctx, od.HTTPMethod, URL, nil)
		if err != nil {
			return nil, e.MakeBplusError(ctx, e.CannotGenerateHTTPRequest1, map[string]interface{}{
				"Error": err.Error()})
		}
	}

	for index, param := range params {
		if od.Params[index].ParamOrigin != fw.PAYLOAD {
			enhanceRequest(param, od.Params[index], req)
		}
	}
	bplusc.CopyHeadersToHTTPRequest(ctx, req)
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return nil, e.MakeBplusError(ctx, e.HTTPCallFailed, map[string]interface{}{
			"Error": err.Error()})
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, e.MakeBplusError(ctx, e.CannotReadResponseBody, map[string]interface{}{
			"Error": err.Error()})
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		fmt.Printf("The status code is %d.Body is %s\n", resp.StatusCode, body)
		return nil, e.MakeBplusErrorWithErrorCode(ctx, resp.StatusCode, e.Non200StatusCodeReturned, map[string]interface{}{
			"StatusCode": resp.StatusCode, "Body": fmt.Sprintf("%s", body),
		})
	}

	if od.OpResponseMaker != nil {
		var responsePayload, _ = od.OpResponseMaker(ctx)
		err = json.Unmarshal(body, &responsePayload)
		if err != nil {
			return nil, e.MakeBplusError(ctx, e.ResponseUnmarshalException, map[string]interface{}{
				"Error": err.Error()})
		}
		return responsePayload, nil
	}
	return nil, nil
}

func createRequest(ctx context.Context, method string, URL string, payload interface{}) (*http.Request, error) {
	var buf *bytes.Buffer
	var err error

	buf, err = constructBytes(payload)
	if err != nil {
		return nil, e.MakeBplusError(ctx, e.CannotGenerateHTTPRequest, map[string]interface{}{
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

func enhanceRequest(param interface{}, pd fw.ParamDescriptor, req *http.Request) {
	// fmt.Printf("Setting header for %s to %#v\n", pd.Name, param)
	req.Header.Set(pd.Name, util.ConvertToString(param, pd.ParamKind))
}
