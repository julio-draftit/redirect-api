package web

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	"database/sql"

	"github.com/Projects-Bots/redirect/infrastructure/database"
	"github.com/Projects-Bots/redirect/infrastructure/repository/access"
	"github.com/Projects-Bots/redirect/infrastructure/repository/redirect"
	"github.com/Projects-Bots/redirect/infrastructure/repository/url"
	"github.com/Projects-Bots/redirect/infrastructure/repository/user"
	accessService "github.com/Projects-Bots/redirect/infrastructure/service/access"
	redirectService "github.com/Projects-Bots/redirect/infrastructure/service/redirect"
	urlService "github.com/Projects-Bots/redirect/infrastructure/service/url"
	userService "github.com/Projects-Bots/redirect/infrastructure/service/user"
	"github.com/Projects-Bots/redirect/infrastructure/web/site"

	"github.com/gin-gonic/gin"
)

func initServices(db *sql.DB) (*urlService.UrlService, *redirectService.RedirectService, *accessService.AccessService, *userService.UserService) {
	urlRepository := url.NewUrlRepository(db)
	redirectRepository := redirect.NewRedirectRepository(db)
	accessRepository := access.NewAccessRepository(db)
	userRepository := user.NewUserRepository(db)

	return urlService.NewUrlService(urlRepository),
		redirectService.NewRedirectService(redirectRepository),
		accessService.NewAccessService(accessRepository),
		userService.NewUserService(userRepository)
}

func setupRouter(urlSrv *urlService.UrlService, redirectSrv *redirectService.RedirectService, accessSrv *accessService.AccessService, userSrv *userService.UserService) *gin.Engine {
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

	urlGroup := router.Group("/url")
	{
		urlGroup.GET("/users/:userID", func(c *gin.Context) {
			urlGET(c, urlSrv)
		})

		urlGroup.POST("/", func(c *gin.Context) {
			urlPOST(c, urlSrv)
		})

		urlGroup.PATCH("/:id", func(c *gin.Context) {
			urlPATCH(c, urlSrv)
		})

		urlGroup.DELETE("/:id", func(c *gin.Context) {
			urlDELETE(c, urlSrv)
		})
	}

	userGroup := router.Group("/user")
	{
		userGroup.GET("/:userID", func(c *gin.Context) {
			userGET(c, userSrv)
		})

		userGroup.GET("/", func(c *gin.Context) {
			userGETList(c, userSrv)
		})

		userGroup.POST("/", func(c *gin.Context) {
			userPOST(c, userSrv)
		})

		userGroup.POST("/auth", func(c *gin.Context) {
			userPOSTAuth(c, userSrv)
		})

		userGroup.PATCH("/:id", func(c *gin.Context) {
			userPATCH(c, userSrv)
		})

		userGroup.DELETE("/:id", func(c *gin.Context) {
			userDELETE(c, userSrv)
		})
	}

	site.NewHandler(*urlSrv, *redirectSrv, *accessSrv, *userSrv).AddRouter(router)

	return router
}

func Handler() {
	db, err := database.NewMysql()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	urlSrv, redirectSrv, accessSrv, userSrv := initServices(db)
	router := setupRouter(urlSrv, redirectSrv, accessSrv, userSrv)

	port := os.Getenv("PORT")
	if port == "" {
		port = "80"
	}

	if err := router.Run(":" + port); err != nil {
		log.Panicf("error: %s", err)
	}
}
