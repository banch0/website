package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"net"
	"net/http"
	"path/filepath"

	"github.com/banch0/mux/cmd/website/app"
	"github.com/banch0/mux/pkg/website/services/burgers"

	"github.com/jackc/pgx/v4/pgxpool"
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

// ErrConnectDBPool ...
var ErrConnectDBPool = errors.New("Can't create connect to db pool")

func start(addr string, dsn string) {
	// router := app.NewExactMux()
	router := app.NewPathResolver()
	pool, err := pgxpool.Connect(context.Background(), dsn)
	if err != nil {
		panic(ErrConnectDBPool)
	}

	burgersSvc := burgers.NewBurgersSvc(pool)
	server := app.NewServer(
		router,
		pool,
		burgersSvc,
		filepath.Join("web", "templates"),
		filepath.Join("web", "assets"),
	)
	// server.InitRoutes()
	server.InitRoutesPath()

	log.Println("Server starting ...")
	panic(http.ListenAndServe(addr, server))
}
