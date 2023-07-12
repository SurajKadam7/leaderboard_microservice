package redis

import (
	"context"

	redis "github.com/redis/go-redis/v9"
	"github.com/surajkadam/youtube_assignment/cache"
	youtubeerror "github.com/surajkadam/youtube_assignment/errors"
	"github.com/surajkadam/youtube_assignment/model"
)

type Redis struct {
	client *redis.Client
}

func New() cache.Cache {
	return &Redis{}
}

func (r *Redis) Start() (err error) {
	rdb := redis.NewClient(&redis.Options{
		PoolSize: 50,
	})

	_, err = rdb.Ping(context.Background()).Result()

	if err != nil {
		panic("Not able to ping to redis")
	}

	r.client = rdb

	return err
}

func (r *Redis) Viewed(ctx context.Context, key string, video string, incr float64) (res float64, err error) {

	res, err = r.client.ZIncrBy(ctx, key, incr, video).Result()

	if err != nil{
		return res, youtubeerror.ErrNotAbleToIncrement
	}

	return res, err

}

func (r *Redis) TopViewed(ctx context.Context, key string, limit int64) (result []model.Response, err error) {

	res, err := r.client.ZRevRangeByScoreWithScores(ctx, key, &redis.ZRangeBy{
		Min:    "-inf",
		Max:    "+inf",
		Offset: 0,
		Count:  limit,
	}).Result()

	if err != nil {
		return nil, youtubeerror.ErrNotAbleToDisplayTopViewed
	}

	result = []model.Response{}

	for _, val := range res {
		result = append(result, model.Response{
			VideoName: val.Member.(string),
			Viewes:    val.Score,
		})
	}

	return result, nil
}

func (r *Redis) ViewCount(ctx context.Context, key string, video string) (viewes float64, err error) {

	viewes, err = r.client.ZScore(ctx, key, video).Result()

	if err != nil {
		if redis.Nil.Error() == err.Error() {

			return 0, youtubeerror.ErrVideoNotFound
		}
		return 0, err
	}

	return viewes, nil
}
