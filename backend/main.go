package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	db "github.com/valentineejk/piple/db/postgres"
	"github.com/valentineejk/piple/internal/handler"
)

func main() {

	
	// Load .env into the process environment. Not fatal if it's missing —
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, relying on existing environment")
	}
	dq, pg := db.Connection()
	defer pg.Close()

	h := handler.New(dq)

		PORT := ":8080"

		router := gin.Default()

		v1 := router.Group("/api/v1")
		v1.GET("/health", h.HealthCheck)

		users := v1.Group("/users") // , handler.AuthMiddleware(), handler.AdminOnly() 
	    users.POST("", h.CreateUser)
	    users.PATCH("/:id", h.UpdateUser)
	    users.DELETE("/:id", h.DeleteUser)

		log.Printf("server running on port %s", PORT)
		router.Run(PORT)

}
