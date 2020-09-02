/**
 * @Time : 2020/9/2 9:24 AM
 * @Author : solacowa@gmail.com
 * @File : http
 * @Software: GoLand
 */

package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"net/http"
)

const (
	rateBucketNum = 1000
)

// transport -> logging -> middleware -> endpoint -> service

func MakeHTTPHandler(logger kitlog.Logger, s Service, dmw []endpoint.Middleware, opts []kithttp.ServerOption, requestId, httpPrefix string) http.Handler {
	ems := []endpoint.Middleware{
		// tokenBucketLimitter(rate.NewLimiter(rate.Every(time.Second*1), rateBucketNum)), // 限流 0
	}

	ems = append(ems, dmw...)

	s = NewLoggingServer(logger, s, requestId)

	eps := NewEndpoint(s, map[string][]endpoint.Middleware{
		"Get":     ems,
		"Del":     ems,
		"Set":     ems,
		"GetCall": ems,
	})

	r := mux.NewRouter()

	r.Handle(fmt.Sprintf("%s/get", httpPrefix), kithttp.NewServer(
		eps.GetEndpoint,
		decodeGetRequest,
		encodeJsonResponse,
		opts...,
	)).Methods(http.MethodPost)

	return r
}

func decodeGetRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req getRequest
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	if req.Key == "" {
		return nil, errors.New("Key不能为空")
	}

	return req, nil
}

func encodeJsonResponse(ctx context.Context, w http.ResponseWriter, response interface{}) (err error) {
	return kithttp.EncodeJSONResponse(ctx, w, response)
}
