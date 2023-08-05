package v1

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tirzasrwn/shopping-cart/internal/controllers/middleware"
	"github.com/tirzasrwn/shopping-cart/internal/handlers"
	"github.com/tirzasrwn/shopping-cart/internal/models"
	"github.com/tirzasrwn/shopping-cart/internal/utils"
)

func getUserEmailFromContex(c *gin.Context) (string, error) {
	_, claims, _ := middleware.AdminAuth.GetTokenFromHeaderAndVerify(c)
	userEmail := fmt.Sprintf("%v", claims["email"])
	if userEmail == "" {
		return "", errors.New("user email in token not found")
	}
	return userEmail, nil
}

// get user information
//
//	@Security		UserAuth
//	@Tags			user
//	@Summary		get user information
//	@Description	this is API to get user information
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
	utils.WriteJSON(c, http.StatusOK, user)
}

// create new user
//
//	@Tags			public
//	@Summary		create new user
//	@Description	this is API to create new user
//	@Param			email		query	string	true	"email"
//	@Param			password	query	string	true	"password"
//	@Produce		json
//	@Router			/user [post]
func PostUser(c *gin.Context) {
	email := c.Query("email")
	password := c.Query("password")
	var user = models.User{
		Email:    email,
		Password: password,
	}
	_, err := handlers.Handlers.InsertUser(&user)
	if err != nil {
		utils.ErrorJSON(c, err)
		return
	}
	data := utils.JSONResponse{
		Error:   false,
		Message: "Success create user!",
		Data:    nil,
	}
	utils.WriteJSON(c, http.StatusOK, data)
}

// get user order
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
	utils.WriteJSON(c, http.StatusOK, products)
}

// get user payment
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
	payment, err := handlers.Handlers.GetUserPayment(email)
	if err != nil {
		utils.ErrorJSON(c, err)
		return
	}
	utils.WriteJSON(c, http.StatusOK, payment)
}

type AuthenticatePayload struct {
	Email    string `json:"email" example:"user1@example.com"`
	Password string `json:"password" example:"user1"`
}

// Login
//
//	@Tags			public
//	@Summary		login
//	@Description	this is api to authenticate user
//	@Param			payload	body	AuthenticatePayload	true	"body payload"
//	@Produce		json
//	@Router			/login [post]
func Authenticate(c *gin.Context) {
	var requestPayload AuthenticatePayload
	err := c.ShouldBind(&requestPayload)
	if err != nil {
		utils.ErrorJSON(c, err, http.StatusBadRequest)
		return
	}
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

	refreshCookie := middleware.AdminAuth.GetRefreshCookie(tokens.RefreshToken)
	http.SetCookie(c.Writer, refreshCookie)

	utils.WriteJSON(c, http.StatusOK, tokens)
}
