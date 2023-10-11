package web

import (
	"github.com/Projects-Bots/redirect/infrastructure/database"
	"github.com/Projects-Bots/redirect/infrastructure/repository/access"
	"github.com/Projects-Bots/redirect/infrastructure/repository/redirect"
	"github.com/Projects-Bots/redirect/infrastructure/repository/url"
	accessService "github.com/Projects-Bots/redirect/infrastructure/service/access"
	redirectService "github.com/Projects-Bots/redirect/infrastructure/service/redirect"
	urlService "github.com/Projects-Bots/redirect/infrastructure/service/url"
	"github.com/Projects-Bots/redirect/infrastructure/web/site"
	"github.com/gin-gonic/gin"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)

func Handler() {
	db, err := database.NewMysql()
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	urlRepository := url.NewUrlRepository(db)
	urlService := urlService.NewUrlService(urlRepository)

	redirectRepository := redirect.NewRedirectRepository(db)
	redirectService := redirectService.NewRedirectService(redirectRepository)

	accessRepository := access.NewAccessRepository(db)
	accessService := accessService.NewAccessService(accessRepository)

	urlHanldler := site.NewHandler(*urlService, *redirectService, *accessService)

	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()

	router.SetFuncMap(template.FuncMap{
		"upper": strings.ToUpper,
	})

	router.LoadHTMLGlob("./infrastructure/web/site/templates/*")

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	urlHanldler.AddRouter(router)

	port := os.Getenv("PORT")
	if port == "" {
		port = "80"
	}

	if err := router.Run(":" + port); err != nil {
		log.Panicf("error: %s", err)
	}
}
