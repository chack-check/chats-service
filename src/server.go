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
	"github.com/chack-check/chats-service/settings"
	"github.com/go-chi/chi"
)

func main() {
	log.SetFlags(log.Lshortfile)

	database.DB.AutoMigrate(&models.Chat{}, &models.Message{}, &models.Reaction{})

	router := chi.NewRouter()

	router.Use(services.UserMiddleware)

	srvV1 := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))
	router.Handle("/api/v1/chats", playground.Handler("GraphQL playground", "/api/v1/chats/query"))

	router.Handle("/api/v1/chats/query", srvV1)

	log.Printf("Server has started on http://0.0.0.0:%d", settings.Settings.PORT)
	listen := fmt.Sprintf(":%d", settings.Settings.PORT)
	log.Fatal(http.ListenAndServe(listen, router))
}
