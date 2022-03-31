/**
 * @Time : 2020/9/2 5:43 PM
 * @Author : solacowa@gmail.com
 * @File : service_test
 * @Software: GoLand
 */

package kitcache

import (
	"context"
	"github.com/go-kit/kit/log"
	"testing"
	"time"
)

var logger log.Logger

func NewSvc() Service {
	logger = log.NewLogfmtLogger(log.StdlibWriter{})
	//rds, err := redisclient.NewRedisClient("127.0.0.1:32679", "123456", "test:", 0, nil)
	//if err != nil {
	//	panic(err)
	//}
	return New(nil)
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
