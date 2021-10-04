package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"mytest/pkg/user_balance/config"
	"mytest/pkg/user_balance/server"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	logger := zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		os.Stdout,
		zap.DebugLevel,
	))
	cfg := config.NewUserBalanceConfig()
	ctx := context.Background()
	db, err := pgxpool.Connect(ctx, cfg.DbConn)
	if err != nil {
		fmt.Println(err)
		return
	}
	userService := server.NewServ(db, logger)
	userHandler := server.NewGin(userService)
	userServer := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Handler: userHandler,
	}
	go func() {
		// service connections
		if err := userServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		}
	}()
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscanll.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	// catching ctx.Done(). timeout of 5 seconds.
	<-ctx.Done()

	cancel()
}
