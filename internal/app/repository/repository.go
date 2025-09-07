package repository

import (
	"fmt"
	"strings"
)

type Repository struct{}

func NewRepository() (*Repository, error) {
	return &Repository{}, nil
}

type Order struct {
	ID    int
	Title string
}

func (r *Repository) GetOrders() ([]Order, error) {
	orders := []Order{
		{
			ID:    1,
			Title: "Order 1",
		},
		{
			ID:    2,
			Title: "Order 2",
		}, {
			ID:    3,
			Title: "Order 3",
		},
	}
	if len(orders) == 0 {
		return nil, fmt.Errorf("orders not found")
	}
	return orders, nil
}

func (r Repository) GetOrder(id int) (Order, error) {
	orders, err := r.GetOrders()
	if err != nil {
		return Order{}, err
	}

	for _, order := range orders {
		if order.ID == id {
			return order, nil
		}
	}
	return Order{}, fmt.Errorf("order not found")
}

func (r Repository) GetOrderByTitle(title string) ([]Order, error) {
	orders, err := r.GetOrders()
	if err != nil {
		return []Order{}, err
	}
	var result []Order
	for _, order := range orders {
		if strings.Contains(strings.ToLower(order.Title), strings.ToLower(title)) {
			result = append(result, order)
		}
	}
	return result, nil
}
