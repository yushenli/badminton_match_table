package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/yushenli/badminton_match_table/web/lib/config"
	"github.com/yushenli/badminton_match_table/web/www/controller"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func setupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/", controller.RenderIndex)
	r.Group("/v1")
	{
		/*v1.GET("book", Controllers.ListBook)
		v1.POST("book", Controllers.AddNewBook)
		v1.GET("book/:id", Controllers.GetOneBook)
		v1.PUT("book/:id", Controllers.PutOneBook)
		v1.DELETE("book/:id", Controllers.DeleteBook)*/
	}

	return r
}

func main() {
	var err error

	config.DB, err = gorm.Open(
		mysql.Open("badminton:badminton@tcp(127.0.0.1:3306)/badminton?charset=utf8&parseTime=True&loc=Local"),
		&gorm.Config{})

	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}

	r := setupRouter()
	r.Run(":9080")
}
