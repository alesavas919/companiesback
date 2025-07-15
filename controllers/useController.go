package controllers

import (
	"context"
	"net/http"

	"companies/database"
	"companies/models"
	"companies/security"

	"github.com/gin-gonic/gin"
)

func UserOpenTest(c *gin.Context) {
	c.JSON(http.StatusOK, "initialized 2")
}

func UserOpenSecretTest(c *gin.Context) {
	c.JSON(http.StatusOK, security.SecretString("test"))
}

func UserCreateUser(c *gin.Context) {
	var user models.User
	errStr := c.ShouldBindJSON(&user)
	if errStr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errStr.Error()})
		return
	}

	_, err := database.DB.Exec(context.Background(),
		"INSERT INTO users (name, email, password) VALUES ($1, $2, $3)",
		user.Name, user.Email, user.Password)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

func UserGetUsers(c *gin.Context) {
	rows, err := database.DB.Query(context.Background(), "SELECT id, name, email FROM users")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		users = append(users, user)
	}
	c.JSON(http.StatusOK, users)
}

func UserGetUserByID(c *gin.Context) {
	id := c.Param("id")

	var user models.User
	err := database.DB.QueryRow(context.Background(),
		"SELECT id, name, email FROM users WHERE id = $1", id, id).
		Scan(&user.ID, &user.Name, &user.Email)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
		return
	}

	c.JSON(http.StatusOK, user)
}
