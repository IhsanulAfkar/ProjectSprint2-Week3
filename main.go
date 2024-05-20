package main

import (
	"Week3/db"
	"Week3/routes"
	"context"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
)
func main(){
	ctx := context.Background()
	db.Init(ctx)
	gin.SetMode(gin.ReleaseMode)
	r := routes.Init()
	r.Run(":8080")
}