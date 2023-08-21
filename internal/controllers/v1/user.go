package v1

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tirzasrwn/shopping-cart/internal/controllers/middleware"
	"github.com/tirzasrwn/shopping-cart/internal/handlers"
	"github.com/tirzasrwn/shopping-cart/internal/models"
	"github.com/tirzasrwn/shopping-cart/internal/utils"
)

type CredentialPayload struct {
	Email    string `json:"email" example:"user1@example.com"`
	Password string `json:"password" example:"user1"`
}

func getUserEmailFromContex(c *gin.Context) (string, error) {
	_, claims, _ := middleware.AdminAuth.GetTokenFromHeaderAndVerify(c)
	userEmail := fmt.Sprintf("%v", claims["email"])
	if userEmail == "" {
		return "", errors.New("user email in token not found")
	}
	return userEmail, nil
}

// GetUserInformation godoc
//
//	@Security		UserAuth
//	@Tags			user
//	@Summary		get all user information
//	@Description	this is api to get all user information
//	@Produce		json
//	@Router			/user [get]
func GetUserInformation(c *gin.Context) {
	var err error
	email, err := getUserEmailFromContex(c)
	if err != nil {
		utils.ErrorJSON(c, err)
		return
	}
	user, err := handlers.Handlers.GetUserByEmail(email)
	if err != nil {
		utils.ErrorJSON(c, err)
		return
	}
	cartID, err := handlers.Handlers.GetUserCartByEmail(email)
	if err != nil {
		utils.ErrorJSON(c, err)
		return
	}
	var userData = struct {
		models.User
		CartID int `json:"cart_id"`
	}{
		User:   *user,
		CartID: cartID,
	}
	data := utils.JSONResponse{
		Error:   false,
		Message: "success get user information",
		Data:    userData,
	}
	utils.WriteJSON(c, http.StatusOK, data)
}

// PostUser godoc
//
//	@Tags			public
//	@Summary		register new user
//	@Description	this is api to register new user
//	@Param			payload	body	CredentialPayload	true	"body payload"
//	@Produce		json
//	@Router			/register [post]
func PostUser(c *gin.Context) {
	requestBody, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		utils.ErrorJSON(c, err)
	}
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(requestBody))
	var requestPayload CredentialPayload
	err = c.ShouldBind(&requestPayload)
	if err != nil {
		utils.ErrorJSON(c, err, http.StatusBadRequest)
		return
	}
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(requestBody))
	var user = models.User{
		Email:    requestPayload.Email,
		Password: requestPayload.Password,
	}
	_, err = handlers.Handlers.InsertUser(&user)
	if err != nil {
		utils.ErrorJSON(c, err)
		return
	}
	data := utils.JSONResponse{
		Error:   false,
		Message: "success create user",
		Data:    nil,
	}
	utils.WriteJSON(c, http.StatusOK, data)
}

// GetUserOrder godoc
//
//	@Security		UserAuth
//	@Tags			user
//	@Summary		get user order
//	@Description	this api is to get user order
//	@Produce		json
//	@Router			/user/order [get]
func GetUserOrder(c *gin.Context) {
	var err error
	email, err := getUserEmailFromContex(c)
	if err != nil {
		utils.ErrorJSON(c, err)
		return
	}
	products, err := handlers.Handlers.GetUserOrder(email)
	if err != nil {
		utils.ErrorJSON(c, err)
		return
	}
	if products == nil {
		data := utils.JSONResponse{
			Error:   false,
			Message: "your cart is empty",
			Data:    nil,
		}
		utils.WriteJSON(c, http.StatusOK, data)
		return
	}
	var totalPayment float64
	for _, p := range products {
		totalPayment = totalPayment + (float64(p.Quantity) * p.Price)
	}
	data := utils.JSONResponse{
		Error:   false,
		Message: fmt.Sprintf("your total payment for checkout is %.2f", totalPayment),
		Data:    products,
	}
	utils.WriteJSON(c, http.StatusOK, data)
}

// GetUserPayment godoc
//
//	@Security		UserAuth
//	@Tags			user
//	@Summary		get user payment
//	@Description	this api is to get user payment
//	@Produce		json
//	@Router			/user/payment [get]
func GetUserPayment(c *gin.Context) {
	var err error
	email, err := getUserEmailFromContex(c)
	if err != nil {
		utils.ErrorJSON(c, err)
		return
	}
	payments, err := handlers.Handlers.GetUserPayment(email)
	if err != nil {
		utils.ErrorJSON(c, err)
		return
	}
	data := utils.JSONResponse{
		Error:   false,
		Message: "success get user payment history",
		Data: map[string]interface{}{
			"payment": payments,
		},
	}
	utils.WriteJSON(c, http.StatusOK, data)
}

// Authenticate godoc
//
//	@Tags			public
//	@Summary		login
//	@Description	this is api to authenticate user then returns jwt token
//	@Param			payload	body	CredentialPayload	true	"body payload"
//	@Produce		json
//	@Router			/login [post]
func Authenticate(c *gin.Context) {
	requestBody, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		utils.ErrorJSON(c, err)
	}
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(requestBody))
	var requestPayload CredentialPayload
	err = c.ShouldBind(&requestPayload)
	if err != nil {
		utils.ErrorJSON(c, err, http.StatusBadRequest)
		return
	}
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(requestBody))
	user, err := handlers.Handlers.GetUserByEmail(requestPayload.Email)
	if err != nil {
		utils.ErrorJSON(c, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}
	valid, err := user.PasswordPlainMatches(requestPayload.Password)
	if err != nil || !valid {
		utils.ErrorJSON(c, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}
	u := middleware.JwtUser{
		ID:    user.ID,
		Email: user.Email,
	}
	tokens, err := middleware.AdminAuth.GenerateTokenPair(&u)
	if err != nil {
		utils.ErrorJSON(c, err, http.StatusBadRequest)
		return
	}

	data := utils.JSONResponse{
		Error:   false,
		Message: "success to login",
		Data:    tokens,
	}

	utils.WriteJSON(c, http.StatusOK, data)
}
