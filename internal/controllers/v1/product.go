package v1

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tirzasrwn/shopping-cart/internal/handlers"
	"github.com/tirzasrwn/shopping-cart/internal/utils"
)

// Get Categories
//
//	@Tags			Product
//	@Summary		Get category
//	@Description	This is API to get category
//	@Param			category_id	query	int	true	"category_id"
//	@Produce		json
//	@Router			/product [get]
func GetProductByCategoryID(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("category_id"))
	if err != nil {
		utils.ErrorJSON(c, err)
		return
	}
	products, err := handlers.Handlers.GetProductByCategory(id)
	if err != nil {
		utils.ErrorJSON(c, err)
		return
	}
	utils.WriteJSON(c, http.StatusOK, products)
}
