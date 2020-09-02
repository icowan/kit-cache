/**
 * @Time : 2020/9/2 5:43 PM
 * @Author : solacowa@gmail.com
 * @File : service_test
 * @Software: GoLand
 */

package cache

import (
	"context"
	"github.com/go-kit/kit/log"
	redisclient "github.com/icowan/redis-client"
	"testing"
	"time"
)

var logger log.Logger

func NewSvc() Service {
	logger = log.NewLogfmtLogger(log.StdlibWriter{})
	rds, err := redisclient.NewRedisClient("127.0.0.1:6379", "", "test:", 0)
	if err != nil {
		panic(err)
	}
	return New(logger, "tarce-id", rds)
}

func TestService_Get(t *testing.T) {
	svc := NewSvc()

	var data string

	res, err := svc.Get(context.Background(), "hello", &data)
	if err != nil {
		t.Error(err)
	}

	println(res)
	t.Log("success", res)
}

func TestService_Set(t *testing.T) {
	svc := NewSvc()
	err := svc.Set(context.Background(), "hello", "world", time.Second*10)
	if err != nil {
		t.Error(err)
	}
	t.Log("success")
}
