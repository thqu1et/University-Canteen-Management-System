package controllers


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
