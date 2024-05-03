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

// SignUp godoc
// @Summary Sign up a new user
// @Description Register a new user with their first name, last name, email, and password
// @Tags user
// @Accept  json
// @Produce  json
// @Param user body struct{FirstName string; LastName string; Email string; Password string} true "SignUp Info"
// @Success 200 {object} object{"message": "user created successfully", "user_id": "ID"}
// @Failure 400 {object} object{"error": "error message"}
// @Router /signup [post]
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

// Login godoc
// @Summary Log in a user
// @Description Log in with email and password
// @Tags user
// @Accept  json
// @Produce  json
// @Param credentials body struct{Email string; Password string} true "Login Credentials"
// @Success 200 {object} object{"token": "JWT Token"}
// @Failure 400 {object} object{"error": "Invalid email or password"}
// @Failure 500 {object} object{"error": "JWT secret is not set"}
// @Router /login [post]
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

// Validate godoc
// @Summary Validate user token
// @Description Checks if the user's token is valid and returns the user info if valid
// @Tags user
// @Accept  json
// @Produce  json
// @Success 200 {object} object{"message": "user info"}
// @Router /validate [get]
func Validate(c *gin.Context) {
	user, _ := c.Get("user")

	c.JSON(http.StatusOK, gin.H{
		"message": user,
	})
}

// GetUser godoc
// @Summary Get a single user
// @Description Retrieves a single user by user ID
// @Tags user
// @Accept  json
// @Produce  json
// @Param id path int true "User ID"
// @Success 200 {object} models.User
// @Failure 404 {object} object{"error": "User not found"}
// @Router /users/{id} [get]
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

// GetUsers godoc
// @Summary Get all users
// @Description Retrieves all registered users
// @Tags user
// @Accept  json
// @Produce  json
// @Success 200 {array} []models.User
// @Failure 500 {object} object{"error": "Failed to retrieve users"}
// @Router /users [get]
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
