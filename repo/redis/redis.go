package rediscache

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	redis "github.com/redis/go-redis/v9"
	youtubeerror "github.com/SurajKadam7/leaderboard_microservice/errors"
	"github.com/SurajKadam7/leaderboard_microservice/model"
	cache "github.com/SurajKadam7/leaderboard_microservice/repo"
)

type Redis struct {
	key    string
	client *redis.Client
}

func New(c *redis.Client, key string) cache.Repository {
	return &Redis{
		client: c,
		key:    key,
	}
}

func (r *Redis) Viewed(ctx context.Context, video string, incr int64) (res float64, err error) {

	// returning result of the lifetime viewes only
	res, err = r.client.ZIncrBy(ctx, r.key, float64(incr), video).Result()

	if err != nil {
		return res, youtubeerror.ErrNotAbleToIncrement
	}

	key := getKey(r.key)
	_, err = r.client.ZIncrBy(ctx, key, float64(incr), video).Result()

	if err != nil {
		err = checkError(err)
		return
	}

	return res, err

}

func (r *Redis) DayTopViewed(ctx context.Context, limit int64) (result []model.ViedeoDetails, err error) {
	key := getKey(r.key)
	result, err = topViewed(ctx, r.client, key, limit)
	return
}

func (r *Redis) LifetimeTopViewed(ctx context.Context, limit int64) (result []model.ViedeoDetails, err error) {

	result, err = topViewed(ctx, r.client, r.key, limit)
	return
}

func (r *Redis) DayViewCount(ctx context.Context, video string) (viewes float64, err error) {
	key := getKey(r.key)
	viewes, err = viewCount(ctx, r.client, video, key)
	return
}

func (r *Redis) LifetimeViewCount(ctx context.Context, video string) (viewes float64, err error) {

	viewes, err = viewCount(ctx, r.client, video, r.key)
	return
}

func (r *Redis) AddVideos(ctx context.Context, videos []model.Video) (err error) {
	key1 := r.key
	key2 := getKey(r.key)

	for _, video := range videos {
		member := redis.Z{
			Score:  0,
			Member: video.Name,
		}
		// adding videos redis ..
		{
			_, err = r.client.ZAdd(ctx, key1, member).Result()
			if err != nil {
				return err
			}
		}

		{
			_, err = r.client.ZAdd(ctx, key2, member).Result()
			if err != nil {
				return err
			}
		}
	}

	return
}

func getKey(key string) string {
	y, m, d := time.Now().Date()
	return fmt.Sprintf("%s:%d-%d-%d", key, y, m, d)
}

func topViewed(ctx context.Context, r *redis.Client, key string, limit int64) (result []model.ViedeoDetails, err error) {
	res, err := r.ZRevRangeByScoreWithScores(ctx, key, &redis.ZRangeBy{
		Min:    "-inf",
		Max:    "+inf",
		Offset: 0,
		Count:  limit,
	}).Result()

	if err != nil {
		err = checkError(err)
		return
	}

	result = []model.ViedeoDetails{}

	for _, val := range res {
		nameStr, ok := val.Member.(string)
		if !ok {
			return result, nil
		}
		result = append(result, model.ViedeoDetails{
			VideoName: nameStr,
			Viewes:    val.Score,
		})
	}

	return result, nil
}

func viewCount(ctx context.Context, r *redis.Client, video string, key string) (viewes float64, err error) {

	viewes, err = r.ZScore(ctx, key, video).Result()

	if err != nil {
		err = checkError(err)
	}

	return viewes, err

}

func checkError(err error) error {
	if err == redis.Nil {
		return youtubeerror.ErrVideoNotFound
	}
	if strings.Contains(errors.Unwrap(err).Error(), "connection refused") {
		return youtubeerror.ErrDBDown
	}
	return err
}
