/**
 * @Time : 2020/9/2 5:43 PM
 * @Author : solacowa@gmail.com
 * @File : grpc_test
 * @Software: GoLand
 */

package cache

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/log/level"
	"github.com/go-kit/kit/transport"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"github.com/google/uuid"
	"github.com/icowan/kit-cache/pb"
	"google.golang.org/grpc"
	googlegrpc "google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/metadata"
	"net"
	"testing"
	"time"
)

func TestMakeGRPCHandler(t *testing.T) {
	grpcOpts := []kitgrpc.ServerOption{
		kitgrpc.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		kitgrpc.ServerBefore(func(ctx context.Context, mds metadata.MD) context.Context {
			ctx = context.WithValue(ctx, "request-id", uuid.New().String())
			return ctx
		}),
	}

	grpcListener, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		_ = logger.Log("transport", "gRPC", "during", "Listen", "err", err)
		panic(err)
	}
	baseServer := googlegrpc.NewServer()
	svc := NewSvc()

	pb.RegisterCacheServer(baseServer, MakeGRPCHandler(logger, svc, nil, grpcOpts, "request-id"))
	_ = level.Info(logger).Log("grpc", "server", "start", "success")
	err = baseServer.Serve(grpcListener)
	panic(err)
}

func TestGrpcServer_Set(t *testing.T) {
	go func() {
		TestMakeGRPCHandler(t)
	}()
	time.Sleep(time.Second * 3)

	ctx, cel := context.WithTimeout(context.Background(), time.Second*30)
	defer cel()

	conn, err := grpc.DialContext(ctx, "localhost:50051",
		grpc.WithInsecure(), grpc.WithBlock(),
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, roundrobin.Name)),
		grpc.WithBlock(), grpc.WithBackoffMaxDelay(time.Second),
	)
	if err != nil {
		_ = level.Error(logger).Log("grpc", "connect", "err", err)
		panic(err)
	}

	defer func() {
		_ = conn.Close()
	}()

	cacheSvc := pb.NewCacheClient(conn)

	res, err := cacheSvc.Set(context.Background(), &pb.SetRequest{
		Key: "g-hello",
		Val: "298fupaoisdjf;ajsdf",
		Exp: int64(time.Second * 50),
	})

	if err != nil {
		t.Error(err)
	}

	t.Log(res)
	t.Log("set", "success")

	var data string

	resp, err := cacheSvc.Get(context.Background(), &pb.GetRequest{
		Key:  "g-hello",
		Data: data,
	})

	if err != nil {
		t.Error(err)
	}

	t.Log(res)
	t.Log("get", "success", "resp", resp.GetData())
	t.Log("get", "success", "data", data)
}
