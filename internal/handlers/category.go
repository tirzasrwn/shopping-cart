package handlers

import "github.com/tirzasrwn/shopping-cart/internal/models"

func (m *module) GetCategory() ([]*models.Category, error) {
	categories, err := m.db.dbrepo.GetCategories()
	if err != nil {
		return nil, err
	}
	return categories, nil
}
