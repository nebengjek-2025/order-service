package helpers

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	error "order-service/src/pkg/http-error"
	"order-service/src/pkg/utils"

	"go.elastic.co/apm/module/apmhttp"
	"golang.org/x/net/context/ctxhttp"
)

type HttpPostFormRequestPayload struct {
	Url      string
	FormData url.Values
	Result   interface{}
}

func HttpPostFormRequest(payload HttpPostFormRequestPayload, ctx context.Context) utils.Result {
	var result utils.Result
	req, err := http.NewRequest("POST", payload.Url, strings.NewReader(payload.FormData.Encode()))

	if err != nil {
		errObj := error.NewInternalServerError()
		errObj.Message = fmt.Sprintf("\"%s\"", utils.ConvertString(err.Error()))
		result.Error = errObj
		return result
	}

	req.Close = true
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// --- do request

	newClient := http.Client{
		Timeout: 15 * time.Second,
	}

	var wrapClient = apmhttp.WrapClient(&newClient)

	resp, err := ctxhttp.Do(ctx, wrapClient, req)

	if err, ok := err.(net.Error); ok && err.Timeout() {
		errObj := error.NewInternalServerError()
		errObj.Message = "request timeout 10s."
		result.Error = errObj
		return result
	}

	if err != nil {
		errObj := error.NewInternalServerError()
		errObj.Message = fmt.Sprintf("\"%s\"", utils.ConvertString(err.Error()))
		result.Error = errObj
		return result
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errObj := error.NewInternalServerError()
		errObj.Message = "request error."
		result.Error = errObj
		return result
	}

	readResp, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		errObj := error.NewInternalServerError()
		errObj.Message = fmt.Sprintf("\"%s\"", utils.ConvertString(err.Error()))
		result.Error = errObj
		return result
	}

	if err := json.Unmarshal(readResp, &payload.Result); err != nil {
		errObj := error.NewInternalServerError()
		errObj.Message = "cannot marshal response payload"
		result.Error = errObj
		return result
	}
	result.Data = payload.Result
	return result
}
