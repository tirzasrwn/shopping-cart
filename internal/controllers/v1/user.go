package v1

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tirzasrwn/shopping-cart/internal/controllers/middleware"
	"github.com/tirzasrwn/shopping-cart/internal/handlers"
	"github.com/tirzasrwn/shopping-cart/internal/models"
	"github.com/tirzasrwn/shopping-cart/internal/utils"
)

// get user information
//
//	@Tags			user
//	@Summary		get user by email
//	@Description	this is API to get user by email
//	@Param			email	query	string	true	"email"
//	@Produce		json
//	@Router			/user [get]
func GetUserByEmail(c *gin.Context) {
	var err error
	email := c.Query("email")
	user, err := handlers.Handlers.GetUserByEmail(email)
	if err != nil {
		utils.ErrorJSON(c, err)
		return
	}
	utils.WriteJSON(c, http.StatusOK, user)
}

// create new user
//
//	@Tags			user
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

// get user order by email
//
//	@Tags			user
//	@Summary		get user order by email
//	@Description	this api is to get user order by email
//	@Param			email	query	string	true	"email"
//	@Produce		json
//	@Router			/user/order [get]
func GetUserOrder(c *gin.Context) {
	var err error
	email := c.Query("email")
	products, err := handlers.Handlers.GetUserOrder(email)
	if err != nil {
		utils.ErrorJSON(c, err)
		return
	}
	utils.WriteJSON(c, http.StatusOK, products)
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
