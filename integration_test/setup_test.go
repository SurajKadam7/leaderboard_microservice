package integration_test

import (
	"context"
	"encoding/json"
	"fmt"

	log2 "log"
	"net/http"
	"os"
	"time"

	"github.com/go-kit/log"
	redis "github.com/redis/go-redis/v9"
	leaderendpoint "github.com/surajkadam/youtube_assignment/leaderboard/endpoint"
	leaderservice "github.com/surajkadam/youtube_assignment/leaderboard/service"
	leadertransport "github.com/surajkadam/youtube_assignment/leaderboard/transport"
	rediscache "github.com/surajkadam/youtube_assignment/repo/redis"
	viewendpoint "github.com/surajkadam/youtube_assignment/view/endpoint"
	viewservice "github.com/surajkadam/youtube_assignment/view/service"
	viewtransport "github.com/surajkadam/youtube_assignment/view/transport"
)

type config struct {
	Key      string `json:"key"`
	Port     string `json:"port"`
	Address  string `json:"address"`
	PoolSize int    `json:"poolSize"`
	UserName string `json:"username"`
	Password string `json:"password"`
}

func getHandler() (http.Handler, func()) {
	data, err := os.ReadFile("config.json")
	if err != nil {
		panic("not able load the configurations")
	}
	c := config{}
	json.Unmarshal(data, &c)

	c.Key = "youtube_assignment_test"
	c.Port = "8080"
	c.Address = "localhost:6379"
	c.PoolSize = 50

	log2.Printf("\nconfigurations : %+v\n\n", c)

	var client *redis.Client
	{
		client = redis.NewClient(&redis.Options{
			Addr:     c.Address,
			PoolSize: c.PoolSize,
			Username: c.UserName,
			Password: c.Password,
		})

		_, err := client.Ping(context.Background()).Result()

		if err != nil {
			panic("Not able to ping to redis")
		}
	}

	rds := rediscache.New(client, c.Key)

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	var viewService viewservice.Service
	{
		viewService = viewservice.New(rds)
		viewService = viewservice.LoggingMiddleware(logger)(viewService)
	}

	viewEndpoint := viewendpoint.New(viewService)
	viewHttpServer := viewtransport.NewHTTPHandler(viewEndpoint, logger)

	var leaderService leaderservice.Service
	{
		leaderService = leaderservice.New(rds)
		leaderService = leaderservice.LoggingMiddleware(logger)(leaderService)
	}

	leaderEndpoint := leaderendpoint.New(leaderService)
	leaderHttpServer := leadertransport.NewHTTPHandler(leaderEndpoint, logger)

	mux := http.NewServeMux()

	mux.Handle("/video/", viewHttpServer)
	mux.Handle("/top/", leaderHttpServer)

	clearRedis := func() {
		client.Del(context.Background(), c.Key, getKey(c.Key))
	}

	return mux, clearRedis
}

func getKey(key string) string {
	y, m, d := time.Now().Date()
	return fmt.Sprintf("%s:%d-%d-%d", key, y, m, d)
}