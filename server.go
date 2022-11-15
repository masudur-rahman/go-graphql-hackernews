package main

import (
	"log"
	"net/http"
	"os"

	"github.com/masudur-rahman/hackernews/graph"
	hackernews "github.com/masudur-rahman/hackernews/graph/generated"
	database "github.com/masudur-rahman/hackernews/internal/pkg/db/migrations/mysql"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi/v5"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	router := chi.NewRouter()

	database.InitDB()
	defer database.CloseDB()

	server := handler.NewDefaultServer(
		hackernews.NewExecutableSchema(
			hackernews.Config{
				Resolvers: &graph.Resolver{},
			},
		),
	)

	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", server)

	log.Printf("Connect to http://localhost:%s/ for GraphQL playground\n", port)
	log.Fatalln(http.ListenAndServe(":"+port, router))
}
