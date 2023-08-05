package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tirzasrwn/shopping-cart/internal/handlers"
	"github.com/tirzasrwn/shopping-cart/internal/utils"
)

// get category
//
//	@Tags			public
//	@Summary		get category
//	@Description	this is api to get category
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
