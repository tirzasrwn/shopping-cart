package v1

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tirzasrwn/shopping-cart/internal/handlers"
	"github.com/tirzasrwn/shopping-cart/internal/utils"
)

// get product
//
//	@Tags			public
//	@Summary		get prouct by category_id
//	@Description	this is api to get product by category_id
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
