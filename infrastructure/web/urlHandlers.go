package web

import (
	"fmt"
	"net/http"
	"strconv"

	urlService "github.com/Projects-Bots/redirect/infrastructure/service/url"

	coreUrl "github.com/Projects-Bots/redirect/internal/core/url"
	"github.com/gin-gonic/gin"
)

func urlGET(c *gin.Context, urlSrv *urlService.UrlService) {
	userID, err := strconv.Atoi(c.Param("userID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	urls, err := urlSrv.GetAllUrlsByUser(c, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to fetch URLs for user: %v", err)})
		return
	}

	c.JSON(http.StatusOK, urls)
}

func urlPOST(c *gin.Context, urlSrv *urlService.UrlService) {
	var url coreUrl.Url
	if err := c.BindJSON(&url); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	newUrl, err := urlSrv.AddUrl(c, url)
	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to fetch URLs for user: %v", err)})
		return
	}

	c.JSON(http.StatusOK, newUrl)
}

func urlPATCH(c *gin.Context, urlSrv *urlService.UrlService) {
	var url coreUrl.Url
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid URL ID"})
		return
	}
	if err := c.BindJSON(&url); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	updatedUrl, err := urlSrv.UpdateUrl(c, id, url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, updatedUrl)
}
func urlDELETE(c *gin.Context, urlSrv *urlService.UrlService) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid URL ID"})
		return
	}

	deletedUrl, err := urlSrv.DeleteUrl(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to fetch URLs for user: %v", err)})
		return
	}

	c.JSON(http.StatusOK, deletedUrl)
}
