package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tirzasrwn/shopping-cart/internal/handlers"
	"github.com/tirzasrwn/shopping-cart/internal/utils"
)

// Get product by category id
//
//	@Tags			Category
//	@Summary		Get product by category id
//	@Description	This is API to get product by category id
//	@Produce		json
//	@Router			/category [get]
func GetCategories(c *gin.Context) {
	categories, err := handlers.Handlers.GetCategory()
	if err != nil {
		utils.ErrorJSON(c, err)
		return
	}
	utils.WriteJSON(c, http.StatusOK, categories)
}
