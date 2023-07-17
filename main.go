package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/log"
	"github.com/oklog/oklog/pkg/group"
	redis "github.com/redis/go-redis/v9"
	leaderendpoint "github.com/surajkadam/youtube_assignment/leaderboard/endpoint"
	leaderservice "github.com/surajkadam/youtube_assignment/leaderboard/service"
	leadertransport "github.com/surajkadam/youtube_assignment/leaderboard/transport"
	rediscache "github.com/surajkadam/youtube_assignment/repository/redis"
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

func main() {
	data, err := os.ReadFile("config.json")
	if err != nil {
		panic("not able load the configurations")
	}
	c := config{}
	json.Unmarshal(data, &c)

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

	// using Server struct so that I can handle the shutdown grasefully
	l, err := net.Listen("tcp", c.Port)
	if err != nil {
		panic(err)
	}

	h := http.Server{
		Handler: mux,
	}

	var g group.Group
	{
		g.Add(
			func() error {
				logger.Log("info", "startig the server", "port :", c.Port)

				return h.Serve(l)
			},
			func(error) {
				logger.Log("info", "server shutdown initialize")

				h.Shutdown(context.Background())

				logger.Log("info", "server shutdown completed...")
			},
		)
	}

	cancelInterrupt := make(chan struct{})

	g.Add(func() error {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		select {
		case sig := <-c:
			return fmt.Errorf("received signal %s", sig)
		case <-cancelInterrupt:
			return nil
		}
	}, func(error) {
		close(cancelInterrupt)
	})
	g.Run()
}
