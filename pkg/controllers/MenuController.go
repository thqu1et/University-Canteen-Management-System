package controllers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"postgresSQLProject/pkg/database"
	"postgresSQLProject/pkg/models"
	"strconv"
)

// GetMenuItems godoc
// @Summary Retrieve all menu items
// @Description Get details of all menu items in the database
// @Tags menu
// @Accept  json
// @Produce  json
// @Success 200 {array} models.MenuItem
// @Failure 500 {object} object {"error": "description of error"}
// @Router /menu/items [get]
func GetMenuItems(c *gin.Context) {
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

// GetMenuItem godoc
// @Summary Retrieve a single menu item
// @Description Get details of a specific menu item by ID
// @Tags menu
// @Accept  json
// @Produce  json
// @Param id path int true "Menu Item ID"
// @Success 200 {object} models.MenuItem
// @Failure 400 {object} object {"error": "Invalid ID format"}
// @Failure 404 {object} object {"error": "Menu item not found"}
// @Failure 500 {object} object {"error": "Error retrieving menu item"}
// @Router /menu/items/{id} [get]
func GetMenuItem(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid ID format",
		})
		return
	}

	var menuItem models.MenuItem

	result := database.DB.First(&menuItem, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Menu item not found",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error retrieving menu item: " + result.Error.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": menuItem,
	})
}

// CreateMenuItem godoc
// @Summary Create a new menu item
// @Description Add a new menu item to the database
// @Tags menu
// @Accept  json
// @Produce  json
// @Param menuItem body models.MenuItem required "Menu Item Info"
// @Success 200 {object} object {"message": "Menu item created successfully"}
// @Failure 400 {object} object {"error": "Failed to read body"}
// @Failure 400 {object} object {"error": "Failed to create menu item"}
// @Router /menu/items [post]
func CreateMenuItem(c *gin.Context) {
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

// UpdateMenuItem godoc
// @Summary Update an existing menu item
// @Description Update details of an existing menu item by ID
// @Tags menu
// @Accept  json
// @Produce  json
// @Param id path int true "Menu Item ID"
// @Param menuItem body map[string]interface{} true "Menu Item Update Data"
// @Success 200 {object} object {"message": "Menu item updated successfully"}
// @Failure 400 {object} object {"error": "Failed to read body"}
// @Failure 404 {object} object {"error": "Menu item not found"}
// @Failure 500 {object} object {"error": "Failed to update menu item"}
// @Router /menu/items/{id} [put]
func UpdateMenuItem(c *gin.Context) {
	id := c.Param("id")
	var existingMenuItem models.MenuItem

	result := database.DB.First(&existingMenuItem, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Menu item not found",
		})
		return
	}

	var updateData map[string]interface{}
	if err := c.BindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body: " + err.Error(),
		})
		return
	}

	if err := database.DB.Model(&existingMenuItem).Updates(updateData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update menu item: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Menu item updated successfully",
	})
}

// DeleteMenuItem godoc
// @Summary Delete a menu item
// @Description Remove a menu item from the database by ID
// @Tags menu
// @Accept  json
// @Produce  json
// @Param id path int true "Menu Item ID"
// @Success 200 {object} object {"message": "Menu item deleted successfully"}
// @Failure 500 {object} object {"error": "Failed to delete menu item"}
// @Router /menu/items/{id} [delete]
func DeleteMenuItem(c *gin.Context) {
	id := c.Param("id")
	result := database.DB.Delete(&models.MenuItem{}, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete menu item: " + result.Error.Error(),
		})
		return
	}
}
