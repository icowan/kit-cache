/**
 * @Time : 2020/9/2 9:24 AM
 * @Author : solacowa@gmail.com
 * @File : grpc
 * @Software: GoLand
 */

package kitcache

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport/grpc"
	"github.com/icowan/kit-cache/pb"
	"time"
)

type grpcServer struct {
	get     grpc.Handler
	set     grpc.Handler
	del     grpc.Handler
	getCall grpc.Handler
}

func (g *grpcServer) Get(ctx context.Context, req *pb.GetRequest) (*pb.Response, error) {
	_, rep, err := g.get.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.Response), nil
}

func (g *grpcServer) Set(ctx context.Context, req *pb.SetRequest) (*pb.Response, error) {
	_, rep, err := g.set.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.Response), nil
}

func MakeGRPCHandler(logger log.Logger, s Service, dmw []endpoint.Middleware, opts []grpc.ServerOption, requestId string) pb.CacheServer {
	var ems []endpoint.Middleware

	ems = append(ems, dmw...)

	s = NewLoggingServer(logger, s, requestId)
	eps := NewEndpoint(s, map[string][]endpoint.Middleware{
		"Get":     ems,
		"Set":     ems,
		"Del":     ems,
		"GetCall": ems,
	})

	return &grpcServer{
		get: grpc.NewServer(
			eps.GetEndpoint,
			decodeGRPCGetRequest,
			encodeResponse,
			opts...,
		),
		set: grpc.NewServer(
			eps.SetEndpoint,
			decodeGRPCSetRequest,
			encodeResponse,
			opts...,
		),
		del: grpc.NewServer(
			eps.DelEndpoint,
			decodeGRPCGetRequest,
			encodeResponse,
			opts...,
		),
		getCall: nil,
	}
}

func decodeGRPCGetRequest(_ context.Context, r interface{}) (interface{}, error) {
	return getRequest{
		Key: r.(*pb.GetRequest).Key,
	}, nil
}

func decodeGRPCSetRequest(_ context.Context, r interface{}) (interface{}, error) {
	return setRequest{
		Key: r.(*pb.SetRequest).Key,
		Val: r.(*pb.SetRequest).Val,
		Exp: time.Duration(r.(*pb.SetRequest).Exp),
	}, nil
}

func encodeResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp, ok := r.(string)
	if !ok {
		resp = ""
	}

	var err error
	return &pb.Response{
		Data: resp,
	}, err
}
