package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/chack-check/chats-service/api/v1/graph"
	"github.com/chack-check/chats-service/api/v1/models"
	"github.com/chack-check/chats-service/api/v1/services"
	"github.com/chack-check/chats-service/database"
	"github.com/go-chi/chi"
)

const defaultPort = "8000"

func main() {
	log.SetFlags(log.Lshortfile)
	port := os.Getenv("PORT")

	if port == "" {
		port = defaultPort
	}

	database.DB.AutoMigrate(&models.Chat{})

	router := chi.NewRouter()

	router.Use(services.UserMiddleware)

	srvV1 := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))
	router.Handle("/api/v1/chats", playground.Handler("GraphQL playground", "/api/v1/chats/query"))

	router.Handle("/api/v1/chats/query", srvV1)

	log.Printf("Server has started on http://0.0.0.0:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
