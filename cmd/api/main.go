package main

import (
	"context"
	"fmt"
	"http_server/internal/config"
	"http_server/internal/transport/http"
	"http_server/internal/transport/http/handler"
)

func main() {

	ctx := context.Background()

	pool, err := config.CheckConnection(ctx)
	if err != nil {
		panic(err)
	}

	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		panic(err)
	}

	dbService, err := config.InitPostgresEnv(ctx, pool)
	if err != nil {
		panic(err)
	}

	httpHandlers := handler.NewHTTPHandlers(dbService)

	httpServer := http.NewHTTPServer(httpHandlers)
	if err := httpServer.StartServer(); err != nil {
		fmt.Printf("failed start http server\n")
		return
	}
}
