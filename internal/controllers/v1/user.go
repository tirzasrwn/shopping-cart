package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tirzasrwn/shopping-cart/internal/handlers"
	"github.com/tirzasrwn/shopping-cart/internal/utils"
)

// Get User By Email
//
//	@Tags			User
//	@Summary		Get user by email
//	@Description	This is API to get user by email
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
