package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tirzasrwn/shopping-cart/internal/handlers"
	"github.com/tirzasrwn/shopping-cart/internal/utils"
)

type InsertOrderPayload struct {
	CardID    int `json:"card_id" example:"1"`
	ProductID int `json:"product_id" example:"5"`
	Quantity  int `json:"quantity" example:"1"`
}

// post or update order
//
//	@Tags			order
//	@Summary		post new order or update the quantity
//	@Description	this api to post new order or update the quantity
//	@Param			payload	body	InsertOrderPayload	true	"body payload"
//	@Produce		json
//	@Router			/order [post]
func InsertOrder(c *gin.Context) {
	var request InsertOrderPayload
	err := c.ShouldBind(&request)
	if err != nil {
		utils.ErrorJSON(c, err)
		return
	}

	orderID, err := handlers.Handlers.InsertOrder(request.CardID, request.ProductID, request.Quantity)
	if err != nil {
		utils.ErrorJSON(c, err)
		return
	}

	data := utils.JSONResponse{
		Error:   false,
		Message: "success insert new order",
		Data:    map[string]interface{}{"order_id": orderID},
	}
	utils.WriteJSON(c, http.StatusOK, data)
}
