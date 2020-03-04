package main

// package
// import
// var + type
// method + function

import (
	"context"
	"crud/cmd/crud/app"
	"crud/pkg/crud/services/burgers"
	"flag"
	"log"
	"net"
	"net/http"
	"path/filepath"

	"github.com/jackc/pgx/v4/pgxpool"
	// "github.com/julienschmidt/httprouter"
)

var (
	host = flag.String("host", "0.0.0.0", "Server host")
	port = flag.String("port", "9999", "Server port")
	dsn  = flag.String("dsn", "postgres://user:pass@localhost:5434/app", "Postgres DSN")
)

func main() {
	flag.Parse()
	addr := net.JoinHostPort(*host, *port)
	start(addr, *dsn)
}

func start(addr string, dsn string) {
	// router := app.NewExactMux()
	router := app.NewPathResolver()
	// Context: <-
	pool, err := pgxpool.Connect(context.Background(), dsn)
	if err != nil {
		panic(err)
	}
	burgersSvc := burgers.NewBurgersSvc(pool)
	server := app.NewServer(
		router,
		pool,
		burgersSvc, // DI + Containers
		filepath.Join("web", "templates"),
		filepath.Join("web", "assets"),
	)
	// server.InitRoutes()
	server.InitRoutesPath()

	log.Println("Server starting ...")
	// server'ы должны работать "вечно"
	panic(http.ListenAndServe(addr, server)) // поднимает сервер на определённом адресе и обрабатывает запросы
}
