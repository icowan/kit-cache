/**
 * @Time : 2020/9/2 9:24 AM
 * @Author : solacowa@gmail.com
 * @File : endpoint
 * @Software: GoLand
 */

package cache

import (
	"context"
	"time"

	"github.com/go-kit/kit/endpoint"
)

type (
	getCallRequest struct {
		Key string `json:"key"`
	}

	getRequest struct {
		Key string `json:"key"`
	}

	setRequest struct {
		Key string        `json:"key"`
		Val interface{}   `json:"val"`
		Exp time.Duration `json:"exp"`
	}
)

type Endpoints struct {
	GetCallEndpoint endpoint.Endpoint
	GetEndpoint     endpoint.Endpoint
	SetEndpoint     endpoint.Endpoint
	DelEndpoint     endpoint.Endpoint
}

func NewEndpoint(s Service, mdw map[string][]endpoint.Middleware) Endpoints {
	eps := Endpoints{
		GetCallEndpoint: makeGetCallEndpoint(s),
		GetEndpoint:     makeGetEndpoint(s),
		SetEndpoint:     makeSetEndpoint(s),
		DelEndpoint:     makeDelEndpoint(s),
	}

	for _, m := range mdw["GetCall"] {
		eps.GetCallEndpoint = m(eps.GetCallEndpoint)
	}
	for _, m := range mdw["Get"] {
		eps.GetEndpoint = m(eps.GetEndpoint)
	}
	for _, m := range mdw["Set"] {
		eps.SetEndpoint = m(eps.SetEndpoint)
	}
	for _, m := range mdw["Del"] {
		eps.DelEndpoint = m(eps.DelEndpoint)
	}

	return eps
}

func makeGetCallEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return nil, nil
	}
}

func makeGetEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getRequest)
		res, err := s.Get(ctx, req.Key, nil)
		return res, err
	}
}

func makeSetEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(setRequest)
		err = s.Set(ctx, req.Key, req.Val, req.Exp)
		return nil, err
	}
}

func makeDelEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getRequest)
		err = s.Del(ctx, req.Key)
		return nil, err
	}
}
