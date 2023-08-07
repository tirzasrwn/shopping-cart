package handlers

import (
	"fmt"

	"github.com/tirzasrwn/shopping-cart/internal/models"
)

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

func (m *module) GetUserCartByEmail(email string) (int, error) {
	cartID, err := m.db.dbrepo.GetUserCartByEmail(email)
	if err != nil {
		return 0, err
	}
	return cartID, nil
}

func (m *module) GetUserOrder(email string) ([]*models.ProductOrder, error) {
	products, err := m.db.dbrepo.GetOrderByUserEmail(email)
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (m *module) GetUserPayment(email string) ([]*models.ProductPayment, error) {
	payment, err := m.db.dbrepo.GetPaymentByUserEmail(email)
	if err != nil {
		return nil, err
	}
	return payment, nil
}

func (m *module) CheckoutOrder(money float64, email string) (float64, error) {
	user, err := m.db.dbrepo.GetUserByEmail(email)
	if err != nil {
		return 0, err
	}
	// Get all the product order.
	products, err := m.db.dbrepo.GetOrderByUserEmail(email)
	if err != nil {
		return 0, err
	}
	// Count the total price.
	var total float64
	for _, p := range products {
		total = total + (float64(p.Quantity) * p.Price)
	}
	// Check the payment.
	var changeMoney float64
	if total > money {
		moneyNeed := total - money
		return 0, fmt.Errorf("insufficient of money, total price is %.2f, you need %.2f more", total, moneyNeed)
	}
	changeMoney = money - total
	// Move order to payment if success and return change money.
	for _, p := range products {
		_, err := m.db.dbrepo.InsertSuccessPayment(user.ID, p.ID, p.Quantity)
		if err != nil {
			return 0, err
		}
		err = m.db.dbrepo.DeleteOder(p.OrderID)
		if err != nil {
			return 0, err
		}
	}
	return changeMoney, nil
}
