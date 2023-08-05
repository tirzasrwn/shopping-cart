package handlers

import "github.com/tirzasrwn/shopping-cart/internal/models"

func (m *module) GetUserByEmail(email string) (user *models.User, err error) {
	user, err = m.db.dbrepo.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}
	return
}

func (m *module) InsertUser(user *models.User) (int, error) {
	userID, err := m.db.dbrepo.InsertUser(user)
	if err != nil {
		return 0, err
	}
	_, err = m.db.dbrepo.CreateCartForNewUser(userID)
	if err != nil {
		return 0, err
	}
	return userID, nil
}

func (m *module) GetUserOrder(email string) ([]*models.Product, error) {
	products, err := m.db.dbrepo.GetOrderByUserEmail(email)
	if err != nil {
		return nil, err
	}
	return products, nil
}
