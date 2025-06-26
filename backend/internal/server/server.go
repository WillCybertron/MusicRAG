package server

import (
	"log"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/One-Frequency/MusicRAG/backend/graph"
	"github.com/One-Frequency/MusicRAG/backend/graph/generated"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// NewServer creates and configures a Gin engine with GraphQL, Playground, health check, and CORS.
func NewServer() *gin.Engine {
    // Set Gin mode from ENV, fallback to Debug
    mode := os.Getenv("GIN_MODE")
    if mode == "" {
        mode = gin.DebugMode
    }
    gin.SetMode(mode)

    r := gin.Default()

    // Trusted proxies from ENV, fallback to localhost
    trustedProxies := os.Getenv("GIN_TRUSTED_PROXIES")
    if trustedProxies == "" {
        trustedProxies = "127.0.0.1"
    }
    if err := r.SetTrustedProxies([]string{trustedProxies}); err != nil {
        log.Fatalf("Could not set trusted proxies: %v", err)
    }

    // CORS middleware (adjust AllowOrigins for production)
    r.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"http://localhost:3000"},
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
    }))

    // Health check endpoint
    r.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{"status": "ok"})
    })

    // GraphQL handler (POST)
    graphqlHandler := handler.NewDefaultServer(
        generated.NewExecutableSchema(generated.Config{
            Resolvers: &graph.Resolver{},
        }),
    )
    r.POST("/graphql", gin.WrapH(graphqlHandler))

    // Playground handler (GET)
    playgroundHandler := playground.Handler("GraphQL Playground", "/graphql")
    r.GET("/", gin.WrapH(playgroundHandler))

    return r
}