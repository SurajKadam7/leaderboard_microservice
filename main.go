package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/log"
	"github.com/oklog/oklog/pkg/group"
	"github.com/surajkadam/youtube_assignment/cache/redis"
	"github.com/surajkadam/youtube_assignment/config"
	leaderendpoint "github.com/surajkadam/youtube_assignment/leaderboard/endpoint"
	leaderservice "github.com/surajkadam/youtube_assignment/leaderboard/service"
	leadertransport "github.com/surajkadam/youtube_assignment/leaderboard/transport"
	viewendpoint "github.com/surajkadam/youtube_assignment/view/endpoint"
	viewservice "github.com/surajkadam/youtube_assignment/view/service"
	viewtransport "github.com/surajkadam/youtube_assignment/view/transport"
)

const (
	key  = "youtube_assignment"
	port = ":8080"
)

func main() {

	cache := redis.New()
	// started redis server
	cache.Start()

	c := config.New(key, cache)

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	var (
		viewService    = viewservice.New(c, logger)
		viewEndpoint   = viewendpoint.New(viewService)
		viewHttpServer = viewtransport.NewHTTPHandler(viewEndpoint, logger)
	)

	var (
		leaderService    = leaderservice.New(c, logger)
		leaderEndpoint   = leaderendpoint.New(leaderService)
		leaderHttpServer = leadertransport.NewHTTPHandler(leaderEndpoint, logger)
	)

	l, err := net.Listen("tcp", port)

	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()

	mux.Handle("/youtube/video/", viewHttpServer)
	mux.Handle("/youtube/top/", leaderHttpServer)

	// using Server struct so that I can handle the shutdown grasefully
	h := http.Server{
		Handler: mux,
	}

	var g group.Group
	{
		g.Add(
			func() error {
				logger.Log("info", "startig the server", "port :", port)

				return h.Serve(l)
			},
			func(error) {
				logger.Log("info", "server shutdown initialize")

				h.Shutdown(context.Background())
				l.Close()

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
