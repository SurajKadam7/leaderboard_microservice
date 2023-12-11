package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"

	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	kitLog "github.com/go-kit/log"
	"github.com/oklog/oklog/pkg/group"
	redis "github.com/redis/go-redis/v9"
	conf "github.com/SurajKadam7/leaderboard_microservice/config"
	leaderendpoint "github.com/SurajKadam7/leaderboard_microservice/leaderboard/endpoint"
	leaderservice "github.com/SurajKadam7/leaderboard_microservice/leaderboard/service"
	leadertransport "github.com/SurajKadam7/leaderboard_microservice/leaderboard/transport"
	rediscache "github.com/SurajKadam7/leaderboard_microservice/repo/redis"
	viewendpoint "github.com/SurajKadam7/leaderboard_microservice/view/endpoint"
	viewservice "github.com/SurajKadam7/leaderboard_microservice/view/service"
	viewtransport "github.com/SurajKadam7/leaderboard_microservice/view/transport"

	capi "github.com/hashicorp/consul/api"
)

func main() {
	// loding the configurations
	consulUrl, err := loadConfig()
	if err != nil {
		log.Fatal(err)
	}

	// using consul for getting the configurations
	var consulClient *capi.Client
	{
		defaultConfig := capi.DefaultConfig()
		defaultConfig.Address = consulUrl
		client, err := capi.NewClient(defaultConfig)

		if err != nil {
			panic(err)
		}

		consulClient = client
	}

	kv, _, err := consulClient.KV().Get("youtube_assignment", nil)

	if err != nil {
		panic(fmt.Errorf("error while featching kv from consul err : %w", err))
	}

	c, err := conf.New(bytes.NewBuffer(kv.Value))

	if err != nil {
		panic(err)
	}

	var client *redis.Client
	{
		client = redis.NewClient(&redis.Options{
			Addr:     c.RedisAddress,
			PoolSize: c.PoolSize,
		})

		_, err := client.Ping(context.Background()).Result()

		if err != nil {
			panic("Not able to ping to redis")
		}
	}

	rds := rediscache.New(client, c.Key)

	var logger kitLog.Logger
	{
		logger = kitLog.NewLogfmtLogger(os.Stderr)
		logger = kitLog.With(logger, "ts", kitLog.DefaultTimestampUTC)
		logger = kitLog.With(logger, "caller", kitLog.DefaultCaller)
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

func loadConfig() (consulUrl string, err error) {
	var config struct {
		Consul string `json:"consul"`
	}
	file, err := os.Open("config.json")
	if err != nil {
		return "", fmt.Errorf("error while opening config : %w", err)
	}
	json.NewDecoder(file).Decode(&config)
	return config.Consul, nil
}
