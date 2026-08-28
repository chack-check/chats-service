package api

import (
	"fmt"
	"log"
	"net/http"

	"chats-service/configs"
	"chats-service/internal/infrastructure/api/graph"
	"chats-service/internal/infrastructure/api/middlewares"
	"chats-service/internal/infrastructure/database"
	"chats-service/internal/infrastructure/rabbit"
	"chats-service/internal/infrastructure/redisdb"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi"
)

func RunApi() {
	defer rabbit.EventsRabbitConnection.Close()
	defer redisdb.RedisConnection.Close()

	configuration := configs.GetAPIConfiguration()
	database.DatabaseConnection.AutoMigrate(&database.Chat{}, &database.Message{}, &database.SavedFile{}, database.Reaction{})

	router := chi.NewRouter()

	router.Use(middlewares.UserMiddleware)
	router.Use(middlewares.CorsMiddleware)

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))

	router.Handle("/api/v1/chats", playground.Handler("GraphQL playground", "/api/v1/chats/query"))
	router.Handle("/api/v1/chats/query", srv)

	listen := fmt.Sprintf(":%d", configuration.Port)
	log.Fatal(http.ListenAndServe(listen, router))
}
