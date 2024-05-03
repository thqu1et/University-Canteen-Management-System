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
	// Define a struct for the request body
	var body struct {
		FirstName string
		LastName  string
		Email     string
		Password  string
	}

	// Bind the incoming JSON to body and handle errors
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to read body: " + err.Error(),
		})
		return
	}

	// Hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to hash password",
		})
		return
	}

	// Create the user in the database
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

	// Respond with success
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

	// Bind the incoming JSON to body and handle errors
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body: " + err.Error(),
		})
		return
	}

	// Look up requested user
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
	// Get user ID from the path parameter
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

	// Respond with user data
	c.JSON(http.StatusOK, user)
}

func GetUsers(c *gin.Context) {
	var users []models.User

	// Retrieve all users from the database
	result := database.DB.Find(&users)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve users: " + result.Error.Error(),
		})
		return
	}

	// Respond with the list of users
	c.JSON(http.StatusOK, users)
}

func GetMenu(c *gin.Context) {
	var menuItems []models.MenuItem

	result := database.DB.Find(&menuItems)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve menu items: " + result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, menuItems)
}

func CreateMenu(c *gin.Context) {
	var menuItem models.MenuItem
	if err := c.BindJSON(&menuItem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body: " + err.Error(),
		})
		return
	}

	result := database.DB.Create(&menuItem)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create menu item: " + result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Menu item created successfully",
	})
}

func UpdateMenu(c *gin.Context) {
	id := c.Param("id")
	var menuItem models.MenuItem

	result := database.DB.First(&menuItem, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Menu item not found",
		})
		return
	}

	if err := c.BindJSON(&menuItem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body: " + err.Error(),
		})
		return
	}

	database.DB.Save(&menuItem)
	c.JSON(http.StatusOK, gin.H{
		"message": "Menu item updated successfully",
	})
}

func DeleteMenu(c *gin.Context) {
	id := c.Param("id")
	result := database.DB.Delete(&models.MenuItem{}, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete menu item: " + result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Menu item deleted successfully",
	})
}
