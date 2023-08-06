package handlers

func (m *module) InsertOrder(cartID int, productID int, quantity int) (int, error) {
	isExist, err := m.db.dbrepo.CheckOrderExist(cartID, productID)
	if err != nil {
		return 0, err
	}
	if !isExist {
		id, err := m.db.dbrepo.InsertOrder(cartID, productID, quantity)
		if err != nil {
			return 0, err
		}
		return id, nil
	}
	id, err := m.db.dbrepo.UpdateQuantity(cartID, productID, quantity)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (m *module) DeleteOrder(orderID int) error {
	err := m.db.dbrepo.DeleteOder(orderID)
	if err != nil {
		return err
	}
	return nil
}
