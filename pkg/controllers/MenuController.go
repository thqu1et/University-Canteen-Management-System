package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"postgresSQLProject/pkg/database"
	"postgresSQLProject/pkg/models"
)

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
