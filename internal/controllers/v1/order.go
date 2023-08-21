package v1

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tirzasrwn/shopping-cart/internal/handlers"
	"github.com/tirzasrwn/shopping-cart/internal/utils"
)

type InsertOrderPayload struct {
	ProductID int `json:"product_id" example:"5"`
	Quantity  int `json:"quantity" example:"1"`
}

// InsertOrder godoc
//
//	@Security		UserAuth
//	@Tags			user
//	@Summary		create new order or update the quantity
//	@Description	this api to post new order or update the quantity
//	@Description	prouduct_id can be found at /product
//	@Param			payload	body	InsertOrderPayload	true	"body payload"
//	@Produce		json
//	@Router			/user/order [post]
func InsertOrder(c *gin.Context) {
	requestBody, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		utils.ErrorJSON(c, err)
	}
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(requestBody))
	var request InsertOrderPayload
	err = c.ShouldBind(&request)
	if err != nil {
		utils.ErrorJSON(c, err)
		return
	}
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(requestBody))

	userEmail, err := getUserEmailFromContex(c)
	if err != nil {
		utils.ErrorJSON(c, err)
		return
	}
	cartID, err := handlers.Handlers.GetUserCartByEmail(userEmail)
	if err != nil {
		utils.ErrorJSON(c, err)
		return
	}

	orderID, err := handlers.Handlers.InsertOrder(cartID, request.ProductID, request.Quantity)
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

// DeleteOrder godoc
//
//	@Security		UserAuth
//	@Tags			user
//	@Summary		delete order by order id
//	@Description	this api to delete order by order_id
//	@Description	order_id can be found at get /user/order
//	@Param			order_id	path	int	true	"order id"	default(1)
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

type CheckoutOrderPayload struct {
	Money float64 `json:"money" example:"100000"`
}

// Checkout godoc
//
//	@Security		UserAuth
//	@Tags			user
//	@Summary		checkout all order in cart
//	@Description	this api to checkout and make payment transactions
//	@Description	total payment can be found at get /user/order
//	@Param			payload	body	CheckoutOrderPayload	true	"body payload"
//	@Produce		json
//	@Router			/user/checkout [post]
func Checkout(c *gin.Context) {
	requestBody, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		utils.ErrorJSON(c, err)
	}
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(requestBody))
	var request CheckoutOrderPayload
	err = c.ShouldBind(&request)
	if err != nil {
		utils.ErrorJSON(c, err)
		return
	}
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(requestBody))
	userEmail, err := getUserEmailFromContex(c)
	if err != nil {
		utils.ErrorJSON(c, err)
		return
	}
	changeMoney, err := handlers.Handlers.CheckoutOrder(request.Money, userEmail)
	if err != nil {
		utils.ErrorJSON(c, err)
		return
	}
	data := &utils.JSONResponse{
		Error:   false,
		Message: fmt.Sprintf("success to checkout, your change money is %.2f", changeMoney),
		Data:    nil,
	}
	utils.WriteJSON(c, http.StatusOK, data)
}
