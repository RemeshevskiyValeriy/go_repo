package main

import (
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"

	"github.com/go-chi/chi/v5"

	"example.com/pz12-rest-vs-graphql/graph"
	"example.com/pz12-rest-vs-graphql/internal/rest"
)

func main() {

	r := chi.NewRouter()

	// REST
	r.Get("/v1/tasks", rest.GetTasks)
	r.Get("/v1/tasks/{id}", rest.GetTask)
	r.Post("/v1/tasks", rest.CreateTask)
	r.Patch("/v1/tasks/{id}", rest.UpdateTask)

	// GraphQL
	srv := handler.NewDefaultServer(
		graph.NewExecutableSchema(
			graph.Config{
				Resolvers: &graph.Resolver{},
			},
		),
	)

	r.Handle("/", playground.Handler("GraphQL playground", "/query"))
	r.Handle("/query", srv)

	log.Println("server started at :8080")

	log.Fatal(http.ListenAndServe(":8080", r))
}