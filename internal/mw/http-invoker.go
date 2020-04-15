package mw

import (
	"bytes"
	"fmt"
	"gitlab.intelligentb.com/devops/bplus/config"
	"gitlab.intelligentb.com/devops/bplus/log"

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
	od := fw.GetOperationDescriptor(ctx)
	resp, err := httpInvoker(ctx, od)

	ctx = bplusc.SetResponsePayload(ctx, resp)
	ctx = bplusc.SetError(ctx, err)
	return ctx
}

func httpInvoker(ctx context.Context, od fw.OperationDescriptor) (interface{}, error) {
	var req *http.Request
	var err error
	var URL = "http://localhost:" + config.Value("bplus.port")
	if s := config.Value(od.Service.Name + "_hostname_port");  s != "" {
		URL = "http://" + s
	}
	URL +=  "/" + od.Service.Name + od.URL
	payload := bplusc.GetPayload(ctx)
	req, err = createRequest(ctx, od.HTTPMethod, URL,payload)
	if err != nil {
		return nil, e.MakeBplusError(ctx, e.CannotGenerateHTTPRequest1, map[string]interface{}{
			"Error": err.Error()})
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
		log.Warnf(ctx,"The status code is %d.Body is %s\n", resp.StatusCode, body)
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
