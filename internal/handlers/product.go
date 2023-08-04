package handlers

import "github.com/tirzasrwn/shopping-cart/internal/models"

func (m *module) GetProductByCategory(id int) ([]*models.Product, error) {
	products, err := m.db.dbrepo.GetProductByCategory(id)
	if err != nil {
		return nil, err
	}
	return products, nil
}
