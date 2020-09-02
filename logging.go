/**
 * @Time : 2020/9/2 9:24 AM
 * @Author : solacowa@gmail.com
 * @File : logging
 * @Software: GoLand
 */

package cache

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

type loggingServer struct {
	logger    log.Logger
	requestId string
	Service
}

func NewLoggingServer(logger log.Logger, s Service, requestId string) Service {
	return &loggingServer{
		logger:    level.Info(logger),
		Service:   s,
		requestId: requestId,
	}
}

func (s *loggingServer) Get(ctx context.Context, key string, data interface{}) (res string, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.requestId, ctx.Value(s.requestId),
			"method", "Get",
			"key", key,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.Get(ctx, key, data)
}