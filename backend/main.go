package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	database "github.com/valentineejk/piple/database"
	"github.com/valentineejk/piple/handler"
)

func main() {

	// Load .env into the process environment. Not fatal if it's missing —
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, relying on existing environment")
	}

	PORT := ":8000"

	q, p := database.Connection()
	defer p.Close()

	h := handler.New(q)

	r := gin.Default()

	v1 := r.Group("/api/v1")

	v1.POST("/employees", h.Create_employee)
	v1.PATCH("/employees/:id", h.Update_employee)
	v1.DELETE("/employees/:id", h.Delete_employee)
	v1.GET("/health", h.HealthCheck)

	r.Run(PORT)

}
