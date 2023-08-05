package v1

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tirzasrwn/shopping-cart/internal/handlers"
	"github.com/tirzasrwn/shopping-cart/internal/utils"
)

// get products
//
//	@Tags			public
//	@Summary		get proucts
//	@Description	this is api to get products
//	@Produce		json
//	@Router			/product [get]
func GetProducts(c *gin.Context) {
	products, err := handlers.Handlers.GetProducts()
	if err != nil {
		utils.ErrorJSON(c, err)
		return
	}
	utils.WriteJSON(c, http.StatusOK, products)
}

// get product
//
//	@Tags			public
//	@Summary		get prouct by category_id
//	@Description	this is api to get product by category_id
//	@Param			category_id	path	int	true	"category_id"	default(1)
//	@Produce		json
//	@Router			/product/{category_id} [get]
func GetProductByCategoryID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("category_id"))
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
