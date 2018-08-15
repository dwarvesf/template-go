package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/kit/log"

	"<%= domainDir + _.folderName %>/src/endpoints"
	serviceHttp "<%= domainDir + _.folderName %>/src/http"
	"<%= domainDir + _.folderName %>/src/middlewares"
	"<%= domainDir + _.folderName %>/src/postgres"
	"<%= domainDir + _.folderName %>/src/service"
	"<%= domainDir + _.folderName %>/src/service/add"
)

func main() {
	var (
		httpAddr = flag.String("addr", ":3000", "HTTP listen address")
	)
	flag.Parse()

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stdout))
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	// FIXME: replace this with `postgres.New()`
	pgdb, close := postgres.NewFake(os.Getenv("PG_DATASOURCE"))
	defer func() {
		if err := close(); err != nil {
			logger.Log("msg", "failed to close postgres connection", "err", err)
		}
	}()

	var s service.Service
	{
		s = service.Service{
			AddService: middlewares.Compose(
				add.NewPGStore(pgdb),
				add.LoggingMiddleware(logger),
				add.ValidationMiddleware(),
			).(add.Service),
		}
	}

	var h http.Handler
	{
		h = serviceHttp.NewHTTPHandler(
			s,
			endpoints.MakeServerEndpoints(s),
			logger,
			true,
		)
	}

	errs := make(chan error)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		logger.Log("transport", "HTTP", "addr", *httpAddr)
		errs <- http.ListenAndServe(*httpAddr, h)
	}()

	logger.Log("exit", <-errs)
}
