package main

import (
	"context"
	"fmt"
	"hexagonal-architexture-utils/config"
	"hexagonal-architexture-utils/internal/adapters/db"
	"hexagonal-architexture-utils/internal/adapters/http"
	"hexagonal-architexture-utils/internal/adapters/http/metrics"
	"hexagonal-architexture-utils/internal/application/core/api"
	domain "hexagonal-architexture-utils/internal/domains"
	"hexagonal-architexture-utils/internal/pkg/logging"
	"hexagonal-architexture-utils/internal/pkg/otel"
	"os"
	"os/signal"
	"syscall"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
)

var version string = "development"

func main() {
	c, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	ctx := context.WithValue(c, domain.XTRACEID, uuid.NewString())
	defer cancel()
	defer logging.Global.Sync()
	defer logging.Global.Info(ctx, "api stopped successfully")

	metrics.InitMetrics()

	err := config.Setup(ctx)
	if err != nil {
		logging.Global.Error(ctx, fmt.Sprintf("failed to setup [%v]", err))
		return
	}
	config.BuildVersion = version

	traceProvider, err := otel.InitProvider(ctx, config.Configuration.ServiceName(), config.Configuration.OtelURL())
	if err != nil {
		logging.Global.Error(ctx, fmt.Sprintf("error initializing telemetry provider [%v]", err))
	}
	defer traceProvider.Shutdown(ctx)

	dbAdapter, err := db.NewAdapter(ctx, &config.Configuration)
	if err != nil {
		logging.Global.Error(ctx, fmt.Sprintf("failed to connect to postgres [%v]", err))
		return
	}
	defer dbAdapter.Close()

	application := api.NewApplication(dbAdapter)

	httpAdapter := http.NewAdapter(application, config.Configuration.HostAddress(), 80)

	g, gCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		httpAdapter.Start(gCtx, nil)
		return nil
	})

	g.Go(func() error {
		<-gCtx.Done()
		httpAdapter.Stop(ctx, config.Configuration.ShutDownTimeout())
		return nil
	})

	if err := g.Wait(); err != nil {
		fmt.Printf("exit reason: %s \n", err)
	}
}
