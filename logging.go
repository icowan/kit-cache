/**
 * @Time : 2020/9/2 9:24 AM
 * @Author : solacowa@gmail.com
 * @File : logging
 * @Software: GoLand
 */

package kitcache

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

type logging struct {
	logger  log.Logger
	next    Service
	traceId string
}

func (s *logging) GetCall(ctx context.Context, key string, call GetCall, exp time.Duration, data interface{}) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "GetCall", "key", key, "call", call, "exp", exp, "data", data,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.GetCall(ctx, key, call, exp, data)
}

func (s *logging) Set(ctx context.Context, key string, v interface{}, exp time.Duration) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "Set", "key", key, "v", v, "exp", exp,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.Set(ctx, key, v, exp)
}

func (s *logging) Get(ctx context.Context, key string, data interface{}) (res string, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "Get", "key", key, "data", data,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.Get(ctx, key, data)
}

func (s *logging) Del(ctx context.Context, key string) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "Del", "key", key,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.Del(ctx, key)
}

func NewLogging(logger log.Logger, traceId string) Middleware {
	logger = log.With(logger, "kitcache", "logging")
	return func(next Service) Service {
		return &logging{
			logger:  level.Info(logger),
			next:    next,
			traceId: traceId,
		}
	}
}
