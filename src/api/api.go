package api

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/gorilla/mux"
	ctrls "github.com/m-faried/api/controllers"
	models "github.com/m-faried/api/models"
	httpRouters "github.com/m-faried/api/presentations/routers"
	rpc "github.com/m-faried/api/presentations/rpc"
	"google.golang.org/grpc"
)

type App struct {
	Port       string
	PortGrpc   string
	DBPath     string
	database   *models.EcommerceDal
	router     *mux.Router
	grpcServer *rpc.GrpcEcommerceServer
}

func (a *App) Initialize() {
	// Config validations
	if a.DBPath == "" || a.Port == "" {
		panic("App configuration is missing")
	}

	// Creating & initializaing the datasource
	var source models.EcommerceDal
	source.Initialize(a.DBPath)
	a.database = &source

	// Creating & initializaing the router
	a.initHttpRouter()

	// Initializing Grpc
	a.initGrpc()
}

func (a *App) RunHttp() {
	// You can do the following as well.
	// http.Handle("/", a.router)
	fmt.Println("Server started and listening on the port", a.Port)
	log.Fatal(http.ListenAndServe(a.Port, a.router))
}

func (a *App) RunGrpc() {
	lis, err := net.Listen("tcp", a.PortGrpc)
	if err != nil {
		log.Fatalf("Cannot create listener: %s", err)
	} else {
		fmt.Println("gRPC started and listending on the port", a.PortGrpc)
	}
	serverRegistrar := grpc.NewServer()
	rpc.RegisterECommerceServer(serverRegistrar, a.grpcServer)
	err = serverRegistrar.Serve(lis)
	fmt.Println("The error", err)
}

func (a *App) initHttpRouter() {

	router := mux.NewRouter()

	router.HandleFunc("/", healthCheck).Methods("GET")

	pc := ctrls.NewProductsController(a.database)
	pr := httpRouters.NewProductsRouter(router, pc)
	pr.InitRoutes()

	oc := ctrls.NewOrdersController(a.database)
	or := httpRouters.NewOrdersRouter(router, oc)
	or.InitRoutes()

	a.router = router
}

func (a *App) initGrpc() {
	pc := ctrls.NewProductsController(a.database)
	oc := ctrls.NewOrdersController(a.database)
	ps := rpc.NewGrpcEcommerceServer(pc, oc)
	a.grpcServer = ps
}

//////////////////// Helper Functions

func healthCheck(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(res, "Hello World")
}
