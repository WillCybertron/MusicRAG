package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/One-Frequency/MusicRAG/backend/graph"
	"github.com/One-Frequency/MusicRAG/backend/graph/generated"
	"github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()

    // Health check route
    r.GET("/health", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{"status": "ok"})
    })

    // GraphQL handler
    gqlHandler := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))
    r.POST("/graphql", func(c *gin.Context) { gqlHandler.ServeHTTP(c.Writer, c.Request) })
    r.GET("/", func(c *gin.Context) {
        playground.Handler("GraphQL", "/graphql").ServeHTTP(c.Writer, c.Request)
    })

    port := os.Getenv("PORT")
    if port == "" { port = "8080" }
    log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
    if err := r.Run(":" + port); err != nil {
        log.Fatal(err)
    }
}