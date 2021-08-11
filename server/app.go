package server

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/DarkSoul94/vicidial_backend/vicidial_backend"
	vicidial_backendhttp "github.com/DarkSoul94/vicidial_backend/vicidial_backend/delivery/http"
	vicidial_backendusecase "github.com/DarkSoul94/vicidial_backend/vicidial_backend/usecase"
	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file" // required
	"github.com/spf13/viper"
)

// App ...
type App struct {
	vicidial_backendUC vicidial_backend.Usecase
	httpServer         *http.Server
}

// NewApp ...
func NewApp() *App {
	uc := vicidial_backendusecase.NewUsecase()
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

func initDB() *sql.DB {
	dbString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s",
		viper.GetString("vicidial_backend.db.login"),
		viper.GetString("vicidial_backend.db.pass"),
		viper.GetString("vicidial_backend.db.host"),
		viper.GetString("vicidial_backend.db.port"),
		viper.GetString("vicidial_backend.db.name"),
		viper.GetString("vicidial_backend.db.args"),
	)
	db, err := sql.Open(
		"mysql",
		dbString,
	)
	if err != nil {
		panic(err)
	}
	runMigrations(db)
	return db
}

func runMigrations(db *sql.DB) {
	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		panic(err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		viper.GetString("vicidial_backend.db.name"),
		driver)
	if err != nil {
		panic(err)
	}
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange && err != migrate.ErrNilVersion {
		fmt.Println(err)
	}
}
