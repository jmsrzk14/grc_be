// @title GRC API
// @version 1.0
// @description API documentation for Governance, Risk, and Compliance system.

// @tag.name TenantsService

// @tag.name PropertiesService

// @tag.name RegulationsService

// @tag.name AssessmentsService

// @tag.name RiskService

package main

import (
	"context"
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
		bc.Server.HTTP.Addr = ":" + port
	}

	// Database Source override (Railway fallback)
	if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
		log.Context(context.Background()).Infof("Database URL detected from environment: %s", dbURL)
		if bc.Data != nil && bc.Data.Database != nil {
			bc.Data.Database.Source = dbURL
		}
	} else {
		log.Context(context.Background()).Info("No DATABASE_URL environment variable found, using config file.")
	}

	log.Context(context.Background()).Infof("Attempting to connect to database: %s", bc.Data.Database.Source)


	// Manual Dependency Injection (tanpa Wire agar cepat)
	d, cleanup, err := data.NewData(bc.Data, logger)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	tenantRepo := data.NewTenantRepo(d, logger)
	propertyRepo := data.NewPropertyRepo(d, logger)
	tpRepo := data.NewTenantPropertyRepo(d, logger)
	regRepo := data.NewRegulationRepo(d, logger)
	itemRepo := data.NewRegulationItemRepo(d, logger)
	mappingRepo := data.NewRegulationPropertyMappingRepo(d, logger)
	sessionRepo := data.NewAssessmentSessionRepo(d, logger)
	resultRepo := data.NewAssessmentResultRepo(d, logger)
	raRepo := data.NewRegulationAssessmentRepo(d, logger)
	authRepo := data.NewAuthRepo(d, logger)
	riskRepo := data.NewRiskRepo(d, logger)
	riskCatRepo := data.NewRiskCategoryRepo(d, logger)
	riskCatTenantRepo := data.NewRiskCategoryTenantRepo(d, logger)

	tenantUC := biz.NewTenantUseCase(tenantRepo, logger)
	propertyUC := biz.NewPropertyUseCase(propertyRepo, tpRepo, logger)
	regUC := biz.NewRegulationUseCase(regRepo, itemRepo, mappingRepo, logger)
	assUC := biz.NewAssessmentUseCase(sessionRepo, resultRepo, raRepo, itemRepo, logger)
	authUC := biz.NewAuthUseCase(authRepo, logger)
	riskUC := biz.NewRiskUseCase(riskRepo, riskCatRepo, riskCatTenantRepo, logger)

	tenantSvc := service.NewTenantService(tenantUC, propertyUC, logger)
	propSvc := service.NewPropertyService(propertyUC, logger)
	regSvc := service.NewRegulationService(regUC, logger)
	assSvc := service.NewAssessmentService(assUC, logger)
	authSvc := service.NewAuthService(authUC, logger)
	riskSvc := service.NewRiskService(riskUC, logger)

	httpSrv := server.NewHTTPServer(bc.Server, tenantSvc, propSvc, regSvc, assSvc, authSvc, riskSvc, logger)

	app := newApp(logger, httpSrv)

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		panic(err)
	}
}
