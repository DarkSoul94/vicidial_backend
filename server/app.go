package server

import (
	"context"
	"log"
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
)

// App ...
type App struct {
	vicidial_backendUC vicidial_backend.Usecase
	httpServer         *http.Server
	httpClient         helper.Helper
}

// NewApp ...
func NewApp() *App {
	clint := helperUC.NewHelper()
	uc := vicidial_backendusecase.NewUsecase(clint)
	return &App{
		vicidial_backendUC: uc,
	}
}

// Run run vicidial_backendlication
func (a *App) Run(port string) error {
	router := gin.Default()
	router.Use(
		gin.Recovery(),
		gin.Logger(),
	)

	vicidial_backendhttp.RegisterHTTPEndpoints(router, a.vicidial_backendUC)

	a.httpServer = &http.Server{
		Addr:           ":" + port,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		if err := a.httpServer.ListenAndServe(); err != nil {
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
