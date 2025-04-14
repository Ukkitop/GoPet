package main

import (
	"context"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/soheilhy/cmux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"realAPI/internal/database"
	"realAPI/internal/handlers"
	"realAPI/internal/routes"
	"realAPI/packages/proto/pb"
)

const mongodb_uri = "mongodb+srv://user:user@godatabase.jia0iua.mongodb.net/?retryWrites=true&w=majority&appName=GoDatabase"

func main() {
	log.SetReportCaller(true)

	dbErr := database.SetupDatabase()
	if dbErr != nil {
		log.Error(dbErr)
	}

	clientOpts := options.Client().ApplyURI(mongodb_uri).SetMaxPoolSize(20)
	pooledClient, err := mongo.Connect(context.TODO(), clientOpts)

	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := pooledClient.Disconnect(context.TODO()); err != nil {
			log.Fatal(err)
		}
	}()

	mongoClient := database.NewServer(pooledClient)

	listener, listenerErr := net.Listen("tcp", ":52423")

	if listenerErr != nil {
		log.Fatal(listenerErr)
	}

	m := cmux.New(listener)
	grpcListener := m.Match(cmux.HTTP2HeaderField("content-type", "application/grpc"))
	httpListener := m.Match(cmux.HTTP1Fast())

	g := new(errgroup.Group)
	g.Go(func() error { return grpcServe(grpcListener) })
	g.Go(func() error { return httpServe(httpListener, mongoClient) })
	g.Go(func() error { return m.Serve() })

	log.Println("run server:", g.Wait())
}

func grpcServe(l net.Listener) error {
	grpcServer := grpc.NewServer()
	pb.RegisterBatchServiceServer(grpcServer, &handlers.Server{})

	return grpcServer.Serve(l)
}

func httpServe(l net.Listener, client *database.Server) error {
	r := mux.NewRouter()

	routes.Routes(r, client)
	httpServer := &http.Server{Handler: r}

	return httpServer.Serve(l)
}
