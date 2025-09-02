package vending_machine_rack

import "fmt"

// IGNORE INVENTORY is using racks to store products

const MaxRackCapacity = 10

type Rack struct {
	ID       int
	Product  *Product // Map of product ID to Product
	Quantity int      // Maximum capacity of the rack // remove qty from product
}

func NewRack(id int, product *Product, quantity int) *Rack {
	fmt.Printf("Added rack with ID %d\n", id)
	return &Rack{ID: id, Product: product, Quantity: quantity}
}

func (r *Rack) IsEmpty() bool {
	return r.Quantity == 0
}

func (r *Rack) IsFull() bool {
	return r.Quantity >= MaxRackCapacity // assuming max capacity of rack is 10
}

func (r *Rack) Capacity() int {
	return MaxRackCapacity - r.Quantity
}

func (r *Rack) IsProductAvailable() bool {
	return r.Quantity > 0
}

func (r *Rack) UpdateRack(product *Product, quantity int) {
	fmt.Printf("Updated rack with ID %d\n", r.ID)
	r.Product = product
	r.Quantity = quantity
}

func (r *Rack) AddProduct(quantity int) error {
	if r.Quantity+quantity <= MaxRackCapacity {
		fmt.Printf("added to rack with ID %d\n", r.ID)
		r.Quantity += quantity
		return nil
	}
	return fmt.Errorf("Rack is full")
}

func (r *Rack) RemoveProduct(quantity int) error {
	if r.Quantity-quantity >= 0 {
		fmt.Printf("removed from rack with ID %d\n", r.ID)
		r.Quantity -= quantity
		return nil
	}
	return fmt.Errorf("Quantity out of range: %d", r.Quantity)
}

func (r *Rack) TransactProduct() (*Product, error) {
	if r.Quantity > 0 {
		r.Quantity--
		product := &Product{
			ID:    r.Product.ID,
			Name:  r.Product.Name,
			Price: r.Product.Price,
		}
		fmt.Printf("Transacted rack with ID %d, remaining quantity: %d\n", r.ID, r.Quantity)
		return product, nil
	}
	return nil, fmt.Errorf("product with rack %d is out of stock", r.ID)
}
