/**
 * @Time : 2020/6/4 4:29 PM
 * @Author : solacowa@gmail.com
 * @File : service
 * @Software: GoLand
 */

package kitcache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	redisclient "github.com/icowan/redis-client"
)

// get回调方法
type GetCall func(key string) (res interface{}, err error)

// cache 模块
type Service interface {
	// 从缓存拿数据，如果没有，执行call回调拿数据后再存入缓存
	// var m types.Member
	// err = s.cacheSvc.GetCall(ctx, cache.Token.String(token), func(key string) (res interface{}, err error) {
	//	  // 可以是任何逻辑
	//    return s.repository.Member().FindByToken(token)
	// }, time.Hour, &m)
	// fmt.Println(m.Username)
	// key: 存储的key
	// call: 回调方法
	// exp: 存储时间
	// data: 返回的数据
	GetCall(ctx context.Context, key string, call GetCall, exp time.Duration, data interface{}) (err error)

	// 设置缓存
	Set(ctx context.Context, key string, v interface{}, exp time.Duration) (err error)

	// get缓存
	Get(ctx context.Context, key string, data interface{}) (res string, err error)

	// 删除缓存
	Del(ctx context.Context, key string) (err error)
}

type service struct {
	logger    log.Logger
	rds       redisclient.RedisClient
	requestId string
}

func (s *service) Del(ctx context.Context, key string) (err error) {
	logger := log.With(s.logger, s.requestId, ctx.Value(s.requestId))
	if err = s.rds.Del(key); err != nil {
		_ = level.Error(logger).Log("rds", "Del", "err", err.Error())
		return err
	}
	return nil
}

func (s *service) Get(ctx context.Context, key string, data interface{}) (res string, err error) {
	logger := log.With(s.logger, s.requestId, ctx.Value(s.requestId))
	res, err = s.rds.Get(key)
	if err != nil {
		_ = level.Error(logger).Log("rds", "Get", "err", err.Error())
		return
	}
	if data != nil {
		_ = json.Unmarshal([]byte(res), &data)
	}
	return
}

func (s *service) Set(ctx context.Context, key string, v interface{}, exp time.Duration) (err error) {
	logger := log.With(s.logger, s.requestId, ctx.Value(s.requestId))
	err = s.rds.Set(key, v, exp)
	if err != nil {
		_ = level.Error(logger).Log("rds", "Set", "err", err.Error())
	}
	return
}

func (s *service) GetCall(ctx context.Context, key string, call GetCall, exp time.Duration, data interface{}) (err error) {
	logger := log.With(s.logger, s.requestId, ctx.Value(s.requestId))

	resp, err := s.rds.Get(key)
	if err != nil {
		_ = level.Warn(logger).Log("rds", "Get", "key", key, "err", err.Error())
	}
	if resp != "" {
		err = json.Unmarshal([]byte(resp), &data)
		return
	}

	result, err := call(key)
	if err != nil {
		_ = level.Warn(logger).Log("method", "call", "key", key, "err", err.Error())
		return
	}

	b, _ := json.Marshal(result)
	_ = json.Unmarshal(b, &data)

	err = s.rds.Set(key, data, exp)

	return err

}

func New(logger log.Logger, requestId string, redis redisclient.RedisClient) Service {
	return &service{
		logger:    logger,
		rds:       redis,
		requestId: requestId,
	}
}
