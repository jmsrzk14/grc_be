package main

import (
	"flag"
	"os"

	"grc_be/internal/conf"
	"grc_be/internal/server"
	"grc_be/internal/data"
	"grc_be/internal/biz"
	"grc_be/internal/service"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/env"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// flagconf is the config flag.
var flagconf string

func init() {
	flag.StringVar(&flagconf, "conf", "configs", "config path, eg: -conf config.yaml")
}

func newApp(logger log.Logger, hs *http.Server) *kratos.App {
	return kratos.New(
		kratos.Name("grc_be"),
		kratos.Version("1.0.0"),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(
			hs,
		),
	)
}

func main() {
	flag.Parse()
	logger := log.With(log.NewStdLogger(os.Stdout),
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
	)

	c := config.New(
		config.WithSource(
			file.NewSource(flagconf),
			env.NewSource("GRC_"),
		),
	)
	defer c.Close()

	if err := c.Load(); err != nil {
		panic(err)
	}

	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		panic(err)
	}

	// Railway dynamic port override
	if port := os.Getenv("PORT"); port != "" {
		bc.Server.Http.Addr = ":" + port
	}

	// Manual Dependency Injection (tanpa Wire agar cepat)
	d, cleanup, err := data.NewData(bc.Data, logger)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	tenantRepo := data.NewTenantRepo(d, logger)
	propertyRepo := data.NewPropertyRepo(d, logger)
	// tpRepo := data.NewTenantPropertyRepo(d, logger) // Optional for now
	regRepo := data.NewRegulationRepo(d, logger)
	itemRepo := data.NewRegulationItemRepo(d, logger)
	mappingRepo := data.NewRegulationPropertyMappingRepo(d, logger)
	sessionRepo := data.NewAssessmentSessionRepo(d, logger)
	resultRepo := data.NewAssessmentResultRepo(d, logger)
	raRepo := data.NewRegulationAssessmentRepo(d, logger)

	tenantUC := biz.NewTenantUseCase(tenantRepo, logger)
	propertyUC := biz.NewPropertyUseCase(propertyRepo, logger)
	regUC := biz.NewRegulationUseCase(regRepo, itemRepo, mappingRepo, logger)
	assUC := biz.NewAssessmentUseCase(sessionRepo, resultRepo, raRepo, itemRepo, logger)

	tenantSvc := service.NewTenantService(tenantUC, logger)
	propSvc := service.NewPropertyService(propertyUC, logger)
	regSvc := service.NewRegulationService(regUC, logger)
	assSvc := service.NewAssessmentService(assUC, logger)

	httpSrv := server.NewHTTPServer(bc.Server, tenantSvc, propSvc, regSvc, assSvc, logger)

	app := newApp(logger, httpSrv)

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		panic(err)
	}
}
