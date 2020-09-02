/**
 * @Time : 2020/9/2 9:24 AM
 * @Author : solacowa@gmail.com
 * @File : middleware
 * @Software: GoLand
 */

package cache

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/pkg/errors"
	"golang.org/x/time/rate"
)

var ErrLimitExceed = errors.New("Rate limit exceed!")

func tokenBucketLimitter(bkt *rate.Limiter) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			if !bkt.Allow() {
				return nil, ErrLimitExceed
			}
			return next(ctx, request)
		}
	}
}
