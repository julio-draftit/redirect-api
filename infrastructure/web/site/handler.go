package site

import (
	"github.com/Projects-Bots/redirect/infrastructure/service/access"
	"github.com/Projects-Bots/redirect/infrastructure/service/redirect"
	"github.com/Projects-Bots/redirect/infrastructure/service/url"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type Handler struct {
	urlService      url.UrlService
	redirectService redirect.RedirectService
	accessService   access.AccessService
}

func NewHandler(service url.UrlService, redirectService redirect.RedirectService, accessService access.AccessService) *Handler {
	return &Handler{
		urlService:      service,
		redirectService: redirectService,
		accessService:   accessService,
	}
}

func (h *Handler) AddRouter(router *gin.Engine) {
	router.GET("/", func(c *gin.Context) {
		h.pageError(c, "Que pena!", "Essa página não existe")
		return
	})

	router.GET("/:url", h.Redirect)
}

func (h *Handler) Redirect(c *gin.Context) {
	link := c.Param("url")

	u, err := h.urlService.GetUrl(c, link)
	if err != nil {
		log.Println("Error getting")
		h.pageError(c, "Que pena!", "Ocorreu um erro ao processar a página")
		return
	}

	if u == nil {
		h.pageError(c, "Que pena!", "A url digitada não existe. Verifique se digitou corretamente.")
		return
	}

	redirect, err := h.redirectService.GetUrl(c, u.ID)
	if err != nil {
		h.pageError(c, "Que pena!", "Ocorreu um erro ao processar a página")
		return
	}

	if redirect == nil {
		h.pageError(c, "Que pena!", "A url digitada não existe. Verifique se digitou corretamente.")
		return
	}

	err = h.redirectService.UpdateUrl(c, redirect.ID)
	if err != nil {
		h.pageError(c, "Que pena!", "Ocorreu um erro ao processar a página")
		return
	}

	_, err = h.accessService.Save(c, redirect.ID)
	if err != nil {
		h.pageError(c, "Que pena!", "Ocorreu um erro ao processar a página")
		return
	}

	c.HTML(http.StatusOK, "redirect.html", gin.H{
		"redirect": redirect.Url,
		"pixel":    u.Pixel,
	})
}

func (h *Handler) pageError(c *gin.Context, title, message string) {
	c.HTML(http.StatusOK, "error.html", gin.H{
		"title":   title,
		"message": message,
	})
}
