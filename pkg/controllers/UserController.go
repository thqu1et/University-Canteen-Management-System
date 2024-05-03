package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"postgresSQLProject/pkg/database"
	"postgresSQLProject/pkg/models"
	"time"
)

func SignUp(c *gin.Context) {
	var body struct {
		FirstName string
		LastName  string
		Email     string
		Password  string
	}

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to read body: " + err.Error(),
		})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to hash password",
		})
		return
	}

	user := models.User{
		FirstName: body.FirstName,
		LastName:  body.LastName,
		Email:     body.Email,
		Password:  string(hash),
	}
	result := database.DB.Create(&user)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to create user: " + result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "user created successfully",
		"user_id": user.ID,
	})
}

func Login(c *gin.Context) {
	var body struct {
		Email    string
		Password string
	}

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body: " + err.Error(),
		})
		return
	}

	user := models.User{}
	result := database.DB.Where("email = ?", body.Email).First(&user)

	if result.Error != nil || user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email or password",
		})
		return
	}

	// Verify the password
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email or password",
		})
		return
	}

	// Generate a JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "JWT secret is not set",
		})
		return
	}

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create token: " + err.Error(),
		})
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(
		"Authorization",
		tokenString,
		3600*24*30,
		"",
		"",
		false,
		true,
	)

	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
	})
}

func Validate(c *gin.Context) {
	user, _ := c.Get("user")

	c.JSON(http.StatusOK, gin.H{
		"message": user,
	})
}

func GetUser(c *gin.Context) {
	userID := c.Param("id")
	var user models.User

	// Retrieve user by ID
	result := database.DB.First(&user, userID)
	if result.Error != nil || result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}

	c.JSON(http.StatusOK, user)
}

func GetUsers(c *gin.Context) {
	var users []models.User

	result := database.DB.Find(&users)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve users: " + result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, users)
}
