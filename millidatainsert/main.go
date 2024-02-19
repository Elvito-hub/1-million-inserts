package main

import (
	"millidatainsert/db"
	"millidatainsert/routes"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	db.InitDb()
	server := gin.Default()

	server.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:3009", "http://localhost:3002", "http://localhost:8081", "http://aquilaspeed-env.eba-fkugkv2i.eu-north-1.elasticbeanstalk.com"},
		AllowMethods: []string{"PUT", "PATCH", "POST", "OPTIONS"},
		// AllowHeaders:     []string{"Origin"},
		AllowHeaders:  []string{"X-Auth-Key", "X-Auth-Secret", "Content-Type"},
		ExposeHeaders: []string{"Content-Length"},
		AllowOriginFunc: func(origin string) bool {
			return origin == "http://localhost:3000"
		},
		MaxAge: 12 * time.Hour,
	}))

	server.POST("/posttodb", routes.HandlePostMilliData)
	server.GET("/fakecsv", routes.GenerateCsv)
	server.POST("/postmillirequests", routes.HandleMilliRequest)
	server.Run(":2030")
}
