package main

import (
	"database/sql"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	graphql_handler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/devfullcycle/20-CleanArch/internal/event/handler"
	"github.com/devfullcycle/20-CleanArch/internal/infra/graph"
	"github.com/devfullcycle/20-CleanArch/internal/infra/grpc/pb"
	"github.com/devfullcycle/20-CleanArch/internal/infra/grpc/service"
	"github.com/devfullcycle/20-CleanArch/internal/infra/web/webserver"
	"github.com/devfullcycle/20-CleanArch/pkg/events"
	"github.com/streadway/amqp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	// mysql
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// Load configuration from environment variables
	dbDriver := "mysql"
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	webServerPort := os.Getenv("WEB_SERVER_PORT")
	grpcServerPort := os.Getenv("GRPC_SERVER_PORT")
	graphQLServerPort := os.Getenv("GRAPHQL_SERVER_PORT")
	rabbitMQHost := "rabbitmq"

	// Retry logic for database connection
	var db *sql.DB
	var err error
	retryInterval := 3 * time.Second
	timeout := 30 * time.Second
	startTime := time.Now()
	//print

	for {
		
		fmt.Printf("Trying to connect to database: %s:%s@tcp(%s:%s)/%s\n", dbUser, dbPassword, dbHost, dbPort, dbName)
		db, err = sql.Open(dbDriver, fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPassword, dbHost, dbPort, dbName))
		if err == nil {
			err = db.Ping()
			if err == nil {
				break
			}
		}

		if time.Since(startTime) > timeout {
			panic(fmt.Sprintf("Failed to connect to database after %v: %v", timeout, err))
		}

		fmt.Println("Retrying database connection...")
		time.Sleep(retryInterval)
	}

	defer db.Close()
	fmt.Println("Connected to the database successfully.")

	rabbitMQChannel := getRabbitMQChannel(rabbitMQHost)

	eventDispatcher := events.NewEventDispatcher()
	eventDispatcher.Register("OrderCreated", &handler.OrderCreatedHandler{
		RabbitMQChannel: rabbitMQChannel,
	})

	createOrderUseCase := NewCreateOrderUseCase(db, eventDispatcher)
	ListAllOrdersUseCase := NewListAllOrdersUseCase(db, eventDispatcher)

	//Web server
	webserver := webserver.NewWebServer(webServerPort)
	webOrderHandler := NewWebOrderHandler(db, eventDispatcher)
	webserver.AddHandler(http.MethodPost, "/order", webOrderHandler.Create)
	webserver.AddHandler(http.MethodGet, "/order", webOrderHandler.ListAll)
	fmt.Println("Starting web server on port", webServerPort)
	go webserver.Start()

	//GRPC server
	grpcServer := grpc.NewServer()
	orderService := service.NewOrderService(*createOrderUseCase, *ListAllOrdersUseCase)
	pb.RegisterOrderServiceServer(grpcServer, orderService)
	reflection.Register(grpcServer)

	fmt.Println("Starting gRPC server on port", grpcServerPort)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcServerPort))
	if err != nil {
		panic(err)
	}
	go grpcServer.Serve(lis)

	//GraphQL Server
	srv := graphql_handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
		CreateOrderUseCase:    *createOrderUseCase,
		ListAllOrdersUseCase:  *ListAllOrdersUseCase,
	}}))
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	fmt.Println("Starting GraphQL server on port", graphQLServerPort)
	http.ListenAndServe(":"+graphQLServerPort, nil)
}

func getRabbitMQChannel(rabbitMQHost string) *amqp.Channel {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://guest:guest@%s:5672/", rabbitMQHost))
	if err != nil {
		panic(err)
	}
	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	return ch
}
