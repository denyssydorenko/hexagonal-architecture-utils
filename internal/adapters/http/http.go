package http

import (
	"context"
	"encoding/json"
	"fmt"
	"hexagonal-architexture-utils/config"
	"hexagonal-architexture-utils/docs"
	"hexagonal-architexture-utils/internal/adapters/http/metrics"
	domain "hexagonal-architexture-utils/internal/domains"
	"hexagonal-architexture-utils/internal/domains/db"
	"hexagonal-architexture-utils/internal/domains/health"
	"hexagonal-architexture-utils/internal/pkg/logging"
	"hexagonal-architexture-utils/internal/ports"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoswagger "github.com/swaggo/echo-swagger"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

const (
	URL    = "url"
	METHOD = "method"
)

type Adapter struct {
	api    ports.ApiPort
	host   string
	port   int
	server *echo.Echo
}

func NewAdapter(api ports.ApiPort, host string, port int) *Adapter {
	adapter := Adapter{
		api:  api,
		host: host,
		port: port,
	}

	// server
	e := echo.New()

	// middleware
	e.Use(middleware.Recover())
	// log request and response
	e.Use(middleware.BodyDumpWithConfig(middleware.BodyDumpConfig{
		Skipper: InfoEndpointsSkipper,
		Handler: func(c echo.Context, reqBody, resBody []byte) {
			ctx := c.Request().Context()
			url := c.Request().URL.String()
			method := c.Request().Method

			reqFields := []zap.Field{zap.String(METHOD, method), zap.String(URL, url)}
			if len(reqBody) > 0 {
				var jsonMapReq interface{}
				json.Unmarshal(reqBody, &jsonMapReq)
				reqFields = append(reqFields, zap.Any("body:", jsonMapReq))
			}
			logging.Global.Info(ctx, "request", reqFields...)

			respFields := []zap.Field{zap.String(METHOD, method), zap.String(URL, url)}
			if len(resBody) > 0 {
				var jsonMapResp interface{}
				json.Unmarshal(resBody, &jsonMapResp)
				respFields = append(respFields, zap.Any("body:", jsonMapResp))
			}
			logging.Global.Info(ctx, "response", respFields...)
		},
	}))
	// add timeout
	e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Timeout:      config.Configuration.Timeout(),
		ErrorMessage: "hexagonal-architecture-utils timeout",
	}))
	// add x-trace-id to response
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			res := c.Response()
			id := req.Header.Get(domain.XTRACEID)
			if id == "" {
				id = uuid.New().String()
			}
			res.Header().Set(domain.XTRACEID, id)

			ctx := context.WithValue(req.Context(), domain.XTRACEID, id)
			request := c.Request().WithContext(ctx)
			c.SetRequest(request)

			return next(c)
		}
	})

	// enable traces
	e.Use(otelecho.Middleware(config.Configuration.ServiceName()))

	// assign server to adapter
	adapter.server = e

	// init swagger
	adapter.initSwagerInfo()

	// add handlers
	adapter.addHandlers()

	return &adapter
}

func InfoEndpointsSkipper(c echo.Context) bool {
	if strings.HasPrefix(c.Path(), config.Configuration.BasePath()+"/metrics") ||
		strings.HasPrefix(c.Path(), config.Configuration.BasePath()+"/swagger") ||
		strings.HasPrefix(c.Path(), config.Configuration.BasePath()+"/health") {
		return true
	}
	return false
}

func (ad *Adapter) initSwagerInfo() {
	docs.SwaggerInfo.Title = "Hexagonal Architecture Utils"
	docs.SwaggerInfo.Description = "Hexagonal Architecture Utils example project"
	docs.SwaggerInfo.Version = "2.0"
	docs.SwaggerInfo.BasePath = config.Configuration.BasePath()
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
}

func (ad *Adapter) addHandlers() {
	ad.server.GET(config.Configuration.BasePath()+"/metrics", metrics.Metrics())
	ad.server.GET(config.Configuration.BasePath()+"/health", ad.Health)
	ad.server.GET(config.Configuration.BasePath()+"/swagger/*", echoswagger.WrapHandler)

	// DB
	ad.server.PUT(config.Configuration.BasePath()+"/v1/db/create", ad.CreateDB)
	ad.server.GET(config.Configuration.BasePath()+"/v1/db/get", ad.GetDB)
	ad.server.GET(config.Configuration.BasePath()+"/v1/db/get-all", ad.GetAllDB)
	ad.server.POST(config.Configuration.BasePath()+"/v1/db/update", ad.UpdateDB)
	ad.server.DELETE(config.Configuration.BasePath()+"/v1/db/delete", ad.DeleteDB)
}

func (ad *Adapter) setTracingXTraceId(ctx context.Context) {
	span := trace.SpanFromContext(ctx)
	xtraceidString := fmt.Sprintf("%s", ctx.Value(domain.XTRACEID))
	span.SetAttributes(attribute.String(domain.XTRACEID, xtraceidString))
}

func (ad *Adapter) updateMetrics(endpoint string, err error, start time.Time) {
	metrics.APIDuration.WithLabelValues(endpoint, strconv.FormatBool(err == nil)).Observe(time.Since(start).Seconds())
	metrics.APIRequests.WithLabelValues(endpoint, strconv.FormatBool(err == nil)).Inc()
}

func (ad *Adapter) Start(ctx context.Context, shutdownCallback func()) {
	err := ad.server.Start(ad.host)
	if err != nil && err != http.ErrServerClosed {
		logging.Global.Error(ctx, fmt.Sprintf("http server encountered an error, initiating shutdown. %v", err))
		shutdownCallback()
	}
}

func (ad *Adapter) Stop(ctx context.Context, timeout time.Duration) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	if ad.server != nil {
		logging.Global.Info(ctx, "shutting down HTTP server")
		if err := ad.server.Shutdown(ctx); err != nil {
			logging.Global.Error(ctx, err.Error())
		}
	}
	logging.Global.Info(ctx, "HTTP server shutdown complete")
}

// Health
// @Summary Service Health
// @Description Service Health
// @Tags api
// @Accept json
// @Produce json
// @Success 200 {object} health.Health
// @Header 200 {string} Token "qwerty"
// @Failure 503 {object} string
// @Router /api/health [get]
func (ad *Adapter) Health(c echo.Context) error {
	ctx := c.Request().Context()
	isHealthy := ad.api.IsHealthy(ctx)

	if !isHealthy {
		return c.JSON(http.StatusInternalServerError, "Service is not ready to accept requests")
	}

	host, _ := os.Hostname()
	health := health.Health{
		Version:           config.BuildVersion,
		Healthy:           isHealthy,
		Host:              host,
		ApplicationConfig: config.Configuration.ApplicationConfiguration(),
		InfraConfig:       config.Configuration.InfrastructureConfiguration(),
	}

	return c.JSON(http.StatusOK, health)
}

// Create
// @Summary Create User
// @Description Create user
// @Tags db
// @Success 200 {object} db.ID
// @Failure 500 {object} error
// @Accept json
// @Produce json
// @Param x-trace-id header string false "x-trace-id"
// @Param request body db.CreateRequest true "Create User Request"
// @Router /api/db/create [put]
func (ad *Adapter) CreateDB(c echo.Context) error {
	start := time.Now()
	ctx := c.Request().Context()

	ad.setTracingXTraceId(ctx)

	var user db.CreateRequest
	err := json.NewDecoder(c.Request().Body).Decode(&user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	resp, err := ad.api.DBCreate(ctx, &user)
	ad.updateMetrics("/api/db/create", err, start)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, resp)
}

// Get
// @Summary Get User by ID
// @Description Get user
// @Tags db
// @Success 200 {object} db.User
// @Failure 500 {object} error
// @Accept json
// @Produce json
// @Param x-trace-id header string false "x-trace-id"
// @Param id path string true "User ID"
// @Router /api/db/get [get]
func (ad *Adapter) GetDB(c echo.Context) error {
	start := time.Now()
	ctx := c.Request().Context()

	ad.setTracingXTraceId(ctx)

	id := c.QueryParam("id")

	resp, err := ad.api.DBGet(ctx, uuid.MustParse(id))
	ad.updateMetrics("/api/db/get", err, start)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	if resp != nil {
		return c.JSON(http.StatusOK, resp)
	}
	return c.JSON(http.StatusNotFound, fmt.Sprintf("user with id %d was not found", id))
}

// Get All
// @Summary Get All Users
// @Description Get all users
// @Tags db
// @Success 200 {object} []db.User
// @Failure 500 {object} error
// @Accept json
// @Produce json
// @Param x-trace-id header string false "x-trace-id"
// @Router /api/db/get-all [get]
func (ad *Adapter) GetAllDB(c echo.Context) error {
	start := time.Now()
	ctx := c.Request().Context()

	ad.setTracingXTraceId(ctx)

	resp, err := ad.api.DBGetAll(ctx)
	ad.updateMetrics("/api/db/get-all", err, start)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	if resp != nil {
		return c.JSON(http.StatusOK, resp)
	}
	return c.JSON(http.StatusNotFound, fmt.Sprintf("no users in the database"))
}

// Update
// @Summary Update User by id
// @Description Update user name and surname by provided id
// @Tags db
// @Success 200 {object} db.ID
// @Failure 500 {object} error
// @Accept json
// @Produce json
// @Param x-trace-id header string false "x-trace-id"
// @Param request body db.UpdateRequest true "Update User Request"
// @Router /api/db/update [post]
func (ad *Adapter) UpdateDB(c echo.Context) error {
	start := time.Now()
	ctx := c.Request().Context()

	ad.setTracingXTraceId(ctx)

	var req db.UpdateRequest
	err := json.NewDecoder(c.Request().Body).Decode(&req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	resp, err := ad.api.DBUpdate(ctx, &req)
	ad.updateMetrics("/api/db/update", err, start)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, resp)
}

// Delete
// @Summary Delete User by id
// @Description Delete user name and surname by provided id
// @Tags db
// @Success 200 {object} bool
// @Failure 500 {object} error
// @Accept json
// @Produce json
// @Param x-trace-id header string false "x-trace-id"
// @Param request body db.DeleteRequest true "Delete User Request"
// @Router /api/db/delete [delete]
func (ad *Adapter) DeleteDB(c echo.Context) error {
	start := time.Now()
	ctx := c.Request().Context()

	ad.setTracingXTraceId(ctx)

	var req db.DeleteRequest
	err := json.NewDecoder(c.Request().Body).Decode(&req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	resp, err := ad.api.DBDelete(ctx, req.Id)
	ad.updateMetrics("/api/db/delete", err, start)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, resp)
}
