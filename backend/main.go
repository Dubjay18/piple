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
	v1.POST("/employees", h.Create_employee)
	v1.PATCH("/employees/:id", h.Update_employee)
	v1.DELETE("/employees/:id", h.Delete_employee)

	log.Printf("server running on port %s", PORT)
	router.Run(PORT)

}
