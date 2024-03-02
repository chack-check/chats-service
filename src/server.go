package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/chack-check/chats-service/api/v1/graph"
	"github.com/chack-check/chats-service/api/v1/models"
	"github.com/chack-check/chats-service/api/v1/services"
	"github.com/chack-check/chats-service/database"
	grpcserver "github.com/chack-check/chats-service/grpc_server"
	"github.com/chack-check/chats-service/middlewares"
	"github.com/chack-check/chats-service/rabbit"
	"github.com/chack-check/chats-service/settings"
	"github.com/chack-check/chats-service/ws"
	"github.com/go-chi/chi"
)

func main() {
	defer rabbit.EventsRabbitConnection.Close()
	log.SetFlags(log.Lshortfile)

	database.DB.AutoMigrate(&models.Chat{}, &models.Message{}, &models.Reaction{})

	router := chi.NewRouter()

	router.Use(services.UserMiddleware)
	router.Use(middlewares.CorsMiddleware)

	srvV1 := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))

	router.Handle("/api/v1/chats", playground.Handler("GraphQL playground", "/api/v1/chats/query"))
	router.Handle("/api/v1/chats/query", srvV1)
	router.HandleFunc("/api/v1/chats/ws", ws.WsHandler)

	go grpcserver.StartServer(settings.Settings.GRPC_SERVER_HOST, settings.Settings.GRPC_SERVER_PORT)

	log.Printf("Server has started on http://0.0.0.0:%d", settings.Settings.PORT)
	listen := fmt.Sprintf(":%d", settings.Settings.PORT)
	log.Fatal(http.ListenAndServe(listen, router))
}
