package cart

import (
	"fmt"

	"github.com/adhupraba/ecom/types"
)

func getCartItemIDs(items []types.CartItem) ([]int, error) {
	productIds := make([]int, len(items))

	for i, item := range items {
		if item.Quantity <= 0 {
			return nil, fmt.Errorf("invalid quantity for the product %d", item.ProductID)
		}

		productIds[i] = item.ProductID
	}

	return productIds, nil
}

func (h *Handler) createOrder(products []types.Product, items []types.CartItem, userId int) (int, float64, error) {
	productMap := make(types.ProductMap)

	for _, product := range products {
		productMap[product.ID] = product
	}

	// check if all products are actually in stock
	if err := checkIfCartIsInStock(items, productMap); err != nil {
		return 0, 0, err
	}

	// calculate the total price
	totalPrice := calculateTotalPrice(items, productMap)

	// reduce quantity of products in db
	for _, item := range items {
		product := productMap[item.ProductID]
		product.Quantity -= item.Quantity

		h.productStore.UpdateProduct(product)
	}

	// create the order
	orderId, err := h.store.CreateOrder(types.Order{
		UserID:  userId,
		Total:   totalPrice,
		Status:  "PENDING",
		Address: "some address",
	})

	if err != nil {
		return 0, 0, err
	}

	// create order items
	for _, item := range items {
		err = h.store.CreateOrderItem(types.OrderItem{
			OrderID:   orderId,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     productMap[item.ProductID].Price,
		})

		if err != nil {
			return 0, 0, err
		}
	}

	return orderId, totalPrice, nil
}

func checkIfCartIsInStock(items []types.CartItem, productMap types.ProductMap) error {
	if len(items) == 0 {
		return fmt.Errorf("cart is empty")
	}

	for _, item := range items {
		product, ok := productMap[item.ProductID]

		if !ok {
			return fmt.Errorf("product %d is not available in the store, please refresh your cart", item.ProductID)
		}

		if product.Quantity < item.Quantity {
			return fmt.Errorf("product %s is not available in the quantity requested", product.Name)
		}
	}

	return nil
}

func calculateTotalPrice(items []types.CartItem, productMap types.ProductMap) float64 {
	var total float64

	for _, item := range items {
		product := productMap[item.ProductID]
		total += product.Price * float64(product.Quantity)
	}

	return total
}
