package mw

import (
	"bytes"
	bplusc "github.com/MenaEnergyVentures/bplus/context"
	fw "github.com/MenaEnergyVentures/bplus/fw"

	"context"
	"encoding/json"
	"fmt"
	util "github.com/MenaEnergyVentures/bplus/util"
	"io/ioutil"
	"net/http"
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
			return nil, fmt.Errorf("#params passed does not match expected. Actual = %d. Expected = %d",
				len(params), len(od.Params))
		}
	}

	var req *http.Request
	var err error
	var URL = "http://localhost:8080" + od.URL
	// We need to loop thru the params twice. Once to create the request with payload and second time
	// to enhance it with Headers.
	for index, param := range params {
		if od.Params[index].ParamOrigin == fw.PAYLOAD {
			req, err = createRequest(ctx, od.HTTPMethod, URL, param)
		}
		if err != nil {
			return nil, fmt.Errorf("unable to generate HTTP request for param %#v. error is %s", param, err.Error())
		}
	}
	if req == nil {
		req, err = createRequest(ctx, od.HTTPMethod, URL, nil)
		if err != nil {
			return nil, fmt.Errorf("unable to generate HTTP request. error is %s", err.Error())
		}
	}

	for index, param := range params {
		if od.Params[index].ParamOrigin != fw.PAYLOAD {
			enhanceRequest(param, od.Params[index], req)
		}
	}
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return nil, fmt.Errorf("http call failed. err = %s", err.Error())
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("cannot read response body err = %s", err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("status code returned is %d. Error = %s", resp.StatusCode, body)
	}

	var responsePayload = od.OpResponseMaker()
	err = json.Unmarshal(body, &responsePayload)
	if err != nil {
		return nil, fmt.Errorf("Unable to unmarshal response payload.Error = %s", err.Error())
	}
	return responsePayload, nil
}

func createRequest(ctx context.Context, method string, URL string, payload interface{}) (*http.Request, error) {
	var buf *bytes.Buffer
	var err error

	buf, err = constructBytes(payload)
	if err != nil {
		return nil, fmt.Errorf("Cannot construct the message from payload %s", payload)
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
