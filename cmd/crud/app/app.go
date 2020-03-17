package app

import (
	"errors"
	"net/http"

	"github.com/mux/pkg/crud/services/burgers"

	"github.com/jackc/pgx/v4/pgxpool"
)

// описание сервиса, который хранит зависимости и выполняет работу
type server struct {
	pool          *pgxpool.Pool
	router        http.Handler
	burgersSvc    *burgers.ServiceBurgers
	templatesPath string
	assetsPath    string
}

// ErrRouteCantBeNil ...
var ErrRouteCantBeNil = errors.New("router can't be nil")

// ErrPoolCantBeNil ...
var ErrPoolCantBeNil = errors.New("pool can't be nil")

// ErrBurgerSvcCantBeNil ...
var ErrBurgerSvcCantBeNil = errors.New("burgersSvc can't be nil")

// ErrTamplateCantBeNil ...
var ErrTamplateCantBeNil = errors.New("templatesPath can't be empty")

// ErrAssetsPathCantBeNil ...
var ErrAssetsPathCantBeNil = errors.New("assetsPath can't be empty")

// NewServer ...
func NewServer(router http.Handler, pool *pgxpool.Pool, burgersSvc *burgers.ServiceBurgers, templatesPath string, assetsPath string) *server {
	if router == nil {
		panic(ErrRouteCantBeNil)
	}
	if pool == nil {
		panic(ErrPoolCantBeNil)
	}
	if burgersSvc == nil {
		panic(ErrBurgerSvcCantBeNil)
	}
	if templatesPath == "" {
		panic(ErrTamplateCantBeNil)
	}
	if assetsPath == "" {
		panic(ErrAssetsPathCantBeNil)
	}

	return &server{
		router:        router,
		pool:          pool,
		burgersSvc:    burgersSvc,
		templatesPath: templatesPath,
		assetsPath:    assetsPath,
	}
}

func (receiver *server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	receiver.router.ServeHTTP(writer, request)
}
