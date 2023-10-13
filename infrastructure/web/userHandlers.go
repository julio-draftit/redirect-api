package web

import (
	"fmt"
	"net/http"
	"strconv"

	userService "github.com/Projects-Bots/redirect/infrastructure/service/user"

	coreUser "github.com/Projects-Bots/redirect/internal/core/user"
	"github.com/gin-gonic/gin"
)

func userGET(c *gin.Context, userSrv *userService.UserService) {
	userID, err := strconv.Atoi(c.Param("userID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid User ID"})
		return
	}

	users, err := userSrv.SelectUser(c, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to fetch Users for user: %v", err)})
		return
	}

	c.JSON(http.StatusOK, users)
}

func userGETList(c *gin.Context, userSrv *userService.UserService) {
	users, err := userSrv.ListUser(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to fetch Users for user: %v", err)})
		return
	}

	c.JSON(http.StatusOK, users)
}

func userPOST(c *gin.Context, userSrv *userService.UserService) {
	var user coreUser.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	newUser, err := userSrv.AddUser(c, user)
	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to fetch Users for user: %v", err)})
		return
	}

	c.JSON(http.StatusOK, newUser)
}

func userPOSTAuth(c *gin.Context, userSrv *userService.UserService) {
	var user coreUser.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	users, err := userSrv.Auth(c, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to fetch Users for user: %v", err)})
		return
	}

	// Se users for nil, significa que o usuário não foi encontrado
	if users == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, users)
}

func userPATCH(c *gin.Context, userSrv *userService.UserService) {
	var user coreUser.User
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid User ID"})
		return
	}
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	updatedUser, err := userSrv.UpdateUser(c, id, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, updatedUser)
}
func userDELETE(c *gin.Context, userSrv *userService.UserService) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid User ID"})
		return
	}

	deletedUser, err := userSrv.DeleteUser(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to fetch Users for user: %v", err)})
		return
	}

	c.JSON(http.StatusOK, deletedUser)
}
