package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/chack-check/chats-service/infrastructure/api/graph"
	"github.com/chack-check/chats-service/infrastructure/api/middlewares"
	"github.com/chack-check/chats-service/infrastructure/api/settings"
	"github.com/chack-check/chats-service/infrastructure/database"
	"github.com/chack-check/chats-service/infrastructure/rabbit"
	"github.com/chack-check/chats-service/infrastructure/redisdb"
	"github.com/go-chi/chi"
)

func RunApi() {
	defer rabbit.EventsRabbitConnection.Close()
	defer redisdb.RedisConnection.Close()

	database.DatabaseConnection.AutoMigrate(&database.Chat{}, &database.Message{}, &database.SavedFile{}, database.Reaction{})

	router := chi.NewRouter()

	router.Use(middlewares.UserMiddleware)
	router.Use(middlewares.CorsMiddleware)

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))

	router.Handle("/api/v1/chats", playground.Handler("GraphQL playground", "/api/v1/chats/query"))
	router.Handle("/api/v1/chats/query", srv)

	listen := fmt.Sprintf(":%d", settings.Settings.APP_PORT)
	log.Fatal(http.ListenAndServe(listen, router))
}
