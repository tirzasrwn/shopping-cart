package handlers

import "github.com/tirzasrwn/shopping-cart/internal/models"

func (m *module) GetUserByEmail(email string) (user *models.User, err error) {
	user, err = m.db.dbrepo.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}
	return
}
