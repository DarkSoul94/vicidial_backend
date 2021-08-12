package server

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/DarkSoul94/vicidial_backend/helper"
	helperUC "github.com/DarkSoul94/vicidial_backend/helper/usecase"
	"github.com/DarkSoul94/vicidial_backend/vicidial_backend"
	vicidial_backendhttp "github.com/DarkSoul94/vicidial_backend/vicidial_backend/delivery/http"
	vicidial_backendusecase "github.com/DarkSoul94/vicidial_backend/vicidial_backend/usecase"
	"github.com/gin-gonic/gin"
	_ "github.com/golang-migrate/migrate/v4/source/file" // required
	"github.com/spf13/viper"
)

// App ...
type App struct {
	vicidial_backendUC vicidial_backend.Usecase
	httpServer         *http.Server
	httpClient         helper.Helper
	socket             string
}

// NewApp ...
func NewApp() *App {
	clint := helperUC.NewHelper()
	uc := vicidial_backendusecase.NewUsecase(clint)
	return &App{
		vicidial_backendUC: uc,
		socket:             viper.GetString("app.socket"),
	}
}

// Run run vicidial_backendlication
func (a *App) Run(port string) error {
	router := gin.New()
	if viper.GetBool("app.release") {
		gin.SetMode(gin.ReleaseMode)
	} else {
		router.Use(gin.Logger())
	}
	router.Use(gin.Recovery())
	apiRouter := router.Group("/api/v1")

	vicidial_backendhttp.RegisterHTTPEndpoints(apiRouter, a.vicidial_backendUC)

	a.httpServer = &http.Server{
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if err := os.RemoveAll(a.socket); err != nil {
		log.Fatal(err)
	}
	unixListener, err := net.Listen("unix", a.socket)
	if err != nil {
		log.Fatal(err)
	}
	os.Chmod(a.socket, 0664)
	defer unixListener.Close()

	go func() {
		if err := a.httpServer.Serve(unixListener); err != nil {
			log.Fatalf("Failed to listen and serve: %+v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Interrupt)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	return a.httpServer.Shutdown(ctx)
}
