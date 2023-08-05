package v1

import (
	"net/http"
	"strconv"

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
//	@Security		UserAuth
//	@Tags			user
//	@Summary		post new order or update the quantity
//	@Description	this api to post new order or update the quantity
//	@Param			payload	body	InsertOrderPayload	true	"body payload"
//	@Produce		json
//	@Router			/user/order [post]
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

// delete order
//
//	@Security		UserAuth
//	@Tags			user
//	@Summary		delete order by order id
//	@Description	this api to delete order by order id
//	@Param			order_id	path	int	true	"order id"
//	@Produce		json
//	@Router			/user/order/{order_id} [delete]
func DeleteOrder(c *gin.Context) {
	pathOrderID := c.Param("order_id")
	orderID, err := strconv.Atoi(pathOrderID)
	if err != nil {
		utils.ErrorJSON(c, err)
		return
	}
	err = handlers.Handlers.DeleteOrder(orderID)
	if err != nil {
		utils.ErrorJSON(c, err)
		return
	}
	data := utils.JSONResponse{
		Error:   false,
		Message: "success delete order",
		Data:    nil,
	}
	utils.WriteJSON(c, http.StatusOK, data)
}
