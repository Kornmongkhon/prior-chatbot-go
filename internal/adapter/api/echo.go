package api

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"prior-chat-bot/configs"
	"prior-chat-bot/internal/adapter/api/handler"
	"prior-chat-bot/internal/adapter/repository"
	"prior-chat-bot/internal/core"
	"prior-chat-bot/internal/core/authentication"
	"prior-chat-bot/internal/core/domain"
	"prior-chat-bot/internal/core/port"
	"prior-chat-bot/internal/core/service"
)

type EchoContext struct {
	echo.Context
	repo port.MyRepo
}

func NewEchoContext(e echo.Context, repo port.MyRepo) *EchoContext {
	return &EchoContext{Context: e, repo: repo}
}

func StartEchoServer() {
	fmt.Println("Start Echo Server")

	// Initialize Echo
	e := echo.New()
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Add context handling and error handling middleware
			err := next(c)
			if err != nil {
				log.Printf("Error occurred in %s %s: %v", c.Request().Method, c.Path(), err)
				switch err.(type) {
				case *echo.HTTPError:
					httpError := err.(*echo.HTTPError)
					return c.JSON(httpError.Code, domain.ResponseT[interface{}]{
						Code:    "E9999",
						Message: "The system has a problem. Please contact the system administrator.",
					})
				default:
					return c.JSON(http.StatusInternalServerError, domain.ResponseT[interface{}]{
						Code:    "E9999",
						Message: "The system has a problem. Please contact the system administrator.",
					})
				}
			}
			return nil
		}
	})

	// Load configurations
	if err := configs.Init(""); err != nil {
		log.Fatal("Configuration initialization failed: ", err)
	}
	cfg := configs.GetConfig()

	configs.SetTimeZone(cfg.Server.TimeZone)

	// Database connection
	db := core.InitDb(cfg.DB.Host, cfg.DB.Port, cfg.Secrets.DbUsername, cfg.Secrets.DbPassword, cfg.DB.Database)
	defer func() {
		if err := db.Close(); err != nil {
			log.Fatal("Error closing db connection:", err)
		}
	}()
	// JWT configuration
	log.Println("NewJwtTokenProvider")
	jwtTokenProvider := authentication.NewJwtTokenProvider(cfg.JWT.Secret, cfg.JWT.ExpirationAccessToken, cfg.JWT.ExpirationRefreshToken)
	log.Println("InterceptorFilter")
	jwtProvider = jwtTokenProvider

	// Initialize repository and service
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			repo := &repository.AuthRepository{DB: db}
			c.Set("repo", repo) // Store repo in context
			return next(c)
		}
	})
	authService := service.NewAuthService(&repository.AuthRepository{DB: db}, jwtTokenProvider)
	log.Println(authService)
	e.Use(InterceptorFilter)
	// Setup routes
	apiV1 := e.Group("/prior_chatbot_api/api/v1")
	apiV1.GET("/health/check", convertEchoHandler(handler.ExecuteHandlerHealthCheck))
	apiV1.POST("/auth/sign-in", convertEchoHandler(func(c port.MyServer) {
		handler.ExecuteHandlerSignIn(c, &cfg)
	}))
	apiV1.GET("/auth/me", convertEchoHandler(func(c port.MyServer) {
		handler.ExecuteHandlerMe(c, &cfg)
	}))
	apiV1.POST("/auth/regenerate-tokens", convertEchoHandler(func(c port.MyServer) {
		handler.ExecuteHandlerRegenerateToken(c, &cfg)
	}))
	apiV1.POST("/auth/sign-up", convertEchoHandler(func(c port.MyServer) {
		handler.ExecuteHandlerSignUp(c, &cfg)
	}))

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}

func (e *EchoContext) GetContext() context.Context {
	return e.Request().Context()
}

func (e *EchoContext) GetRequest() *http.Request {
	return e.Request()
}

func (e *EchoContext) GetRepo() port.MyRepo {
	return e.repo
}

func (e *EchoContext) BindRequest(v interface{}) error {
	if err := e.Bind(&v); err != nil {
		return e.String(http.StatusBadRequest, "bad request")
	}
	return nil
}

func (e *EchoContext) ToResponse(code int, statusCode string, msg interface{}, data interface{}) error {
	response := port.Response{
		Code:    statusCode,
		Message: msg,
		Data:    data,
	}
	return e.JSON(code, response)
}

func convertEchoHandler(handler func(port.MyServer)) echo.HandlerFunc {
	return func(c echo.Context) error {
		repo := c.Get("repo").(*repository.AuthRepository)
		echoCtx := NewEchoContext(c, repo)
		handler(echoCtx)
		return nil
	}
}
