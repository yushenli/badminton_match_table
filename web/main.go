package main

import (
	"flag"
	"log"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/yushenli/badminton_match_table/web/lib/config"
	"github.com/yushenli/badminton_match_table/web/www/controller"
	"github.com/yushenli/badminton_match_table/web/www/formatter"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var staticDirFlag = flag.String("static_dir", "./static", "The root directory which contains all the static files")
var templatesDirFlag = flag.String("templates_dir", "./templates", "The root directory which contains all the template files")

func setupRouter() *gin.Engine {
	r := gin.Default()
	formatter.RegisterFormatters(r)

	r.LoadHTMLGlob(filepath.Join(*templatesDirFlag, "*.html"))

	r.GET("/", controller.RenderIndex)
	r.GET("/index.html", controller.RenderIndex)
	r.GET("/rules.html", controller.RenderRules)
	r.GET("/event/:key", controller.RenderEvent)
	r.GET("/admin/change_match_status", controller.ChangeMatchStatus)
	r.GET("/admin/change_break_status", controller.ChangeBreakStatus)
	r.GET("/admin/complete_round", controller.CompleteRound)
	r.GET("/admin/schedule", controller.ScheduleCurrentRound)

	staticFiles := []string{}
	for _, staticFile := range staticFiles {
		r.StaticFile("/"+staticFile, filepath.Join(*staticDirFlag, staticFile))
	}
	staticSubdirs := []string{
		"css",
		"images",
		"js",
	}
	for _, staticDir := range staticSubdirs {
		r.Static("/"+staticDir, filepath.Join(*staticDirFlag, staticDir))
	}

	return r
}

func main() {
	var err error
	flag.Parse()

	config.DB, err = gorm.Open(
		mysql.Open("badminton:badminton@tcp(127.0.0.1:3306)/badminton?charset=utf8&parseTime=True&loc=Local"),
		&gorm.Config{})

	if err != nil {
		log.Printf("Unable to connect to database: %v", err)
		config.DB = nil
	}

	r := setupRouter()
	r.Run(":9080")
}
