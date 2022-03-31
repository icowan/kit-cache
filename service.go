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
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"strings"
	"time"
)

// GetCall get回调方法
type GetCall func(key string) (res interface{}, err error)

type Middleware func(Service) Service

// Service cache 模块
type Service interface {
	// GetCall 从缓存拿数据，如果没有，执行call回调拿数据后再存入缓存
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
	// Set 设置缓存
	Set(ctx context.Context, key string, v interface{}, exp time.Duration) (err error)
	// Get get缓存
	Get(ctx context.Context, key string, data interface{}) (res string, err error)
	// Del 删除缓存
	Del(ctx context.Context, key string) (err error)
}

type service struct {
	rds redis.Cmdable
}

func (s *service) Del(ctx context.Context, key string) (err error) {
	if err = s.rds.Del(ctx, key).Err(); err != nil {
		return err
	}
	return nil
}

func (s *service) Get(ctx context.Context, key string, data interface{}) (res string, err error) {
	res = s.rds.Get(ctx, key).Val()
	if !strings.EqualFold(res, "") && data != nil {
		_ = json.Unmarshal([]byte(res), &data)
	}
	return
}

func (s *service) Set(ctx context.Context, key string, v interface{}, exp time.Duration) (err error) {
	err = s.rds.Set(ctx, key, v, exp).Err()
	return
}

func (s *service) GetCall(ctx context.Context, key string, call GetCall, exp time.Duration, data interface{}) (err error) {
	resp := s.rds.Get(ctx, key).Val()
	if resp != "" {
		err = json.Unmarshal([]byte(resp), &data)
		return
	}

	result, err := call(key)
	if err != nil {
		err = errors.Wrap(err, "call")
		return
	}

	b, _ := json.Marshal(result)
	_ = json.Unmarshal(b, &data)

	err = s.rds.Set(ctx, key, data, exp).Err()

	return err

}

func New(redis redis.Cmdable) Service {
	return &service{
		rds: redis,
	}
}
